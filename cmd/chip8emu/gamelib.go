/*
gamelib
This hardcodes games in this program
*/
package main

import (
	_ "embed"
	"image/color"

	"github.com/hjkoskel/wiointerface"
)

type GameSettings struct {
	ChipSettings *ChipRunSettings
	ColorGuide   map[uint16]byte
	Code         []byte
}

/*****
BLINKY
*****/
//go:embed roms/BLINKY
var gamecode_blinky []byte

var game_blinky = GameSettings{
	ChipSettings: &ChipRunSettings{
		BorderColor: wiointerface.RGBATo565(color.RGBA{R: 64, G: 64, B: 64}),
		Keymask: []uint32{
			0,
			wiointerface.KEYMASK_B,    //1
			0,                         //2
			wiointerface.KEYMASK_UP,   //3
			0,                         //4
			0,                         //5
			wiointerface.KEYMASK_DOWN, //6
			wiointerface.KEYMASK_LEFT,
			wiointerface.KEYMASK_RIGHT,
			0,                      //9
			0,                      //10
			0,                      //11
			0,                      //12
			0,                      //13
			0,                      //14
			wiointerface.KEYMASK_A, //15,
		},
		Palette: []uint16{
			0,
			wiointerface.RGBATo565(color.RGBA{R: 0, G: 0xFF, B: 0}),
			wiointerface.RGBATo565(color.RGBA{R: 0xFF, G: 0xFF, B: 0}),
			wiointerface.RGBATo565(color.RGBA{R: 0, G: 0xFF, B: 0xFF}),
			wiointerface.RGBATo565(color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF}),
		},
	},
	ColorGuide: map[uint16]byte{
		0x909: 2,
		0x913: 2,
		0x90D: 2,
		0x903: 2,
		0x917: 2,
		0x915: 2,
		0x905: 2,
		0x911: 2,
		0x907: 2,
		0x8E5: 2,
		0x901: 2,
		0x8ED: 2,
		0x8F9: 2,
	},
	Code: gamecode_blinky,
}

/*
TETRIS,
The only good thing what came from russia
*/
//go:embed roms/TETRIS
var gamecode_tetris []byte

var game_tetris = GameSettings{
	ChipSettings: &ChipRunSettings{
		BorderColor: wiointerface.RGBATo565(color.RGBA{R: 64, G: 64, B: 64}),
		Keymask: []uint32{
			0,
			0,                          //1
			0,                          //2
			0,                          //3
			wiointerface.KEYMASK_UP,    //4
			wiointerface.KEYMASK_LEFT,  //5
			wiointerface.KEYMASK_RIGHT, //6
			wiointerface.KEYMASK_DOWN,
			0,
			0, //9
			0, //10
			0, //11
			0, //12
			0, //13
			0, //14
			0, //15,
		},
		Palette: []uint16{
			0,
			wiointerface.RGBATo565(color.RGBA{R: 0, G: 0xFF, B: 0}),
			wiointerface.RGBATo565(color.RGBA{R: 0xFF, G: 0xFF, B: 0}),
			wiointerface.RGBATo565(color.RGBA{R: 0, G: 0xFF, B: 0xFF}),
			wiointerface.RGBATo565(color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF}),
		},
	},
	ColorGuide: map[uint16]byte{
		0x30C: 3,
		0x310: 3,
		0x330: 2,
		0x324: 2,
		0x32C: 2,
	},
	Code: gamecode_tetris,
}

/*
MISSILE
*/
//go:embed roms/MISSILE
var gamecode_missile []byte

