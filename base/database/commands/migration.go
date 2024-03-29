package commands

import (
	"log"

	"github.com/backend-timedoor/gtimekeeper-framework/base/database"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	DefaultTimeFormat = "20060102150405"
)

func GetMigration() *migrate.Migrate {
	db := database.DBDriverAnchor
	driver, err := db.GetDriver()

	if err != nil {
		log.Fatal("failed to create migration instance:", err)
	}

	migration, _ := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		db.GetConnection(),
		driver,
	)

	return migration
}
