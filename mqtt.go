package main

import (
	"flag"
	_ "net/http"
	_ "net/http/pprof"

	"fmt"
	"strconv"

	"log"

	"github.com/d2r2/go-bsbmp"
	"github.com/d2r2/go-i2c"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/m-pavel/go-hassio-mqtt/pkg"
)

type BsBmpService struct {
	addr  *uint
	line  *int
	i2c   *i2c.I2C
	bmp   *bsbmp.BMP
	debug bool
}
type Request struct {
	Temperature  float32 `json:"temperature"`
	PressurePa   float32 `json:"pressure-pa"`
	PressureMmHg float32 `json:"pressure-mmhg"`
	Altitude     float32 `json:"altitude"`
}

func (ts *BsBmpService) PrepareCommandLineParams() {
	ts.addr = flag.Uint("i2c-addr", 0, "I2C address (hex)")
	ts.line = flag.Int("i2c-line", 1, "I2C line")
}
func (ts BsBmpService) Name() string { return "bsbmp" }

func (ts *BsBmpService) Init(client MQTT.Client, topic, topicc, topica string, debug bool, ss ghm.SendState) error {
	var err error
	addr, err := strconv.ParseInt(fmt.Sprintf("%d", *ts.addr), 16, 64)
	if err != nil {
		return err
	}
	if ts.i2c, err = i2c.NewI2C(uint8(addr), *ts.line); err != nil {
		return err
	}
	ts.bmp, err = bsbmp.NewBMP(bsbmp.BMP180, ts.i2c)
	ts.debug = debug
	return err
}

func (ts BsBmpService) Do() (interface{}, error) {
	req := Request{}
	var err error
	if req.Temperature, err = ts.bmp.ReadTemperatureC(bsbmp.ACCURACY_STANDARD); err != nil {
		return nil, err
	}
	if req.PressurePa, err = ts.bmp.ReadPressurePa(bsbmp.ACCURACY_STANDARD); err != nil {
		return nil, err
	}
	if req.PressureMmHg, err = ts.bmp.ReadPressureMmHg(bsbmp.ACCURACY_STANDARD); err != nil {
		return nil, err
	}
	if req.Altitude, err = ts.bmp.ReadAltitude(bsbmp.ACCURACY_STANDARD); err != nil {
		return nil, err
	}
	if ts.debug {
		log.Println(req)
	}
	return &req, err
}

func (ts BsBmpService) Close() error {
	return ts.i2c.Close()
}

func main() {
	ghm.NewStub(&BsBmpService{}).Main()
}
