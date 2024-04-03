//go:build wioterminal

/*
Hw based implementatio of interface

Using tags...
*/

package wiointerface

import (
	"machine"
	"time"

	"tinygo.org/x/drivers"
)

var cmdBuf [6]byte

type Config struct {
	Width            int16
	Height           int16
	Rotation         Rotation
	DisplayInversion bool
}

type driver interface {
	configure(config *Config)
	write8(b byte)
	write8n(b byte, n int)
	write8sl(b []byte)
	write16(data uint16)
	write16n(data uint16, n int)
	write16sl(data []uint16)
	write16sl_bytes(data []byte)
}

type WioDisplayHW struct {
	width    int16
	height   int16
	rotation Rotation
	driver   driver

	x0, x1 int16 // cached address window; prevents useless/expensive
	y0, y1 int16 // syscalls to PASET and CASET

	dc  machine.Pin
	cs  machine.Pin
	rst machine.Pin
	rd  machine.Pin

	backlight machine.Pin

	//Buttons TODO PISTÄ MUUALLE!

	/*
		btnA machine.Pin
		btnB machine.Pin
		btnC machine.Pin

		btnLeft   machine.Pin
		btnRight  machine.Pin
		btnUp     machine.Pin
		btnDown   machine.Pin
		btnCenter machine.Pin
	*/
}

func delay(m int) { //TODO time.sleep?

	t := time.Now().UnixNano() + int64(time.Duration(m*1000)*time.Microsecond)
	for time.Now().UnixNano() < t {
	}
}

var initCmd = []byte{
	0xEF, 3, 0x03, 0x80, 0x02,
	0xCF, 3, 0x00, 0xC1, 0x30,
	0xED, 4, 0x64, 0x03, 0x12, 0x81,
	0xE8, 3, 0x85, 0x00, 0x78,
	0xCB, 5, 0x39, 0x2C, 0x00, 0x34, 0x02,
	0xF7, 1, 0x20,
	0xEA, 2, 0x00, 0x00,
	PWCTR1, 1, 0x23, // Power control VRH[5:0]
	PWCTR2, 1, 0x10, // Power control SAP[2:0];BT[3:0]
	VMCTR1, 2, 0x3e, 0x28, // VCM control
	VMCTR2, 1, 0x86, // VCM control2
	MADCTL, 1, 0x48, // Memory Access Control
	VSCRSADD, 1, 0x00, // Vertical scroll zero
	PIXFMT, 1, 0x55,
	FRMCTR1, 2, 0x00, 0x18,
	DFUNCTR, 3, 0x08, 0x82, 0x27, // Display Function Control
	0xF2, 1, 0x00, // 3Gamma Function Disable
	GAMMASET, 1, 0x01, // Gamma curve selected
	GMCTRP1, 15, 0x0F, 0x31, 0x2B, 0x0C, 0x0E, 0x08, // Set Gamma
	0x4E, 0xF1, 0x37, 0x07, 0x10, 0x03, 0x0E, 0x09, 0x00,
	GMCTRN1, 15, 0x00, 0x0E, 0x14, 0x03, 0x11, 0x07, // Set Gamma
	0x31, 0xC1, 0x48, 0x08, 0x0F, 0x0C, 0x31, 0x36, 0x0F,
}

var buf [64]byte

type spiDriver struct {
	bus drivers.SPI
}

func (pd *spiDriver) configure(config *Config) {
}

func (pd *spiDriver) write8(b byte) {
	buf[0] = b
	pd.bus.Tx(buf[:1], nil)
}

func (pd *spiDriver) write8n(b byte, n int) {
	buf[0] = b
	for i := 0; i < n; i++ {
		pd.bus.Tx(buf[:1], nil)
	}
}

func (pd *spiDriver) write8sl(b []byte) {
	pd.bus.Tx(b, nil)
}

func (pd *spiDriver) write16(data uint16) {
	buf[0] = uint8(data >> 8)
	buf[1] = uint8(data)
	pd.bus.Tx(buf[:2], nil)
}

func (pd *spiDriver) write16n(data uint16, n int) {
	for i := 0; i < len(buf); i += 2 {
		buf[i] = uint8(data >> 8)
		buf[i+1] = uint8(data)
	}

	for i := 0; i < (n >> 5); i++ {
		pd.bus.Tx(buf[:], nil)
	}

	pd.bus.Tx(buf[:n%64], nil)
}

func (pd *spiDriver) write16sl(data []uint16) {
	for i, c := 0, len(data); i < c; i++ {
		buf[0] = uint8(data[i] >> 8)
		buf[1] = uint8(data[i])
		pd.bus.Tx(buf[:2], nil)
	}
}

func (pd *spiDriver) write16sl_bytes(data []byte) {
	for i := 0; i < len(data)/2; i++ {
		pd.bus.Tx(data[i*2:(2+i*2)], nil)
	}
}

