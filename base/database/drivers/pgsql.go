package drivers

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4/database"
	pg "github.com/golang-migrate/migrate/v4/database/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgsqlDriver struct {
	Host     string
	Username string
	Password string
	Database string
	Port     int
	Config   pg.Config
}

func (d *PgsqlDriver) GetConnection() string {
	return "postgres"
}

func (d *PgsqlDriver) GetSqlDb() *sql.DB {
	db, _ := sql.Open(d.GetConnection(), d.GetDsn())

	return db
}

func (d *PgsqlDriver) GetDriver() (database.Driver, error) {
	return pg.WithInstance(d.GetSqlDb(), &d.Config)
}

func (d *PgsqlDriver) GetGormDialect() gorm.Dialector {
	return postgres.Open(d.GetDsn())
}

func (d *PgsqlDriver) GetDsn() string {
	// "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		d.Host,
		d.Username,
		d.Password,
		d.Database,
		d.Port,
	)
}
