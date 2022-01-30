package service

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var (
	Mylog   *log.Logger
	InitLog sync.Once
)

func LogInit() {
	InitLog.Do(func() {
		LogConfiguration()
	})
}

func LogConfiguration() {
	log := logrus.New()

	log.Out = os.Stdout

	log.Formatter = &logrus.JSONFormatter{}

	Mylog = log
}
