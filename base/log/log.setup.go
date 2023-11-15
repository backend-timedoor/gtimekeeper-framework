package log

import (
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func BootLog() *logrus.Logger {
	log := logrus.New()

	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcname, filename
		},
	})

	log.AddHook(&DebugStackHook{})

	return log
}