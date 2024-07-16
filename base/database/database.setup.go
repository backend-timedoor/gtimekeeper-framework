package database

import (
	"database/sql"

	bm "github.com/backend-timedoor/gtimekeeper-framework/base/database/mongo"
	"github.com/backend-timedoor/gtimekeeper-framework/base/database/redis"
	"github.com/backend-timedoor/gtimekeeper-framework/container"
	db "github.com/golang-migrate/migrate/v4/database"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

const ContainerName string = "db"

var DBDriverAnchor DatabaseDriver

type Database struct {
	*gorm.DB
	WithContext *gorm.DB
	Mongo       *mongo.Client
	Redis       *redis.Redis
	Config      *Config
}

type DatabaseDriver interface {
	GetConnection() string
	GetSqlDb() *sql.DB
	GetDriver() (db.Driver, error)
	GetGormDialect() gorm.Dialector
	GetDsn() string
}

type Config struct {
	Driver     DatabaseDriver
	GormConfig *gorm.Config
	Mongo      string
	Redis      *redis.Config
}

func New(config *Config) *Database {
	DBDriverAnchor = config.Driver
	var (
		err error
		db  = &Database{
			Config: config,
		}
		l = container.Log()
	)

	db.DB, err = gorm.Open(config.Driver.GetGormDialect(), config.GormConfig)
	if err != nil {
		l.Infof("Error connecting to database: %s", err.Error())
	}

	if config.Mongo != "" {
		db.Mongo, err = bm.New(config.Mongo)
		if err != nil {
			l.Infof("Error connecting to MongoDB: %s", err.Error())
		}
	}

	if config.Redis != nil {
		db.Redis, err = redis.New(config.Redis)
		if err != nil {
			l.Infof("Error connecting to Redis: %s", err.Error())
		}

	}

	container.Set(ContainerName, db)

	return db
}

//func (d *Database) () *Database {
//    db, _ := container.Get(ContainerName)
//    return db.(*Database)
//}
