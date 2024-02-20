package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

const EXTENTION = ".sql"

type MigrationCreateCommand struct{}

func (m *MigrationCreateCommand) Signature() string {
	return "migrate:create"
}

func (m *MigrationCreateCommand) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "name", Aliases: []string{"n"}},
	}
}

func (m *MigrationCreateCommand) Handle(c *cli.Context) error {
	name := "create_table"

	if  c.NArg() >= 1 {
		name = c.Args().First()
	} else if c.String("name") != "" {
		name = c.String("name")
	}

	// fmt.Println(c.Args(), c.Args().Slice(), c.NArg())

	createCmd(name)

	return nil
}

func createCmd(name string) error {
	var err error
	dir := filepath.Clean("database/migrations")

	// create dir
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	for _, direction := range []string{"up", "down"} {
		basename := generateName(name, direction)
		filename := filepath.Join(dir, basename)

		if err = createFile(direction, filename, name); err != nil {
			return err
		}

		// absPath, _ := filepath.Abs(filename)
		log.Println(basename)
	}

	return nil
}

func generateName(name string, direction string) string {
	format := DefaultTimeFormat
	startTime := time.Now()
	version := startTime.Format(format)
	ext := "." + strings.TrimPrefix(EXTENTION, ".")

	return fmt.Sprintf("%s_%s.%s%s", version, name, direction, ext)
}

func createFile(direction string, filename string, nameArg string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}

	tableName := extractTableName(nameArg)

	if extractAction(nameArg) == "create" {
		if direction == "up" {
			upStatement(f, tableName)
		} else {
			downStatement(f, tableName)
		}
	}

	return f.Close()
}

func upStatement(file *os.File, table string) {
	statement := fmt.Sprintf(`CREATE TABLE %s (
	id SERIAL PRIMARY KEY,
	created_at TIMESTAMPTZ,
	updated_at TIMESTAMPTZ
);`, table)

	_, err := file.WriteString(statement)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func downStatement(file *os.File, table string) {
	statement := fmt.Sprintf("DROP TABLE IF EXISTS %s;", table)

	_, err := file.WriteString(statement)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func extractTableName(filename string) string {
	parts := strings.Split(filename, "_")

	if len(parts) > 1 && parts[0] == "create" && parts[1] == "table" {
		return strings.Join(parts[2:], "_")
	}

	return "table"
}

func extractAction(input string) string {
	parts := strings.Split(input, "_")
	if len(parts) > 0 {
		return parts[0]
	}

	return ""
}