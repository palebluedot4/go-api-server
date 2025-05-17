package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
	logFile *os.File
}

var (
	instance *Logger
	once     sync.Once
)

type Config struct {
	Level string `mapstructure:"level"`
}

func Instance(level ...string) *Logger {
	once.Do(func() {
		l := logrus.New()

		l.SetLevel(logrus.InfoLevel)
		l.SetOutput(os.Stderr)
		l.SetReportCaller(true)

		l.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat:   time.RFC3339Nano,
			DisableHTMLEscape: true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				functionName := filepath.Base(f.Function)
				fileName := filepath.Base(f.File)
				return fmt.Sprintf("%s()", functionName), fmt.Sprintf("%s:%d", fileName, f.Line)
			},
		})

		if len(level) > 0 && level[0] != "" {
			parsedLevel, err := logrus.ParseLevel(strings.ToLower(level[0]))
			if err == nil {
				l.SetLevel(parsedLevel)
			}
		}

		// TODO: Add file logging

		instance = &Logger{
			Logger:  l,
			logFile: nil,
		}
	})
	return instance
}

func Shutdown() {
	if instance != nil && instance.logFile != nil {
		instance.logFile.Close()
		instance.logFile = nil
	}
}

func Configure(cfg Config) {
	if instance == nil {
		Instance()
	}

	level, err := logrus.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		level = logrus.InfoLevel
	}
	instance.SetLevel(level)
}
