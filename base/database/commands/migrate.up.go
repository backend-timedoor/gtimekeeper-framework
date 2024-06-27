package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/urfave/cli/v2"
)

type MigrationUpCommand struct{}

func (m *MigrationUpCommand) Signature() string {
	return "migrate:up"
}

func (m *MigrationUpCommand) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (m *MigrationUpCommand) Handle(c *cli.Context) (err error) {
	start := time.Now()
	step := c.Int("step")
	migration := GetMigration("up")

	if step > 0 {
		err = migration.Steps(step)
	} else {
		err = migration.Up()
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("failted to run migration:", err)
	}

	duration := time.Since(start).Seconds() * 1000

	fmt.Print(DotMessage("Migration up successfully", duration))

	return nil
}
