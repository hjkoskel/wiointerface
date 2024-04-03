//go:build wioterminal

package main

import (
	_ "embed"
	"fmt"

	"github.com/hjkoskel/wiointerface"
)

func main() {
	disp := &wiointerface.WioDisplayHW{}
	initErr := disp.Init(wiointerface.Rotation0)
	if initErr != nil {
		fmt.Printf("ERROR:%s\n", initErr)
	}

	var game *GameSettings = &game_tetris
	runner, errRunnerInit := NewChipRunner(*game.ChipSettings)
	if errRunnerInit != nil {
		fmt.Printf("runner init fail %s\n", &errRunnerInit)
	}

	chp := InitChip8(game.Code, game.ColorGuide)
	runner.Clear(disp)
	for {
		runner.Run(&chp, disp)
	}
}
