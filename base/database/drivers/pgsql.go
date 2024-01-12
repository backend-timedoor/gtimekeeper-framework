package drivers

import (
	"database/sql"
	"fmt"

	"github.com/backend-timedoor/gtimekeeper-framework/app"
	"github.com/golang-migrate/migrate/v4/database"
	pg "github.com/golang-migrate/migrate/v4/database/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgsqlDriver struct{}

func (d *PgsqlDriver) GetConnection() string {
	return "postgres"
}

func (d *PgsqlDriver) GetSqlDb() *sql.DB {
	db, _ := sql.Open(d.GetConnection(), d.GetDsn())

	return db
}

func (d *PgsqlDriver) GetDriver() (database.Driver, error) {
	return pg.WithInstance(d.GetSqlDb(), &pg.Config{})
}

func (d *PgsqlDriver) GetGormDialect() gorm.Dialector {
	return postgres.Open(d.GetDsn())
}

func (d *PgsqlDriver) GetDsn() string {
	config := app.Config
	pgsql := config.Get("database.pgsql").(map[string]any)

	// "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		pgsql["host"],
		pgsql["username"],
		pgsql["password"],
		pgsql["database"],
		pgsql["port"],
	)
}