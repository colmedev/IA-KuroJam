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
	// Query to get the basic career information
	query := `SELECT
		id, title, description, personality_description, education, average_salary,
		lower_salary, highest_salary
		FROM careers
		ORDER BY embedding <=> $1 LIMIT 5`

	careers := make([]Career, 0)

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := cs.db.SelectContext(ctxTimeout, &careers, query, pgvector.NewVector(embedding))
	if err != nil {
		return []Career{}, fmt.Errorf("%w: %w", ErrGettingSuggestions, err)
	}

	// Fetch related data for each career
	for i := range careers {
		// Fetch tasks
		tasksQuery := `SELECT id, career_id, task_description FROM tasks WHERE career_id = $1`
		err = cs.db.SelectContext(ctxTimeout, &careers[i].Tasks, tasksQuery, careers[i].ID)
		if err != nil {
			return []Career{}, fmt.Errorf("%w: %w", ErrGettingSuggestions, err)
		}

		// Fetch knowledge categories and areas
		knowledgeCategoriesQuery := `SELECT id, career_id, category_name FROM knowledge_categories WHERE career_id = $1`
		err = cs.db.SelectContext(ctxTimeout, &careers[i].Knowledge, knowledgeCategoriesQuery, careers[i].ID)
		if err != nil {
			return []Career{}, fmt.Errorf("%w: %w", ErrGettingSuggestions, err)
		}
		for j := range careers[i].Knowledge {
			knowledgeAreasQuery := `SELECT area_name FROM knowledge_areas WHERE category_id = $1`
			err = cs.db.SelectContext(ctxTimeout, &careers[i].Knowledge[j].Areas, knowledgeAreasQuery, careers[i].Knowledge[j].ID)
			if err != nil {
				return []Career{}, fmt.Errorf("%w: %w", ErrGettingSuggestions, err)
			}
		}

		// Fetch ability categories and areas
		abilityCategoriesQuery := `SELECT id, career_id, category_name FROM ability_categories WHERE career_id = $1`
		err = cs.db.SelectContext(ctxTimeout, &careers[i].Abilities, abilityCategoriesQuery, careers[i].ID)
		if err != nil {
			return []Career{}, fmt.Errorf("%w: %w", ErrGettingSuggestions, err)
		}
		for j := range careers[i].Abilities {
			abilityAreasQuery := `SELECT area_name FROM ability_areas WHERE category_id = $1`
			err = cs.db.SelectContext(ctxTimeout, &careers[i].Abilities[j].Areas, abilityAreasQuery, careers[i].Abilities[j].ID)
			if err != nil {
				return []Career{}, fmt.Errorf("%w: %w", ErrGettingSuggestions, err)
			}

			fmt.Println(careers[i].Abilities)
		}

		// Fetch skill categories and areas
		skillCategoriesQuery := `SELECT id, career_id, category_name FROM skill_categories WHERE career_id = $1`
		err = cs.db.SelectContext(ctxTimeout, &careers[i].SkillCategories, skillCategoriesQuery, careers[i].ID)
		if err != nil {
			return []Career{}, fmt.Errorf("%w: %w", ErrGettingSuggestions, err)
		}
		for j := range careers[i].SkillCategories {
			skillAreasQuery := `SELECT area_name FROM skill_areas WHERE category_id = $1`
			err = cs.db.SelectContext(ctxTimeout, &careers[i].SkillCategories[j].Areas, skillAreasQuery, careers[i].SkillCategories[j].ID)
			if err != nil {
				return []Career{}, fmt.Errorf("%w: %w", ErrGettingSuggestions, err)
			}
		}

		// Fetch technology categories and areas
		technologyCategoriesQuery := `SELECT id, career_id, category_name FROM technology_categories WHERE career_id = $1`
		err = cs.db.SelectContext(ctxTimeout, &careers[i].TecnologyCategories, technologyCategoriesQuery, careers[i].ID)
		if err != nil {
			return []Career{}, fmt.Errorf("%w: %w", ErrGettingSuggestions, err)
		}
		for j := range careers[i].TecnologyCategories {
			technologyAreasQuery := `SELECT area_name FROM technology_areas WHERE category_id = $1`
			err = cs.db.SelectContext(ctxTimeout, &careers[i].TecnologyCategories[j].Areas, technologyAreasQuery, careers[i].TecnologyCategories[j].ID)
			if err != nil {
				return []Career{}, fmt.Errorf("%w: %w", ErrGettingSuggestions, err)
			}
		}

		// Fetch personality attributes
		personalityAttributesQuery := `SELECT attribute_name FROM personality_attributes WHERE career_id = $1`
		err = cs.db.SelectContext(ctxTimeout, &careers[i].Personality.Attributes, personalityAttributesQuery, careers[i].ID)
		if err != nil {
			return []Career{}, fmt.Errorf("%w: %w", ErrGettingSuggestions, err)
		}
	}

	return careers, nil
}
