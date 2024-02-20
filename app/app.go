package app

import (
	"github.com/backend-timedoor/gtimekeeper-framework/base/config"
	"github.com/backend-timedoor/gtimekeeper-framework/base/contracts"
	"github.com/backend-timedoor/gtimekeeper-framework/base/database"
	"github.com/backend-timedoor/gtimekeeper-framework/base/job"
	"github.com/backend-timedoor/gtimekeeper-framework/base/mail"
	"github.com/backend-timedoor/gtimekeeper-framework/base/validation"
	"github.com/sirupsen/logrus"
)

var (
	Config     *config.Config
	Validation *validation.Validation
	DB         *database.Database
	Log        *logrus.Logger
	Server     contracts.Server
	Job        *job.Job
	Kafka      contracts.Kafka
	Mail       *mail.Email
)
