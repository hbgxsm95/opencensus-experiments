// Copyright 2017, OpenCensus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Starter codes for the iot application under the OpenCensus framework

package main

import (
	"fmt"
	"github.com/d2r2/go-dht"
	"github.com/d2r2/go-logger"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"math"
	"time"
)

var lg = logger.NewPackageLogger("main",
	logger.DebugLevel,
	// logger.InfoLevel,
)

const (
	high int = 1
	low  int = 0
)

func main() {
	go readTemperature()
	board := raspi.NewAdaptor()
	ads1015 := i2c.NewADS1015Driver(board)
	soundSensor := aio.NewGroveSoundSensorDriver(ads1015, "0")
	lightSensor := aio.NewGroveLightSensorDriver(ads1015, "1")

	work := func() {
		gobot.Every(1*time.Second, func() {
			soundStrength, soundErr := readSound(soundSensor)
			lightStrength, lightErr := lightSensor.Read()
			if soundErr != nil || lightErr != nil {
				log.Fatalf("Could not read value from sound / light sensors\n")
			} else {
				fmt.Printf("Sound: %d, Light: %d\n", soundStrength, lightStrength)
			}
		})
	}

	robot := gobot.NewRobot("PinVoltageCollection",
		[]gobot.Connection{board},
		[]gobot.Device{ads1015, soundSensor, lightSensor},
		work,
	)

	robot.Start()

}

func readSound(sensor *aio.GroveSoundSensorDriver) (int, error) {
	min := math.MaxInt32
	max := math.MinInt32
	for i := 0; i < 100; i++ {
		strength, err := sensor.Read()
		if err != nil {
			log.Fatalf("Couldn't read data from the sensor\n")
		} else {
			if strength > max {
				max = strength
			}
			if strength < min {
				min = strength
			}
		}
	}
	return max - min, nil
}
func readTemperature() {
	for range time.Tick(5 * time.Second) {
		defer logger.FinalizeLogger()
		// Uncomment/comment next line to suppress/increase verbosity of output
		// logger.ChangePackageLogLevel("dht", logger.InfoLevel)

		sensorType := dht.DHT11
		// Read DHT11 sensor data from pin 4, retrying 10 times in case of failure.
		// You may enable "boost GPIO performance" parameter, if your device is old
		// as Raspberry PI 1 (this will require root privileges). You can switch off
		// "boost GPIO performance" parameter for old devices, but it may increase
		// retry attempts. Play with this parameter.
		temperature, humidity, retried, err :=
			dht.ReadDHTxxWithRetry(sensorType, 4, false, 10)
		if err != nil {
			lg.Fatal(err)
		}
		// print temperature and humidity
		lg.Infof("Sensor = %v: Temperature = %v*C, Humidity = %v%% (retried %d times)",
			sensorType, temperature, humidity, retried)
	}
}
