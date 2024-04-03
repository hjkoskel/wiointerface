//go:build !wioterminal

package main

import (
	"fmt"
	"os"

	"github.com/hjkoskel/wiointerface"
)

func main() {
	fmt.Printf("ajetaan PC:llä")
	sdlwio, errInit := wiointerface.InitSdlwio()
	if errInit != nil {
		fmt.Printf("err init %s\n", errInit)
		return
	}

	sdlwio.Landscape = true

	errDemorun := MarioFill(sdlwio)
	if errDemorun != nil {
		fmt.Printf("Running demo failed %s\n", errDemorun)
		os.Exit(-1)
	}

}
