package contracts

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/hibiken/asynq"
	"github.com/segmentio/kafka-go"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

type ServiceProvider interface {
	Boot()
	Register()
}

type Config interface {
	Env(envName string, defaultValue ...any) any
	Add(name string, configuration any)
	Get(path string, defaultValue ...any) any
	GetString(path string, defaultValue ...any) string
	GetInt(path string, defaultValue ...any) int
	GetBool(path string, defaultValue ...any) bool
}

type Server interface {
	Grpc() ServerHandle
	Http() ServerHandle
	RegisterCustomValidation([]CustomeValidation)
}

type ServerHandle interface {
	Start()
	Run(string)
}

type Commands interface {
	Signature() string
	Flags() []cli.Flag
	Handle(*cli.Context) error
}

type DatabaseDriver interface {
	GetConnection() string
	GetSqlDb() *sql.DB
	GetDriver() (database.Driver, error)
	GetGormDialect() gorm.Dialector
	GetDsn() string
}

type Queue interface {
	Signature() string
	Options() []asynq.Option
	Handle(context.Context, *asynq.Task) error
}

type Schedule interface {
	Signature() string
	Options() []asynq.Option
	Schedule() string
	Handle(context.Context, *asynq.Task) error
}

type CustomeValidation interface {
	Signature() string
	Handle(validator.FieldLevel) bool
}

type Kafka interface {
	Produce(context.Context, ...kafka.Message) error
}
