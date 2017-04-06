package logger

import (
	log "github.com/Sirupsen/logrus"
	"os"
)

type ExtraFileds map[string]interface{}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)
}

func Debug(extra map[string]interface{}, format string, args ...interface{}) {
	log.WithFields(extra).Debugf(format, args...)
}

func Info(extra map[string]interface{}, format string, args ...interface{}) {
	log.WithFields(extra).Infof(format, args...)
}

func Error(extra map[string]interface{}, format string, args ...interface{}) {
	log.WithFields(extra).Errorf(format, args...)
}
