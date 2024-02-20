package drivers

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4/database"
	my "github.com/golang-migrate/migrate/v4/database/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlDriver struct {
	Host     string
	Username string
	Password string
	Database string
	Port     int
}

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
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		d.Username,
		d.Password,
		d.Host,
		d.Port,
		d.Database,
	)
}
