/*
Thin wrapper for wio terminal.

Instead of having custom made machine package etc... have some interfaces
For some reason machine package is in included in devices. Very non testable situation

# Interfaces also allow more high level abstraction instead of playng with device registers

# Actual implementation on hardware is simpler than what comes from ili9341
*/
package wiointerface

import "image/color"

type Rotation uint8

// Clockwise rotation of the screen.
const (
	Rotation0 = iota
	Rotation90
	Rotation180
	Rotation270
	Rotation0Mirror
	Rotation90Mirror
	Rotation180Mirror
	Rotation270Mirror
)

const (
	KEYMASK_UP     uint32 = 1 << 0
	KEYMASK_DOWN   uint32 = 1 << 1
	KEYMASK_LEFT   uint32 = 1 << 2
	KEYMASK_RIGHT  uint32 = 1 << 3
	KEYMASK_CENTER uint32 = 1 << 4

	KEYMASK_A uint32 = 1 << 5
	KEYMASK_B uint32 = 1 << 6
	KEYMASK_C uint32 = 1 << 7
)

type WioInterface interface {
	Backlight(on bool) error

	//Display SPI
	Init(rotation Rotation) error
	SetWindow(x int16, y int16, w int16, h int16) error

	StartWrite() error
	EndWrite() error
	Write8bit(arr []byte) error
	Write16bitbytes(arr []byte) error
	Write16bit(arr []uint16) error

	Sleep(sleeping bool) error

	SetRotation(rotation Rotation) error

	SetScrollArea(topFixedArea int16, bottomFixedArea int16) error

	SetScroll(line int16)
	// Keys
	GetWioKeys() uint32
}

/*
What ILI9341 does


*/

func RGBATo565(c color.RGBA) uint16 {
	r, g, b, _ := c.RGBA()
	return uint16((r & 0xF800) +
		((g & 0xFC00) >> 5) +
		((b & 0xF800) >> 11))
}
