package commands

import (
	"fmt"
	"strconv"
	"time"

	"log"

	"github.com/golang-migrate/migrate/v4"

	"github.com/urfave/cli/v2"
)

type MigrationForceCommand struct{}

func (m *MigrationForceCommand) Signature() string {
	return "migrate:force"
}

func (m *MigrationForceCommand) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (m *MigrationForceCommand) Handle(c *cli.Context) (err error) {
	start := time.Now()
	migration := GetMigration("up")
	version := c.Args().First()

	// check version and convert to int
	if version == "" {
		log.Fatal("version is required")
	}

	// convert to int64
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		log.Fatal("version must be a number")
	}

	if err := migration.Force(versionInt); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failted to run migration:", err)
	}

	duration := time.Since(start)

	message := fmt.Sprintf("Version %v dirty reset successfully", version)

	fmt.Print(DotMessage(message, duration.Seconds()*1000))

	return nil
}
