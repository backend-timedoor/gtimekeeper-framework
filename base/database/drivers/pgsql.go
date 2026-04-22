package drivers

import (
	"database/sql"
	"fmt"
	"strings"

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
	SslMode  string
	Options  string
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
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		d.Host,
		d.Username,
		d.Password,
		d.Database,
		d.Port,
		d.getSSLMode(),
	)

	if d.shouldSendOptions() {
		dsn += fmt.Sprintf(" options='%s'", d.Options)
	}

	return dsn
}

func (d *PgsqlDriver) getSSLMode() string {
	if d.SslMode == "" {
		return "disable"
	}

	return d.SslMode
}

func (d *PgsqlDriver) shouldSendOptions() bool {
	if strings.TrimSpace(d.Options) == "" {
		return false
	}

	// RDS Proxy rejects PostgreSQL startup command-line options.
	return !strings.Contains(d.Host, ".proxy-")
}
