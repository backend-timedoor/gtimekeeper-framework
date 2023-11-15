package commands

import (
	"fmt"
	"log"

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
	step := c.Int("step")
	migration := GetMigration()

	if step > 0 {
		err = migration.Steps(step)
	} else {
		err = migration.Up()
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("failted to run migration:", err)
	}

	fmt.Println("database migrated successfully")

	return nil
}
