//go:build !wioterminal

package wiointerface

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/png"

	"github.com/veandco/go-sdl2/sdl"
)

/*
TODO BACKLIGHT
-backlight korostus
-Turn sdlwio f "flip"
*/

type Sdlwio struct {
	window       *sdl.Window
	renderer     *sdl.Renderer
	backlightNow bool

	imgsurf *sdl.Surface

	winNow_x     int16
	winNow_y     int16
	winNow_w     int16
	winNow_h     int16
	winIndex     int
	writeStarted bool

	winScale        int
	backgroundFrame SdlwioBackgroundFrame //BackgroundFrame is around actual display area

	Landscape bool
	Flipped   bool
}

type SdlwioBackgroundFrame struct {
	texture    *sdl.Texture
	sourceArea image.Rectangle
}

const (
	WIODISPLAY_W = 320
	WIODISPLAY_H = 240
	TITLEDISPLAY = "wio simulator"
)

//go:embed wiobackground.png
var backgroundpngbytes []byte

const (
	BACKGROUNDPNG_VIEW_X0 = 85
	BACKGROUNDPNG_VIEW_Y0 = 111
	BACKGROUNDPNG_VIEW_X1 = 917
	BACKGROUNDPNG_VIEW_Y1 = 660
)

func InitSdlwio() (*Sdlwio, error) {
	var err error
	result := Sdlwio{}

	bgImage, errPng := png.Decode(bytes.NewBuffer(backgroundpngbytes))
	if errPng != nil {
		return nil, fmt.Errorf("background png decoding error %s", errPng)
	}
	bgBorder := bgImage.Bounds()

	bgSurface, errBgSurface := sdl.CreateRGBSurface(0, int32(bgBorder.Dx()), int32(bgBorder.Dy()), 32, 0, 0, 0, 0)
	if errBgSurface != nil {
		return nil, fmt.Errorf("create background surface fail err=%s", errBgSurface)
	}
	//Copy
	for y := int(0); y < int(bgSurface.H); y++ {
		for x := int(0); x < int(bgSurface.W); x++ {
			bgSurface.Set(x, y, bgImage.At(x, y))
		}
	}

	if sdl.Init(sdl.INIT_EVERYTHING) != nil {
		return &result, sdl.GetError()
	}

	result.window, err = sdl.CreateWindow(TITLEDISPLAY,
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WIODISPLAY_W, WIODISPLAY_H,
		sdl.WINDOW_SHOWN)
	if err != nil {
		return &result, fmt.Errorf("Failed to create window: %s\n", err)
	}

	result.window.SetResizable(true)

	result.renderer, err = sdl.CreateRenderer(result.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return &result, fmt.Errorf("Failed to create renderer: %s\n", err)
	}

	var errBackgroundTexture error
	result.backgroundFrame.texture, errBackgroundTexture = result.renderer.CreateTextureFromSurface(bgSurface)
	if errBackgroundTexture != nil {
		return nil, fmt.Errorf("Failed to create texture: %s\n", errBackgroundTexture)
	}
	result.backgroundFrame.sourceArea = bgSurface.Bounds()

	//p.imgbuf = image.NewRGBA(image.Rect(0, 0, WIODISPLAY_W, WIODISPLAY_H))
	var errsurf error
	result.imgsurf, errsurf = sdl.CreateRGBSurface(0, WIODISPLAY_W, WIODISPLAY_H, 32, 0, 0, 0, 0)
	return &result, errsurf
}

func (p *Sdlwio) Close() error {
	e := p.window.Destroy()
	if e != nil {
		return e
	}

	return p.renderer.Destroy()
}

func (p *Sdlwio) Init(rotation Rotation) error {

	return nil
}

func rotateRect(r *sdl.Rect) {
	v := r.X
	r.X = r.Y
	r.Y = v

	v = r.W
	r.W = r.H
	r.H = v

}

