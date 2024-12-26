package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"

	"github.com/Vikuuu/gobase"
)

func handlerHelp(cmds *commands, s *state, cmd command) error {
	fmt.Println("Welcome to Gobase! An automated migration tool")
	fmt.Println("Usage: gobase <command>")
	fmt.Println("Available commands: ")

	for name, info := range cmds.registeredCommands {
		fmt.Printf("%-10s %-30s %s\n", name, info.usage, info.description)
	}

	return nil
}

func handlerInit(s *state, cmd command) error {
	dbConn, err := gobase.SqliteConn(gobase.DBFILENAME)
	if err != nil {
		return err
	}

	err = gobase.CreateMetaTable(dbConn)
	if err != nil {
		return err
	}
	return nil
}

func handlerMigrate(s *state, cmd command) error {
	schemaFileName := s.cfg.SchemaData
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}
	migrationDir := filepath.Join(workingDir, s.cfg.Migration)

	dbConn, err := gobase.SqliteConn(gobase.DBFILENAME)
	if err != nil {
		return err
	}

	if len(cmd.Args) > 0 && cmd.Args[0] == "down" {
		var files []string
		err := filepath.Walk(migrationDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, info.Name())
			}
			return nil
		})
		if err != nil {
			return err
		}

		sort.Slice(files, func(i, j int) bool {
			// Extract the numerical part of the filenames
			re := regexp.MustCompile(`^(\d+)`)
			numI, _ := strconv.Atoi(re.FindString(files[i]))
			numJ, _ := strconv.Atoi(re.FindString(files[j]))

			// Compare the numbers in decreasing order
			return numI > numJ
		})

		fileName := files[0]

		err = gobase.DownMigrate(dbConn, filepath.Join(migrationDir, fileName))
		if err != nil {
			return err
		}

		return nil
	}

	latestVer, err := gobase.GetLatestId(dbConn)
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("%03d_migration.sql", latestVer+1)

	newJSON, changeJSON, err := gobase.MigrationFile(dbConn, schemaFileName, migrationDir, fileName)
	if err != nil {
		return err
	}

	err = gobase.UpMigrate(dbConn, filepath.Join(migrationDir, fileName), newJSON, changeJSON)
	if err != nil {
		return err
	}

	return nil
}
