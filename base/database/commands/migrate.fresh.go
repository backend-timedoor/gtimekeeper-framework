package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/urfave/cli/v2"
)

type MigrationFreshCommand struct{}

func (m *MigrationFreshCommand) Signature() string {
	return "migrate:fresh"
}

func (m *MigrationFreshCommand) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (m *MigrationFreshCommand) Handle(c *cli.Context) (err error) {
	start := time.Now()
	migration := GetMigration("up")

	if err := migration.Drop(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failted to run migration:", err)
	}

	duration := time.Since(start).Seconds() * 1000

	fmt.Print(DotMessage("Drop all tables successfully", duration))

	return nil
}