func (p *Sdlwio) updateview() error {
	texture, err := p.renderer.CreateTextureFromSurface(p.imgsurf)
	if err != nil {
		return fmt.Errorf("Failed to create texture: %s\n", err)
	}
	defer texture.Destroy()
	winW, winH, winSizeErr := p.renderer.GetOutputSize()
	if winSizeErr != nil {
		return fmt.Errorf("updateview failed on GetOutputSize %s", winSizeErr)
	}

	//Background
	ratio := float64(p.backgroundFrame.sourceArea.Dx()) / float64(p.backgroundFrame.sourceArea.Dy())

	bgSrc := sdl.Rect{0, 0, int32(p.backgroundFrame.sourceArea.Dx()), int32(p.backgroundFrame.sourceArea.Dy())}

	targetHeight := min(winH, int32(float64(winW)/ratio)) //Take all.. or limit width
	targetWidth := int32(float64(targetHeight) * ratio)
	bgDst := sdl.Rect{0, 0, targetWidth, targetHeight} //WIODISPLAY_W, WIODISPLAY_H}

	//If fullscreen then

	viewSrc := sdl.Rect{0, 0, WIODISPLAY_W, WIODISPLAY_H}
	scale := float64(targetHeight) / float64(p.backgroundFrame.sourceArea.Dy())
	viewDst := sdl.Rect{int32(BACKGROUNDPNG_VIEW_X0 * scale),
		int32(BACKGROUNDPNG_VIEW_Y0 * scale),
		int32(scale * (BACKGROUNDPNG_VIEW_X1 - BACKGROUNDPNG_VIEW_X0)),
		int32(scale * (BACKGROUNDPNG_VIEW_Y1 - BACKGROUNDPNG_VIEW_Y0))}
	//Scale

	if !p.Landscape {
		scale = 1 / scale
		viewDst = sdl.Rect{int32(BACKGROUNDPNG_VIEW_Y0 * scale),
			int32(BACKGROUNDPNG_VIEW_X0 * scale),
			int32(scale * (BACKGROUNDPNG_VIEW_Y1 - BACKGROUNDPNG_VIEW_Y0)),
			int32(scale * (BACKGROUNDPNG_VIEW_X1 - BACKGROUNDPNG_VIEW_X0))}

	}

	//

	p.renderer.Clear()
	//p.renderer.Copy(p.backgroundFrame.texture, &bgSrc, &bgDst)
	//p.renderer.Copy(texture, &viewSrc, &viewDst)

	if p.Landscape {
		if p.Flipped {
			p.renderer.CopyEx(p.backgroundFrame.texture, &bgSrc, &bgDst, 180, nil, sdl.FLIP_NONE)
			viewDst.X = bgDst.W - viewDst.W - viewDst.X
			viewDst.Y = bgDst.H - viewDst.H - viewDst.Y
			p.renderer.CopyEx(texture, &viewSrc, &viewDst, 180, nil, sdl.FLIP_NONE)
		} else {
			p.renderer.CopyEx(p.backgroundFrame.texture, &bgSrc, &bgDst, 0, nil, sdl.FLIP_NONE)
			p.renderer.CopyEx(texture, &viewSrc, &viewDst, 0, nil, sdl.FLIP_NONE)
		}
	} else {
		p.renderer.CopyEx(p.backgroundFrame.texture, &bgSrc, &bgDst, 90, nil, sdl.FLIP_NONE)
		viewDst.X = 0
		viewDst.Y = 0
		a := viewDst.W
		viewDst.W = viewDst.H
		viewDst.H = a
		//viewDst.H = bgDst.W
		p.renderer.CopyEx(texture, &viewSrc, &viewDst, 90, nil, sdl.FLIP_NONE)
	}

	p.renderer.Present()
	return nil
}

func (p *Sdlwio) Backlight(on bool) error {
	p.backlightNow = on

	return nil
}
func (p *Sdlwio) SetWindow(x int16, y int16, w int16, h int16) error {
	p.winNow_x = x
	p.winNow_y = y

	p.winNow_w = w
	p.winNow_h = h
	return nil
}

func (p *Sdlwio) StartWrite() error {
	p.winIndex = 0
	if p.writeStarted {
		return fmt.Errorf("already drawing")
	}
	p.writeStarted = true
	return nil
}

func (p *Sdlwio) EndWrite() error {
	p.writeStarted = false
	return p.updateview()
}
func (p *Sdlwio) Write8bit(arr []byte) error { //8bit mode?
	//fmt.Printf("arr koko %v\n", len(arr))
	return p.Write8bit(arr)
}

