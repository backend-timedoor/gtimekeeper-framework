package database

import (
	"context"

	"github.com/backend-timedoor/gtimekeeper-framework/app"
	"github.com/backend-timedoor/gtimekeeper-framework/base/contracts"
	"github.com/backend-timedoor/gtimekeeper-framework/base/database/drivers"
	"github.com/backend-timedoor/gtimekeeper-framework/utils/app/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)


func BootDatabase() *database.Database {
	driver := GetDatabaseDriver()
	
	g, _ := gorm.Open(driver.GetGormDialect(), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})

	m, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(app.Config.GetString("database.mongo")))

	return &database.Database{
		DB: g,
		Mongo: m,
	}
}

func GetDatabaseDriver() contracts.DatabaseDriver {
	config := app.Config

	switch config.Get("database.connection") {
	case "mysql":
		return &drivers.MysqlDriver{}
	case "pgsql":
		return &drivers.PgsqlDriver{}
	default:
		panic("Database driver not found use mysql or pgsql")
	}
}