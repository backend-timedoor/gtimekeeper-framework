package log

import (
	"github.com/backend-timedoor/gtimekeeper-framework/container"
	"github.com/sirupsen/logrus"
)

const ContainerName string = "log"

func New() *logrus.Logger {
	log := logrus.New()

	container.Set(ContainerName, log)

	return log
}
