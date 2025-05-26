package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	ctx := context.Background()

	DSN := os.Getenv("DSN")
	if DSN == "" {
		fmt.Printf("Failed to get DSN")
		os.Exit(1)
	}

	pool, err := pgxpool.New(ctx, DSN)
	if err != nil {
		fmt.Printf("Failed to connect to db")
		os.Exit(1)
	}
	defer pool.Close()

	files, err := readSQLFiles("migrations")
	if err != nil {
		fmt.Printf("Failed to read migration files: %v\n", err)
		os.Exit(1)
	}

	for _, file := range files {
		fmt.Printf("Applying migration: %s\n", file)

		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Failed to read %s: %v\n", file, err)
			os.Exit(1)
		}

		_, err = pool.Exec(ctx, string(content))
		if err != nil {
			fmt.Printf("Failed to execute %s: %v\n", file, err)
			os.Exit(1)
		}
	}

	fmt.Printf("All migrations applied successfully")
}

func readSQLFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".sql") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(files)
	return files, nil
}
