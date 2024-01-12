package drivers

import (
	"database/sql"
	"fmt"

	"github.com/backend-timedoor/gtimekeeper/app"
	"github.com/golang-migrate/migrate/v4/database"
	my "github.com/golang-migrate/migrate/v4/database/postgres"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlDriver struct{}

func (d *MysqlDriver) GetConnection() string {
	return "mysql"
}

func (d *MysqlDriver) GetSqlDb() *sql.DB {
	db, _ := sql.Open(d.GetConnection(), d.GetDsn())

	return db
}

func (d *MysqlDriver) GetDriver() (database.Driver, error) {
	return my.WithInstance(d.GetSqlDb(), &my.Config{})
}

func (d *MysqlDriver) GetGormDialect() gorm.Dialector {
	return mysql.Open(d.GetDsn())
}

func (d *MysqlDriver) GetDsn() string {
	config := app.Config
	pgsql := config.Get("database.mysql").(map[string]any)

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		pgsql["host"],
		pgsql["username"],
		pgsql["password"],
		pgsql["database"],
		pgsql["port"],
	)
}