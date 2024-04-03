//go:build !wioterminal

package main

import (
	"fmt"

	"github.com/hjkoskel/wiointerface"
)

func main() {
	disp, errInit := wiointerface.InitSdlwio()
	if errInit != nil {
		fmt.Printf("err init %s\n", errInit)
		return
	}

	//disp.Flipped = true
	disp.Landscape = true

	spy := InitChipSpy() //PC can rune spy.. helps checking keys and colors

	var game *GameSettings = &game_blinky
	runner, errRunnerInit := NewChipRunner(*game.ChipSettings)

	if errRunnerInit != nil {
		fmt.Printf("runner init fail %s\n", errRunnerInit)
		return
	}

	chp := InitChip8(game.Code, game.ColorGuide)

	runner.Clear(disp)
	for {
		runErr := runner.Run(&chp, disp)
		if runErr != nil {
			fmt.Printf("run error %s\n", runErr)
			return
		}

		//if false {
		keysChanged := spy.SpyKeys(&chp)
		spritesChanged := spy.SpySprites(&chp)
		if keysChanged || spritesChanged {
			fmt.Printf("\n%s\n", spy)
		}
		//}

		//fmt.Printf("pc=0x%04X\n", chp.PC)
	}

}
