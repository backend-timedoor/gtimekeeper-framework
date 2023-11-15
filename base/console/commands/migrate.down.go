package commands

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/urfave/cli/v2"
)

type MigrationDownCommand struct{}

func (m *MigrationDownCommand) Signature() string {
	return "migrate:down"
}

func (m *MigrationDownCommand) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{Name: "step", Aliases: []string{"s"}},
	}
}

func (m *MigrationDownCommand) Handle(c *cli.Context) (err error) {
	step := c.Int("step")
	migration := GetMigration()

	if step > 0 {
		s := -(step)
		err = migration.Steps(s)
	} else {
		err = migration.Down()
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("failted to run migration:", err)
	}

	fmt.Println("database migrated successfully")

	return nil
}
