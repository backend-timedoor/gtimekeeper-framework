package database

import (
	"database/sql"
	"time"

	bm "github.com/backend-timedoor/gtimekeeper-framework/base/database/mongo"
	"github.com/backend-timedoor/gtimekeeper-framework/base/database/redis"
	"github.com/backend-timedoor/gtimekeeper-framework/container"
	db "github.com/golang-migrate/migrate/v4/database"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

const ContainerName string = "db"

var DBDriverAnchor DatabaseDriver

type Database struct {
	*gorm.DB
	Mongo  *mongo.Client
	Redis  *redis.Redis
	Config *Config
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
	ReadDriver DatabaseDriver
	GormConfig *gorm.Config
	Pool       *PoolConfig
	Mongo      string
	Redis      *redis.Config
}

type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
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

	db.DB, err = gorm.Open(config.Driver.GetGormDialect(), cloneGormConfig(config.GormConfig))
	if err != nil {
		l.Infof("Error connecting to database: %s", err.Error())
	}

	resolvedPool := resolvePoolConfig(config.Pool)
	if err = applyPoolConfig(db.DB, resolvedPool); err != nil {
		l.Infof("Error configuring database pool: %s", err.Error())
	}

	if config.ReadDriver != nil && config.ReadDriver.GetDsn() != config.Driver.GetDsn() {
		resolver := dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{
				config.ReadDriver.GetGormDialect(),
			},
			Policy: dbresolver.RandomPolicy{},
		}).
			SetMaxIdleConns(resolvedPool.MaxIdleConns).
			SetMaxOpenConns(resolvedPool.MaxOpenConns).
			SetConnMaxLifetime(resolvedPool.ConnMaxLifetime).
			SetConnMaxIdleTime(resolvedPool.ConnMaxIdleTime)

		err = db.DB.Use(resolver)
		if err != nil {
			l.Infof("Error configuring read/write database resolver: %s", err.Error())
		}
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

func cloneGormConfig(config *gorm.Config) *gorm.Config {
	if config == nil {
		return &gorm.Config{}
	}

	c := *config
	return &c
}

func resolvePoolConfig(pool *PoolConfig) *PoolConfig {
	resolved := &PoolConfig{
		MaxOpenConns:    40,
		MaxIdleConns:    20,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
	}

	if pool == nil {
		return resolved
	}

	if pool.MaxOpenConns > 0 {
		resolved.MaxOpenConns = pool.MaxOpenConns
	}
	if pool.MaxIdleConns >= 0 {
		resolved.MaxIdleConns = pool.MaxIdleConns
	}
	if pool.ConnMaxLifetime > 0 {
		resolved.ConnMaxLifetime = pool.ConnMaxLifetime
	}
	if pool.ConnMaxIdleTime > 0 {
		resolved.ConnMaxIdleTime = pool.ConnMaxIdleTime
	}

	if resolved.MaxIdleConns > resolved.MaxOpenConns {
		resolved.MaxIdleConns = resolved.MaxOpenConns
	}

	return resolved
}

func applyPoolConfig(gormDB *gorm.DB, pool *PoolConfig) error {
	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(pool.MaxOpenConns)
	sqlDB.SetMaxIdleConns(pool.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(pool.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(pool.ConnMaxIdleTime)

	return nil
}

//func (d *Database) () *Database {
//    db, _ := container.Get(ContainerName)
//    return db.(*Database)
//}
