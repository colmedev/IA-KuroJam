package careers

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pgvector/pgvector-go"
	"golang.org/x/net/context"
)

type store interface {
	GetSuggestions(ctx context.Context, embedding []float32) ([]Career, error)
}

type careerStore struct {
	db *sqlx.DB
}

func newCareerStore(db *sqlx.DB) careerStore {
	return careerStore{
		db: db,
	}
}

var (
	ErrGettingSuggestions = errors.New("error executing suggestions query")
)

func (cs *careerStore) GetSuggestions(ctx context.Context, embedding []float32) ([]Career, error) {
	query := `SELECT
		title, description, personality_description, education, average_salary,
		lower_salary, highest_salary
		FROM careers
		ORDER BY embedding <=> $1 LIMIT 3
	`

	careers := make([]Career, 3)

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := cs.db.SelectContext(ctxTimeout, &careers, query, pgvector.NewVector(embedding))
	if err != nil {
		return []Career{}, fmt.Errorf("%w: %w", ErrGettingSuggestions, err)
	}

	return careers, nil
}
