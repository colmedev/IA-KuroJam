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
	"github.com/pgvector/pgvector-go"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	dbDsn := ""
	maxOpenConns := 25
	maxIdleConns := 25
	maxIdleTime := "15m"
	llmApiKey := ""

	flag.StringVar(&dbDsn, "db-dsn", "falso", "PostgreSQL DSN")
	flag.StringVar(&llmApiKey, "llm-api-key", "", "LLM API Key")

	flag.Parse()

	db, err := openDb(dbDsn, maxOpenConns, maxIdleConns, maxIdleTime)
	if err != nil {
		log.Fatal(err)
	}

	client := openai.NewClient(llmApiKey)

	file, err := os.Open("./processed_careers.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	careerQuery := `INSERT INTO careers (title, description, personality_description,
		education, average_salary, lower_salary, highest_salary, embedding)
		VALUES
		(:title, :description, :personality_description, :education,
			:average_salary, :lower_salary, :highest_salary, :embedding)
		RETURNING id`

	preparedCareerQuery, err := db.PrepareNamed(careerQuery)
	if err != nil {
		log.Fatal(err)
	}

	tasksQuery := `INSERT INTO tasks (career_id, task_description)
		VALUES ($1, $2)`

	knowledgeCategoryQuery := `INSERT INTO knowledge_categories
		(career_id, category_name) VALUES ($1, $2)
		RETURNING id`

	knowledgeAreaQuery := `INSERT INTO knowledge_areas
		(category_id, area_name) VALUES($1, $2)`

	abilityCategoryQuery := `INSERT INTO ability_categories
		(career_id, category_name) VALUES ($1, $2)
		RETURNING id`

	abilityAreaQuery := `INSERT INTO ability_areas
		(category_id, area_name) VALUES ($1, $2)`

	skillCategoryQuery := `INSERT INTO skill_categories
		(career_id, category_name) VALUES ($1, $2)
		RETURNING id`

	skillAreaQuery := `INSERT INTO skill_areas
		(category_id, area_name) VALUES ($1, $2)`

	technologyCategoryQuery := `INSERT INTO technology_categories
		(career_id, category_name) VALUES ($1, $2)
		RETURNING id`

	technologyAreaQuery := `INSERT INTO technology_areas
		(category_id, area_name) VALUES ($1, $2)`

	personalityAttributesQuery := `INSERT INTO personality_attributes
		(career_id, attribute_name) VALUES ($1, $2)`

	for scanner.Scan() {
		line := scanner.Bytes()

		embedding := getDataEmbedding(client, string(line))

		var career careers.Career

		err := json.Unmarshal(line, &career)
		if err != nil {
			log.Fatal(err)
		}

		career.Embedding = pgvector.NewVector(embedding)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err = preparedCareerQuery.QueryRowxContext(ctx, career).Scan(&career.ID)
		if err != nil {
			log.Fatal(err)
		}

		// Tasks
		for _, t := range career.TasksString {
			_, err := db.ExecContext(ctx, tasksQuery, career.ID, t)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Knowledge
		for _, k := range career.Knowledge {
			err := db.QueryRowxContext(ctx, knowledgeCategoryQuery, career.ID, k.Name).Scan(&k.ID)
			if err != nil {
				log.Fatal(err)
			}

			for _, ka := range k.Areas {
				_, err := db.ExecContext(ctx, knowledgeAreaQuery, k.ID, ka)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		// Abilities
		for _, a := range career.Abilities {
			err := db.QueryRowxContext(ctx, abilityCategoryQuery, career.ID, a.Name).Scan(&a.ID)
			if err != nil {
				log.Fatal(err)
			}

			for _, aa := range a.Areas {
				fmt.Println(a, aa)
				_, err := db.ExecContext(ctx, abilityAreaQuery, a.ID, aa)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		// Skills
		for _, s := range career.SkillCategories {
			err := db.QueryRowxContext(ctx, skillCategoryQuery, career.ID, s.Name).Scan(&s.ID)
			if err != nil {
				log.Fatal(err)
			}

			for _, sa := range s.Areas {
				_, err := db.ExecContext(ctx, skillAreaQuery, s.ID, sa)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		// Technology
		for _, t := range career.TecnologyCategories {
			err := db.QueryRowxContext(ctx, technologyCategoryQuery, career.ID, t.Name).Scan(&t.ID)
			if err != nil {
				log.Fatal(err)
			}

			for _, ta := range t.Areas {
				_, err := db.ExecContext(ctx, technologyAreaQuery, t.ID, ta)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		// Personality
		for _, p := range career.Personality.Attributes {
			_, err := db.ExecContext(ctx, personalityAttributesQuery, career.ID, p)
			if err != nil {
				log.Fatal(err)
			}
		}
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
