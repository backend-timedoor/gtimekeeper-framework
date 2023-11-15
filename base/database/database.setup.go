package database

import (
	"github.com/backend-timedoor/gtimekeeper/app"
	"github.com/backend-timedoor/gtimekeeper/base/contracts"
	"github.com/backend-timedoor/gtimekeeper/base/database/drivers"
	"gorm.io/gorm"
)


func BootDatabase() *gorm.DB {
	driver := GetDatabaseDriver()
	
	g, _ := gorm.Open(driver.GetGormDialect(), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})

	return g
}

func GetDatabaseDriver() contracts.DatabaseDriver {
	config := app.Config

	switch config.Get("database.connection") {
	case "mysql":
		return &drivers.MysqlDriver{}
	case "pgsql":
		return &drivers.PgsqlDriver{}
	default:
		panic("Database driver not found")
	}
}