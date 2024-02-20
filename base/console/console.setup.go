package console

import (
	"os"

	"github.com/backend-timedoor/gtimekeeper-framework/base/contracts"
	cmd "github.com/backend-timedoor/gtimekeeper-framework/base/database/commands"
	"github.com/urfave/cli/v2"
)

func New(commands []contracts.Commands) {
	c := &cli.App{}

	c.Name = "GTime Keeper Project"
	c.UsageText = "g-time-keeper [global options] command [options] [arguments...]"

	commands = append(commands, []contracts.Commands{
		&cmd.MigrationCreateCommand{},
		&cmd.MigrationUpCommand{},
		&cmd.MigrationDownCommand{},
	}...)

	for _, command := range commands {
		command := command

		cliCommand := cli.Command{
			Name:  command.Signature(),
			Flags: command.Flags(),
			Action: func(c *cli.Context) error {
				return command.Handle(c)
			},
		}

		c.Commands = append(c.Commands, &cliCommand)
	}

	Run(c)
}

func Run(c *cli.App) {
	args := os.Args

	if len(args) >= 2 {
		if args[1] == "gtime" {
			cliArgs := append([]string{args[0]}, args[2:]...)

			if err := c.Run(cliArgs); err != nil {
				panic(err.Error())
			}

			os.Exit(0)
		}
	}
}
