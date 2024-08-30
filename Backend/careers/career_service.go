package careers

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
)

type Service interface {
	GetSimilarity(ctx context.Context, embedding []float32) ([]Career, error)
}

type CareerService struct {
	store store
}

func NewCareerService(db *sqlx.DB) *CareerService {
	s := newCareerStore(db)

	return &CareerService{
		store: &s,
	}
}

func (cs *CareerService) GetSimilarity(ctx context.Context, embedding []float32) ([]Career, error) {
	return cs.store.GetSuggestions(ctx, embedding)
}
