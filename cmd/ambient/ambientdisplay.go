package main

import (
	"fmt"
	"time"

	"github.com/hjkoskel/wiointerface"

	"github.com/hjkoskel/BME280golib"
	"github.com/hjkoskel/gomonochromebitmap"
)

var bigFont = gomonochromebitmap.GetFont_11x16()

type AmbientDisplay struct {
	Latest     BME280golib.HumTempPressureMeas
	TextBitmap gomonochromebitmap.MonoBitmap
}

func (p *AmbientDisplay) AddData(tNow time.Time, amb BME280golib.HumTempPressureMeas) {
	p.Latest = amb
}

func (p *AmbientDisplay) RenderToDisplay(w wiointerface.WioInterface) error {
	//CLEAR
	w.SetWindow(0, 0, 320, 240)
	w.StartWrite()
	for i := 0; i < 320*240; i++ {
		w.Write16bit([]uint16{0}) //TODO OPTIMIZE!!!
	}
	w.EndWrite()

	p.TextBitmap.Print(fmt.Sprintf("T:%.2f C\nH:%.2f %%\nP:%.3f kPa", p.Latest.Temperature, p.Latest.Rh, p.Latest.Pressure/1000),
		bigFont, 16, 1, p.TextBitmap.Bounds(), true, true, false, true)
	w.SetWindow(0, 0, 320, 240)
	w.StartWrite()
	for y := 0; y < 200; y++ {
		for x := 0; x < 320; x++ {
			if p.TextBitmap.GetPix(x, y) {
				w.Write16bit([]uint16{0xFFFF}) //TODO.... optimoidumpi rutiini. Kerralla enemmän pisteitä
			} else {
				w.Write16bit([]uint16{0})
			}
		}
	}
	w.EndWrite()

	return nil
}

var initialConfig = BME280golib.BME280Config{
	Oversample_humidity:    BME280golib.OVRSAMPLE_1,
	Oversample_pressure:    BME280golib.OVRSAMPLE_1,
	Oversample_temperature: BME280golib.OVRSAMPLE_1,
	Mode:                   BME280golib.MODE_NORMAL,
	Standby:                BME280golib.STANDBYDURATION_500,
	Filter:                 BME280golib.FILTER_NO}

func RunAmbientMeter(w wiointerface.WioInterface, bmeDevice BME280golib.BME280Device) error {
	resetErr := bmeDevice.SoftReset()
	if resetErr != nil {
		return fmt.Errorf("Failed on reset %v", resetErr.Error())
	}
	errConf := bmeDevice.Configure(initialConfig)
	if errConf != nil {
		return errConf
	}

	view := AmbientDisplay{TextBitmap: gomonochromebitmap.NewMonoBitmap(320, 240, false)}

	for {
		measResult, errMeasResult := bmeDevice.Read()
		if errMeasResult != nil {
			return errMeasResult
		}
		view.AddData(time.Now(), measResult)
		view.RenderToDisplay(w)
		time.Sleep(time.Millisecond * 1500)
	}

}
