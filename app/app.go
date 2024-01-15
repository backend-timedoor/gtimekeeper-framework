package app

import (
	"github.com/backend-timedoor/gtimekeeper-framework/base/contracts"
	"github.com/backend-timedoor/gtimekeeper-framework/utils/app/database"
	"github.com/sirupsen/logrus"
)


var (
	Config contracts.Config
	DB     *database.Database
	Log    *logrus.Logger
	Server   contracts.Server
	Cache  contracts.Cache
	Queue contracts.Queue
	Schedule contracts.Schedule
	Kafka contracts.Kafka
	Mail contracts.Mail
)