func (p *Sdlwio) Write16bitbytes(arr []byte) error {
	arr16 := make([]uint16, len(arr)/2)
	for i := range arr16 {
		arr16[i] = uint16(arr[i*2])<<8 | uint16(arr[i*2+1])
	}
	return p.Write16bit(arr16)
}
func (p *Sdlwio) Write16bit(arr []uint16) error {
	//fmt.Printf("arr16=0x%X\n", arr)
	//fmt.Printf("WINDOW x:%v y:%v w:%v h:%v\n", p.winNow_x, p.winNow_y, p.winNow_w, p.winNow_h)
	for i, v := range arr {
		r := (uint32(byte(v>>11)) * 255) / uint32(0x1F)
		g := (uint32((v>>5)&0x3F) * 255) / uint32(0x3F)
		b := (uint32(v&0x1F) * 255) / uint32(0x1F)

		pos := i + int(p.winIndex)

		x := int(p.winNow_x) + pos%int(p.winNow_w)
		y := int(p.winNow_y) + pos/int(p.winNow_w)
		//fmt.Printf("i=%v x=%v y=%v v=0x%X (R:0x%X G:0x%X B:0x%X)\n", i, x, y, v, r, g, b)
		//fmt.Printf("pos=%v x=%v y=%v\n", pos, x, y)
		p.imgsurf.Set(
			x, y,
			color.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: 255})
	}
	p.winIndex += len(arr)
	return nil
}

/*
func RGBATo565(c color.RGBA) uint16 {
	r, g, b, _ := c.RGBA()
	return uint16((r & 0xF800) +
		((g & 0xFC00) >> 5) +
		((b & 0xF800) >> 11))
}
*/

func (p *Sdlwio) Sleep(sleeping bool) error {
	return nil
}

func (p *Sdlwio) SetRotation(rotation Rotation) error {
	//TODO SELECT CASE ERI SUUNTAAN... ERITAVOIN TÄYTTÄÄ.  Vs mitenpäin pidetään näyttöä
	return nil
}

func (p *Sdlwio) SetScrollArea(topFixedArea int16, bottomFixedArea int16) error {
	return nil
}

func (p *Sdlwio) SetScroll(line int16) {

}

func (p *Sdlwio) GetWioKeys() uint32 {

	sdl.PollEvent()

	keystate := sdl.GetKeyboardState()
	//fmt.Printf("STATE=%#v\n", keystate)
	result := uint32(0)

	if 0 < keystate[sdl.SCANCODE_UP] {
		result |= KEYMASK_UP
	}
	if 0 < keystate[sdl.SCANCODE_DOWN] {
		result |= KEYMASK_DOWN
	}
	if 0 < keystate[sdl.SCANCODE_LEFT] {
		result |= KEYMASK_LEFT
	}
	if 0 < keystate[sdl.SCANCODE_RIGHT] {
		result |= KEYMASK_RIGHT
	}
	if 0 < keystate[sdl.SCANCODE_RETURN] { //Push down joystick
		result |= KEYMASK_CENTER
	}

	if 0 < keystate[sdl.SCANCODE_A] { //Buttons on top
		result |= KEYMASK_A
	}
	if 0 < keystate[sdl.SCANCODE_S] { //Buttons on top
		result |= KEYMASK_B
	}
	if 0 < keystate[sdl.SCANCODE_D] { //Buttons on top
		result |= KEYMASK_C
	}
	return result
}

/*
func main() {
	sdlwio, errInit := InitSdlwio()
	if errInit != nil {
		fmt.Printf("err init %s\n", errInit)
		return
	}

	sdlwio.renderer.SetDrawColor(0, 0, 0, 255)
	sdlwio.renderer.Clear()

	sdlwio.renderer.SetDrawColor(255, 255, 255, 255)
	sdlwio.renderer.DrawPoint(150, 300)

	sdlwio.renderer.SetDrawColor(0, 0, 255, 255)
	sdlwio.renderer.DrawLine(0, 0, 200, 200)

	sdlwio.renderer.Present()
	sdl.Delay(16)

	for {
		event := sdl.PollEvent()
		if event != nil {
			fmt.Printf("event=%#v\n", event)
		}

		switch event.(type) {
		case *sdl.QuitEvent:
			sdlwio.Close()
			return
		}
	}

}
*/
