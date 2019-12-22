package main

import (
	"context"
	"fmt"
	"time"

	logger "github.com/d2r2/go-logger"

	firebase "firebase.google.com/go"
	"github.com/d2r2/go-bsbmp"
	"github.com/d2r2/go-i2c"
	"google.golang.org/api/option"
)

// RoomCondition is
type RoomCondition struct {
	now         time.Time
	temperature string
	pressure    string
	humidity    string
}

var now time.Time

func main() {
	logger.ChangePackageLogLevel("i2c", logger.InfoLevel)
	logger.ChangePackageLogLevel("bsbmp", logger.InfoLevel)
	logInit()

	condition, err := fetchCondition()
	if err != nil {
		log.Error(err)
	}

	recordCondition(&condition)

}

func fetchCondition() (RoomCondition, error) {
	now = time.Now()
	log.Info(now)

	// Create new connection to i2c-bus on 1 line with address 0x76.
	// Use i2cdetect utility to find device address over the i2c-bus
	i2c, err := i2c.NewI2C(0x76, 1)
	if err != nil {
		log.Fatal(err)
	}
	defer i2c.Close()

	sensor, err := bsbmp.NewBMP(bsbmp.BME280, i2c) // signature=0x60
	if err != nil {
		log.Fatal(err)
	}

	_, err = sensor.ReadSensorID()
	if err != nil {
		log.Fatal(err)
	}

	err = sensor.IsValidCoefficients()
	if err != nil {
		log.Error(err)
	}

	// Read temperature in celsius degree
	t, err := sensor.ReadTemperatureC(bsbmp.ACCURACY_STANDARD)
	if err != nil {
		log.Error(err)
	}
	fmt.Printf("Temperature: %v *C\n", t)

	// Read atmospheric pressure in pascal
	p, err := sensor.ReadPressurePa(bsbmp.ACCURACY_LOW)
	if err != nil {
		log.Error(err)
	}
	hpa := fmt.Sprintf("%.0f", p/100)
	fmt.Printf("Pressure = %v hPa\n", hpa)

	// Read atmospheric pressure in mmHg
	supported, h1, err := sensor.ReadHumidityRH(bsbmp.ACCURACY_LOW)
	if supported {
		if err != nil {
			log.Error(err)
		}
		fmt.Printf("Humidity = %v %%\n", h1)
	}

	condition := RoomCondition{
		now:         now,
		temperature: fmt.Sprintf("%.1f", t-5),
		pressure:    hpa,
		humidity:    fmt.Sprintf("%.1f", h1+20),
	}

	return condition, nil
}

func recordCondition(condition *RoomCondition) {
	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsFile("./credential.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	// insert
	_, _, err = client.Collection("conditions").Add(ctx, map[string]interface{}{
		"createdAt":   condition.now,
		"temperature": condition.temperature,
		"pressure":    condition.pressure,
		"humidity":    condition.humidity,
	})
}
