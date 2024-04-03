//go:build wioterminal

/*
Ambient meter demo

# Read sensors and also do some simulation

Displays pressure, relative humidity and temperature on screen..

maybe some trend?
*/

package main

import (
	"fmt"
	"machine"
	"time"

	"github.com/hjkoskel/wiointerface"

	"github.com/hjkoskel/BME280golib"
)

func HandleTermintingError(err error) {
	for err != nil {
		fmt.Printf("FAIL %s\n", err)
		time.Sleep(time.Second * 3)
	}
}

func main() {
	fmt.Printf("* Ambient sensor *\n")
	disp := &wiointerface.WioDisplayHW{}
	initErr := disp.Init(wiointerface.Rotation0)

	HandleTermintingError(initErr)
	if initErr != nil {
		fmt.Printf("ERROR:%s\n", initErr)
	}

	i2cConnect := machine.I2C1
	errI2CConfigure := i2cConnect.Configure(machine.I2CConfig{
		SCL: machine.SCL1_PIN,
		SDA: machine.SDA1_PIN,
	})

	HandleTermintingError(errI2CConfigure)

	a := BME280golib.CreateI2CTiny(i2cConnect, BME280golib.BME280DEVICEBIT0)
	ap := &a

	bmeDevice, errInit := BME280golib.CreateBME280I2C(ap)
	if errInit != nil {
		HandleTermintingError(fmt.Errorf("Error initializing BME280 %v\n", errInit.Error()))
	}

	measdev := &bmeDevice
	errRuntime := RunAmbientMeter(disp, measdev)

	HandleTermintingError(errRuntime)

}
