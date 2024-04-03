//go:build wioterminal

package main

import (
	"fmt"
	"os"

	"github.com/hjkoskel/wiointerface"
)

func main() {

	//Init what happens on wio terminal side
	disp := &wiointerface.WioDisplayHW{}

	//disp := InitDisplay()
	fmt.Printf("GOING TO INIT\n")
	initErr := disp.Init(wiointerface.Rotation0)
	if initErr != nil {
		fmt.Printf("ERROR:%s\n", initErr)
	}

	errDemorun := MarioFill(disp)
	if errDemorun != nil {
		fmt.Printf("Running demo failed %s\n", errDemorun)
		os.Exit(-1)
	}

}
