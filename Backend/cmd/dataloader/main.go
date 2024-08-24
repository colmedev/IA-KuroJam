package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/colmedev/IA-KuroJam/Backend/careers"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	dbDsn := ""
	maxOpenConns := 25
	maxIdleConns := 25
	maxIdleTime := "15m"

	flag.StringVar(&dbDsn, "db-dsn", "falso", "PostgreSQL DSN")

	flag.Parse()

	db, err := openDb(dbDsn, maxOpenConns, maxIdleConns, maxIdleTime)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("./careers.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	query := `INSERT INTO careers (title) VALUES (:title) RETURNING id`

	preparedQuery, err := db.PrepareNamed(query)
	if err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		line := scanner.Bytes()

		var career careers.Career

		err := json.Unmarshal(line, &career)
		if err != nil {
			log.Fatal(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = preparedQuery.QueryRowxContext(ctx, career).Scan(&career.ID)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func openDb(dbDsn string, maxOpenConns int, maxIdleConns int, maxIdleTime string) (*sqlx.DB, error) {
	fmt.Println(dbDsn)
	db, err := sqlx.Open("postgres", dbDsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
