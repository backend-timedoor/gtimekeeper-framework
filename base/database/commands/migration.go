package commands

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/backend-timedoor/gtimekeeper-framework/base/database"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"golang.org/x/term"
)

var (
	DefaultTimeFormat = "20060102150405"
)

func GetMigration(path string) *migrate.Migrate {
	db := database.DBDriverAnchor
	driver, err := db.GetDriver()

	if err != nil {
		log.Fatal("failed to create migration instance:", err)
	}

	migration, _ := migrate.NewWithDatabaseInstance(
		"file://database/migrations/"+path,
		db.GetConnection(),
		driver,
	)

	return migration
}

func DotMessage(message string, args ...any) string {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Failed to get terminal size:", err)
		return message
	}

	messageLen := len(message)

	if args != nil {
		messageLen = messageLen + len(fmt.Sprintf("(%.2fms)", args[0])) + 2
	}

	dots := width - messageLen

	return fmt.Sprintf("%s %s (%.2fms)\n", message, strings.Repeat(".", dots), args[0])
}
