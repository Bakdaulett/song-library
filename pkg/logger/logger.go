package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var log = logrus.New()

func InitLogger() {
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	logLevel := viper.GetString("log_level")
	if logLevel == "" {
		logLevel = "debug"
	}

	switch logLevel {
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	default:
		log.SetLevel(logrus.DebugLevel)
	}
}

func Info(msg string) {
	log.Info(msg)
}

func Debug(msg string) {
	log.Debug(msg)
}

func Error(msg string) {
	log.Error(msg)
}
