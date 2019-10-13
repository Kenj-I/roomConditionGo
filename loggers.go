//refer https://re-engines.com/2018/11/05/go言語のエラーハンドリングとログローテーション/
package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	envProduction  = "production"
	envStaging     = "staging"
	envDevelopment = "development"
)

var log = logrus.New()

// logrusの初期設定
func logInit() {
	var environment string

	switch os.Getenv("GO_ENV") {
	case envProduction, envStaging, envDevelopment:
		environment = os.Getenv("GO_ENV")
	default:
		environment = envDevelopment
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})

	if environment == envProduction {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if environment == envProduction {
		// TODO: logの出力
		// productionはlogfileに保存してログローテートする
		ljack := &lumberjack.Logger{
			Filename:   "/home/pi/houseCondition/log/condition.log",
			MaxSize:    10,
			MaxAge:     15,
			MaxBackups: 5,
			LocalTime:  true,
		}
		logrus.SetOutput(ljack)
	} else {
		logrus.SetOutput(os.Stdout)
	}
}
