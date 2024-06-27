package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/backend-timedoor/gtimekeeper-framework/utils/helper"
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
	start := time.Now()
	name := "create_table"

	if c.NArg() >= 1 {
		name = c.Args().First()
	} else if c.String("name") != "" {
		name = c.String("name")
	}

	createCmd(name)

	duration := time.Since(start).Seconds() * 1000

	fmt.Print(DotMessage("Migration created", duration))

	return nil
}

func createCmd(name string) error {
	var err error
	baseDir := filepath.Clean("database/migrations")

	// Create the base directory if it doesn't exist
	if err = os.MkdirAll(baseDir, os.ModePerm); err != nil {
		return err
	}

	// Define migration directions
	directions := []string{"up", "down"}

	// Loop through the directions and create directories for each
	for _, direction := range directions {
		basename := generateName(name, direction)
		directionPath := filepath.Join(baseDir, direction)
		dirPath := filepath.Join(directionPath, basename)

		// Create the direction-specific directory
		if err = os.MkdirAll(directionPath, os.ModePerm); err != nil {
			return err
		}

		if err = createFile(direction, dirPath, name); err != nil {
			return err
		}

		// Print the name of the created file
		message := fmt.Sprintf("[%s] %s created successfully", dirPath, direction)
		log.Println(message)
	}

	return nil
}

func generateName(name string, direction string) string {
	format := DefaultTimeFormat
	startTime := time.Now()
	version := startTime.Format(format)
	ext := strings.TrimPrefix(EXTENTION, ".")

	return fmt.Sprintf("%s_%s.%s.%s", version, name, direction, ext)
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

	// Ensure there are at least three parts: "create", "<TableName>", "table"
	if len(parts) >= 3 && parts[0] == "create" && parts[len(parts)-1] == "table" {
		tableName := strings.Join(parts[1:len(parts)-1], "_")
		return helper.Pluralize(tableName)
	}

	return "tables"
}

func extractAction(input string) string {
	parts := strings.Split(input, "_")
	if len(parts) > 0 {
		return parts[0]
	}

	return ""
}
