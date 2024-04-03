package main

import (
	"fmt"
	"time"

	"github.com/hjkoskel/wiointerface"

	_ "embed"
)

//go:embed a.pic
var demospritedata []byte

func MarioFill(disp wiointerface.WioInterface) error {
	//time.Sleep(time.Second * 3)
	rulla := int16(0)
	for {
		tStart := time.Now()
		for y := 0; y < 240-32; y += 32 {
			for x := 0; x < 320-16; x += 16 {
				disp.SetWindow(int16(x), int16(y), 16, 32)
				disp.StartWrite()
				disp.Write16bitbytes(demospritedata)
				disp.EndWrite()
			}
		}
		fmt.Printf("runtime %s\n", time.Since(tStart))
		disp.SetScroll(rulla)
		rulla = (rulla + 1) % 320

	}
	//time.Sleep(time.Second * 9999)

	return nil

}
