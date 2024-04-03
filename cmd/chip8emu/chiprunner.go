/*
Chiprunner renders.. calls chip virtual engine
and handles keyboard  and keyboard&video settings

And sounds?
*/
package main

import (
	"fmt"
	"time"

	"github.com/hjkoskel/wiointerface"
)

type ChipRunSettings struct {
	BorderColor uint16
	Keymask     []uint32
	Palette     []uint16
}

type ChipRunner struct {
	Settings        ChipRunSettings
	InSuperchipMode bool
}

func NewChipRunner(settings ChipRunSettings) (ChipRunner, error) {
	if len(settings.Keymask) == 0 {
		return ChipRunner{}, fmt.Errorf("no keys defined")
	}
	if len(settings.Palette) < 2 {
		return ChipRunner{}, fmt.Errorf("at least 2 colors needed to be defined")
	}
	return ChipRunner{Settings: settings}, nil
}

func (p *ChipRunner) Clear(disp wiointerface.WioInterface) {
	disp.SetWindow(0, 0, 320, 240)
	disp.StartWrite()
	for i := 0; i < 320*240; i++ {
		disp.Write16bit([]uint16{p.Settings.BorderColor}) //TODO OPTIMIZE!!!
	}
	disp.EndWrite()
}

func (p *ChipRunner) Render(chp *Chip8, disp wiointerface.WioInterface) error {
	if p.InSuperchipMode != chp.HiRes {
		p.Clear(disp) //Mode changing..redraw
	}
	p.InSuperchipMode = chp.HiRes

	if !p.InSuperchipMode {
		disp.SetWindow(0, 0, 64*4, 32*4)
		disp.StartWrite()
		for y := 0; y < 32; y++ {
			for rep := 0; rep < 4; rep++ {
				for x := 0; x < 64; x++ {
					c := p.Settings.Palette[chp.Videomem[x][y]]
					disp.Write16bit([]uint16{c, c, c, c})

				}
			}
		}
		disp.EndWrite()
	} else {
		//fmt.Printf("SUPERCHIP!\n")
		disp.SetWindow(0, 0, 128*2, 64*2)
		disp.StartWrite()
		for y := 0; y < 64; y++ {
			for rep := 0; rep < 2; rep++ {
				for x := 0; x < 128; x++ {
					c := p.Settings.Palette[chp.Videomem[x][y]]
					disp.Write16bit([]uint16{c, c})
				}
			}
		}
		disp.EndWrite()
	}
	return nil
}

func (p *ChipRunner) Run(chp *Chip8, disp wiointerface.WioInterface) error {
	timeRunned := time.Now()

	chp.UpdateTimer(time.Now())
	doDisplayUpdate, errExec := chp.ExecOp()
	if errExec != nil {
		return errExec
	}
	if doDisplayUpdate {
		p.Render(chp, disp)
	}

	for time.Since(timeRunned) < time.Millisecond*2 { //500Hz loop
		keys := disp.GetWioKeys()
		for i, mask := range p.Settings.Keymask {
			chp.Keyboard[i] = 0 < keys&mask
		}
	}
	return nil
}