func (p *WioDisplayHW) Init(rotation Rotation) error {
	machine.SPI3.Configure(machine.SPIConfig{
		SCK:       machine.LCD_SCK_PIN,
		SDO:       machine.LCD_SDO_PIN,
		SDI:       machine.LCD_SDI_PIN,
		Frequency: 40000000,
	})

	p.dc = machine.LCD_DC
	p.cs = machine.LCD_SS_PIN
	p.rst = machine.LCD_RESET
	p.rd = machine.NoPin

	p.driver = &spiDriver{bus: machine.SPI3}

	/*
		WIO_5S_PRESS
		WIO_5S_DOWN
		WIO_5S_RIGHT
		WIO_5S_LEFT
	*/

	machine.WIO_KEY_A.Configure(machine.PinConfig{Mode: machine.PinInput})
	machine.WIO_KEY_B.Configure(machine.PinConfig{Mode: machine.PinInput})
	machine.WIO_KEY_C.Configure(machine.PinConfig{Mode: machine.PinInput})

	machine.WIO_5S_LEFT.Configure(machine.PinConfig{Mode: machine.PinInput})
	machine.WIO_5S_RIGHT.Configure(machine.PinConfig{Mode: machine.PinInput})
	machine.WIO_5S_UP.Configure(machine.PinConfig{Mode: machine.PinInput})
	machine.WIO_5S_DOWN.Configure(machine.PinConfig{Mode: machine.PinInput})
	machine.WIO_5S_PRESS.Configure(machine.PinConfig{Mode: machine.PinInput})

	//-----------------

	config := Config{
		Width:            320,
		Height:           240,
		Rotation:         Rotation270, //TODO FIKSAA!!!
		DisplayInversion: false,
	}

	p.width = config.Width
	p.height = config.Height
	p.rotation = config.Rotation

	// try to pick an initial cache miss for one of the points
	p.x0, p.x1 = -(p.width + 1), p.x0
	p.y0, p.y1 = -(p.height + 1), p.y0

	output := machine.PinConfig{machine.PinOutput}

	// configure chip select if there is one
	p.cs.Configure(output)
	p.cs.High() // deselect

	p.dc.Configure(output)
	p.dc.High() // data mode

	// driver-specific configuration
	p.driver.configure(&config)

	if p.rd != machine.NoPin {
		p.rd.Configure(output)
		p.rd.High()
	}

	// configure hardware reset if there is one
	p.rst.Configure(output)
	p.rst.High()
	delay(100)
	p.rst.Low()
	delay(100)
	p.rst.High()
	delay(200)

	if config.DisplayInversion {
		initCmd = append(initCmd, INVON, 0x80)
	}

	initCmd = append(initCmd,
		SLPOUT, 0x80, // Exit Sleep
		DISPON, 0x80, // Display on
		0x00, // End of list
	)
	for i, c := 0, len(initCmd); i < c; {
		cmd := initCmd[i]
		if cmd == 0x00 {
			break
		}
		x := initCmd[i+1]
		numArgs := int(x & 0x7F)
		p.sendCommand(cmd, initCmd[i+2:i+2+numArgs])
		if x&0x80 > 0 {
			delay(150)
		}
		i += numArgs + 2
	}

	p.SetRotation(p.rotation)

	p.backlight = machine.LCD_BACKLIGHT
	p.backlight.Configure(machine.PinConfig{machine.PinOutput})

	p.backlight.High()

	return nil
}

func (d *WioDisplayHW) SetRotation(rotation Rotation) error {
	madctl := uint8(0)
	switch rotation % 8 {
	case Rotation0:
		madctl = MADCTL_MX | MADCTL_BGR
	case Rotation90:
		madctl = MADCTL_MV | MADCTL_BGR
	case Rotation180:
		madctl = MADCTL_MY | MADCTL_BGR | MADCTL_ML
	case Rotation270:
		madctl = MADCTL_MX | MADCTL_MY | MADCTL_MV | MADCTL_BGR | MADCTL_ML
	case Rotation0Mirror:
		madctl = MADCTL_BGR
	case Rotation90Mirror:
		madctl = MADCTL_MY | MADCTL_MV | MADCTL_BGR | MADCTL_ML
	case Rotation180Mirror:
		madctl = MADCTL_MX | MADCTL_MY | MADCTL_BGR | MADCTL_ML
	case Rotation270Mirror:
		madctl = MADCTL_MX | MADCTL_MY | MADCTL_MV | MADCTL_BGR | MADCTL_ML
	}
	cmdBuf[0] = madctl
	d.sendCommand(MADCTL, cmdBuf[:1])
	d.rotation = rotation
	return nil
}

func (p *WioDisplayHW) Backlight(on bool) error {
	if on {
		p.backlight.High()
	}
	p.backlight.Low()

	return nil
}

