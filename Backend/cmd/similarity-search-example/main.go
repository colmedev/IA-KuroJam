package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/colmedev/IA-KuroJam/Backend/careers"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pgvector/pgvector-go"
	"github.com/sashabaranov/go-openai"
)

func main() {
	dbDsn := ""
	maxOpenConns := 25
	maxIdleConns := 25
	maxIdleTime := "15m"
	llmApiKey := ""
	inputText := ""

	flag.StringVar(&dbDsn, "db-dsn", "", "PostgreSQL DSN")
	flag.StringVar(&llmApiKey, "llm-api-key", "", "LLM API Key")
	flag.StringVar(&inputText, "input-text", "", "Query Input Text")

	flag.Parse()

	db, err := openDb(dbDsn, maxOpenConns, maxIdleConns, maxIdleTime)
	if err != nil {
		log.Fatal(err)
	}

	client := openai.NewClient(llmApiKey)

	em := getDataEmbedding(client, inputText)

	v := pgvector.NewVector(em)

	query := `SELECT title FROM careers ORDER BY embedding <=> $1 LIMIT 10`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryxContext(ctx, query, v)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var career careers.Career

		if err := rows.Scan(&career.Title); err != nil {
			log.Fatal(err)
		}

		fmt.Println(career.Title)
	}

}

func getDataEmbedding(client *openai.Client, str string) []float32 {
	queryReq := openai.EmbeddingRequest{
		Input: []string{str},
		Model: openai.LargeEmbedding3,
	}

	queryResponse, err := client.CreateEmbeddings(context.Background(), queryReq)
	if err != nil {
		log.Fatal("Error creating query embedding:", err)
	}

	return queryResponse.Data[0].Embedding
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
