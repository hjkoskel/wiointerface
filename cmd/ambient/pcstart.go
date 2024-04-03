//go:build !wioterminal

package main

import (
	"fmt"
	"math/rand"

	"github.com/hjkoskel/wiointerface"

	"github.com/hjkoskel/BME280golib"
)

type FakeAmbientSensor struct {
}

func (p *FakeAmbientSensor) Close() error {
	return nil
}
func (p *FakeAmbientSensor) Configure(config BME280golib.BME280Config) error {
	return nil
}
func (p *FakeAmbientSensor) Read() (BME280golib.HumTempPressureMeas, error) {
	return BME280golib.HumTempPressureMeas{
		Temperature: 23 + rand.Float64()*5,
		Rh:          33 + rand.Float64()*0.4,
		Pressure:    100000 + rand.Float64()*1000,
	}, nil
}
func (p *FakeAmbientSensor) SoftReset() error {
	return nil
}
func (p *FakeAmbientSensor) GetCalibration() (BME280golib.CalibrationRegs, error) { //Gets latest values
	return BME280golib.CalibrationRegs{}, nil
}

func main() {
	disp, errInit := wiointerface.InitSdlwio()
	if errInit != nil {
		fmt.Printf("err init %s\n", errInit)
		return
	}
	disp.Landscape = true
	fakesensor := &FakeAmbientSensor{}
	errRuntime := RunAmbientMeter(disp, fakesensor)
	if errRuntime != nil {
		fmt.Printf("%s\n", errRuntime)
		return
	}

}