var game_missile = GameSettings{
	ChipSettings: &ChipRunSettings{
		BorderColor: wiointerface.RGBATo565(color.RGBA{R: 64, G: 64, B: 64}),
		Keymask: []uint32{
			0,
			0,                          //1
			0,                          //2
			0,                          //3
			wiointerface.KEYMASK_UP,    //4
			wiointerface.KEYMASK_LEFT,  //5
			wiointerface.KEYMASK_RIGHT, //6
			wiointerface.KEYMASK_DOWN,
			0,
			0, //9
			0, //10
			0, //11
			0, //12
			0, //13
			0, //14
			0, //15,
		},
		Palette: []uint16{
			0,
			wiointerface.RGBATo565(color.RGBA{R: 0, G: 0xFF, B: 0}),
		},
	},
	ColorGuide: map[uint16]byte{},
	Code:       gamecode_missile,
}

/*
TANK
*/
//go:embed roms/TANK
var gamecode_tank []byte

var game_tank = GameSettings{
	ChipSettings: &ChipRunSettings{
		BorderColor: wiointerface.RGBATo565(color.RGBA{R: 64, G: 64, B: 64}),
		Keymask: []uint32{
			0,
			0,                          //1
			0,                          //2
			0,                          //3
			wiointerface.KEYMASK_UP,    //4
			wiointerface.KEYMASK_LEFT,  //5
			wiointerface.KEYMASK_RIGHT, //6
			wiointerface.KEYMASK_DOWN,
			0,
			0, //9
			0, //10
			0, //11
			0, //12
			0, //13
			0, //14
			0, //15,
		},
		Palette: []uint16{
			0,
			wiointerface.RGBATo565(color.RGBA{R: 0, G: 0xFF, B: 0}),
		},
	},
	ColorGuide: map[uint16]byte{},
	Code:       gamecode_tank,
}

//go:embed roms/PONG
var gamecode_pong []byte

var game_pong = GameSettings{
	ChipSettings: &ChipRunSettings{
		BorderColor: wiointerface.RGBATo565(color.RGBA{R: 64, G: 64, B: 64}),
		Keymask: []uint32{
			wiointerface.KEYMASK_UP,
			wiointerface.KEYMASK_UP,    //1
			wiointerface.KEYMASK_UP,    //2
			wiointerface.KEYMASK_UP,    //3
			wiointerface.KEYMASK_UP,    //4
			wiointerface.KEYMASK_LEFT,  //5
			wiointerface.KEYMASK_RIGHT, //6
			wiointerface.KEYMASK_DOWN,
			wiointerface.KEYMASK_UP,
			wiointerface.KEYMASK_UP, //9
			wiointerface.KEYMASK_UP, //10
			wiointerface.KEYMASK_UP, //11
			wiointerface.KEYMASK_UP, //12
			wiointerface.KEYMASK_UP, //13
			wiointerface.KEYMASK_UP, //14
			wiointerface.KEYMASK_UP, //15,
		},
		Palette: []uint16{
			0,
			wiointerface.RGBATo565(color.RGBA{R: 0xFF, G: 0xFF, B: 0}),
		},
	},
	ColorGuide: map[uint16]byte{},
	Code:       gamecode_pong,
}

//go:embed roms/sokoban.ch8
var gamecode_sokoban []byte

var game_sokoban = GameSettings{
	ChipSettings: &ChipRunSettings{
		BorderColor: wiointerface.RGBATo565(color.RGBA{R: 64, G: 64, B: 64}),
		Keymask: []uint32{
			wiointerface.KEYMASK_B,     //0
			0,                          //1
			0,                          //2
			0,                          //3
			0,                          //4
			wiointerface.KEYMASK_UP,    //5
			0,                          //6
			wiointerface.KEYMASK_LEFT,  //7
			wiointerface.KEYMASK_DOWN,  //8
			wiointerface.KEYMASK_RIGHT, //9
			wiointerface.KEYMASK_A,     //10
			0,                          //11
			0,                          //12
			0,                          //13
			0,                          //14
			0,                          //15,
		},
		Palette: []uint16{
			0,
			wiointerface.RGBATo565(color.RGBA{R: 0xFF, G: 0xFF, B: 0}),
		},
	},
	ColorGuide: map[uint16]byte{},
	Code:       gamecode_sokoban,
}
