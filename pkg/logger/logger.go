package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

func InitLogger(level string) *logrus.Logger {
	log := logrus.New()
	logLvl, err := logrus.ParseLevel(level)
	if err != nil {
		logLvl = logrus.InfoLevel
	}

	log.SetLevel(logLvl)
	log.SetOutput(os.Stdout)

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return log
}
