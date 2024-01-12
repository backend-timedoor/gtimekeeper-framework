package app

import (
	"github.com/backend-timedoor/gtimekeeper/base/contracts"
	"github.com/backend-timedoor/gtimekeeper/utils/database"
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
)


