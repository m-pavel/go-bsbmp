package main

import (
	"flag"
	_ "net/http"
	_ "net/http/pprof"

	"github.com/d2r2/go-bsbmp"
	"github.com/d2r2/go-i2c"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/m-pavel/go-hassio-mqtt/pkg"
)

type BsBmpService struct {
	addr *int
	line *int
	i2c  *i2c.I2C
	bmp  *bsbmp.BMP
}
type Request struct {
	Temperature float32 `json:"temperature"`
	Pressure    float32 `json:"pressure"`
	Altitude    float32 `json:"altitude"`
}

func (ts BsBmpService) PrepareCommandLineParams() {
	flag.Int("i2c-addr", 0, "I2C address")
	flag.Int("i2c-line", 1, "I2C line")
}
func (ts BsBmpService) Name() string { return "bsbmp" }

func (ts *BsBmpService) Init(client MQTT.Client, topic, topicc, topica string, debug bool, ss ghm.SendState) error {
	var err error
	if ts.i2c, err = i2c.NewI2C(0x76, 1); err != nil {
		return err
	}
	ts.bmp, err = bsbmp.NewBMP(bsbmp.BMP280, ts.i2c)
	return err
}

func (ts BsBmpService) Do() (interface{}, error) {
	req := Request{}
	var err error
	if req.Temperature, err = ts.bmp.ReadTemperatureC(bsbmp.ACCURACY_STANDARD); err != nil {
		return nil, err
	}
	if req.Pressure, err = ts.bmp.ReadPressurePa(bsbmp.ACCURACY_STANDARD); err != nil {
		return nil, err
	}
	if req.Altitude, err = ts.bmp.ReadAltitude(bsbmp.ACCURACY_STANDARD); err != nil {
		return nil, err
	}
	return &req, err
}

func (ts BsBmpService) Close() error {
	return ts.i2c.Close()
}

func main() {
	ghm.NewStub(&BsBmpService{}).Main()
}
