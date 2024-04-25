package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate"

	_ "github.com/golang-migrate/migrate/database/sqlite3"

	_ "github.com/golang-migrate/migrate/source/file"
)

func main() {
	var storagePath, migratorPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "storage path")
	flag.StringVar(&migratorPath, "migrations-path", "", "migrator path")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "migrations table")
	flag.Parse()

	if storagePath == "" || migratorPath == "" {
		panic("storage-path or migrator-path are required")
	}

	m, err := migrate.New("file://"+migratorPath, fmt.Sprintf("sqlite3://%s?x-migrations_table=%s", storagePath, migrationsTable))

	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}

	fmt.Println("migrations applied")

}
