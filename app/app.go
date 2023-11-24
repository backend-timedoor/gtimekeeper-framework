package app

import (
	"github.com/backend-timedoor/gtimekeeper/base/contracts"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)


var (
	Config contracts.Config
	DB     *gorm.DB
	Log    *logrus.Logger
	Server   contracts.Server
	Cache  contracts.Cache
	Queue contracts.Queue
	Schedule contracts.Schedule
)