func (p *WioDisplayHW) SetWindow(x int16, y int16, w int16, h int16) error {
	//x += d.columnOffset
	//y += d.rowOffset
	x1 := x + w - 1
	if x != p.x0 || x1 != p.x1 {
		cmdBuf[0] = uint8(x >> 8)
		cmdBuf[1] = uint8(x)
		cmdBuf[2] = uint8(x1 >> 8)
		cmdBuf[3] = uint8(x1)
		p.sendCommand(CASET, cmdBuf[:4])
		p.x0, p.x1 = x, x1
	}
	y1 := y + h - 1
	if y != p.y0 || y1 != p.y1 {
		cmdBuf[0] = uint8(y >> 8)
		cmdBuf[1] = uint8(y)
		cmdBuf[2] = uint8(y1 >> 8)
		cmdBuf[3] = uint8(y1)
		p.sendCommand(PASET, cmdBuf[:4])
		p.y0, p.y1 = y, y1
	}
	p.sendCommand(RAMWR, nil)
	return nil
}

//go:inline
func (p *WioDisplayHW) StartWrite() error {
	p.cs.Low()
	return nil
}

//go:inline
func (p *WioDisplayHW) EndWrite() error {
	p.cs.High()
	return nil

}
func (p *WioDisplayHW) Write8bit(arr []byte) error {
	p.driver.write8sl(arr)
	//p.driver.write8sl(arr)
	//TODO MIKSI KIRJOITUS PITÄÄ TEHDÄ NIINKUIN TEHDÄÄN

	/*arr16 := make([]uint16, len(arr)/2)
	for i := range arr16 {
		arr16[i] = uint16(arr[i*2+1])<<8 | uint16(arr[i*2+0])
	}*/

	//return p.Write16bit(arr16)

	return nil

}
func (p *WioDisplayHW) Write16bit(arr []uint16) error {
	p.driver.write16sl(arr)
	return nil
}

func (p *WioDisplayHW) Write16bitbytes(arr []byte) error {
	p.driver.write16sl_bytes(arr)
	return nil

}

func (p *WioDisplayHW) sendCommand(cmd byte, data []byte) {
	p.StartWrite()
	p.dc.Low()
	p.driver.write8(cmd)
	p.dc.High()
	if data != nil {
		p.driver.write8sl(data)
	}
	p.EndWrite()
}

func (p *WioDisplayHW) Sleep(enableSleep bool) error {
	if enableSleep {
		// Shut down LCD panel.
		p.sendCommand(SLPIN, nil)
		time.Sleep(5 * time.Millisecond) // 5ms required by the datasheet
	} else {
		// Turn the LCD panel back on.
		p.sendCommand(SLPOUT, nil)
		// Note: the ili9341 documentation says that it is needed to wait at
		// least 120ms before going to sleep again. Sleeping here would not be
		// practical (delays turning on the screen too much), so just hope the
		// screen won't need to sleep again for at least 120ms.
		// In practice, it's unlikely the user will set the display to sleep
		// again within 120ms.
	}
	return nil

}

func (p *WioDisplayHW) SetScrollArea(topFixedArea int16, bottomFixedArea int16) error {
	if p.height < 320 {
		// The screen doesn't use the full 320 pixel height.
		// Enlarge the bottom fixed area to fill the 320 pixel height, so that
		// bottomFixedArea starts from the visible bottom of the screen.
		bottomFixedArea += 320 - p.height
	}
	cmdBuf[0] = uint8(topFixedArea >> 8)
	cmdBuf[1] = uint8(topFixedArea)
	cmdBuf[2] = uint8((320 - topFixedArea - bottomFixedArea) >> 8)
	cmdBuf[3] = uint8(320 - topFixedArea - bottomFixedArea)
	cmdBuf[4] = uint8(bottomFixedArea >> 8)
	cmdBuf[5] = uint8(bottomFixedArea)
	p.sendCommand(VSCRDEF, cmdBuf[:6])

	return nil
}

// SetScroll sets the vertical scroll address of the display.
func (p *WioDisplayHW) SetScroll(line int16) {
	cmdBuf[0] = uint8(line >> 8)
	cmdBuf[1] = uint8(line)
	p.sendCommand(VSCRSADD, cmdBuf[:2])
}

func (p *WioDisplayHW) GetWioKeys() uint32 {
	result := uint32(0)

	if !machine.WIO_5S_UP.Get() {
		result |= KEYMASK_UP
	}
	if !machine.WIO_5S_DOWN.Get() {
		result |= KEYMASK_DOWN
	}
	if !machine.WIO_5S_LEFT.Get() {
		result |= KEYMASK_LEFT
	}
	if !machine.WIO_5S_RIGHT.Get() {
		result |= KEYMASK_RIGHT
	}
	if !machine.WIO_5S_PRESS.Get() { //Push down joystick
		result |= KEYMASK_CENTER
	}

	if !machine.WIO_KEY_A.Get() { //Buttons on top
		result |= KEYMASK_A
	}
	if !machine.WIO_KEY_B.Get() { //Buttons on top
		result |= KEYMASK_B
	}
	if !machine.WIO_KEY_C.Get() { //Buttons on top
		result |= KEYMASK_C
	}
	return result
}
