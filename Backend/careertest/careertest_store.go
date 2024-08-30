package careertest

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// TODO: Implemenbt service

type store interface {
	Insert(ctx context.Context, ct *CareerTest) error
	Update(ctx context.Context, ct *CareerTest) error
	Get(ctx context.Context, careerTestId int64) (*CareerTest, error)
	GetActive(ctx context.Context, userId int64) (*CareerTest, error)
}

type careerTestStore struct {
	db *sqlx.DB
}

func newCareerTestStore(db *sqlx.DB) careerTestStore {
	return careerTestStore{
		db: db,
	}
}

// Errors
var (
	ErrEncoding       = errors.New("data encoding error")
	ErrEditConflict   = errors.New("conflict error")
	ErrRecordNotFound = errors.New("career test not found")
)

func (cts *careerTestStore) Insert(ctx context.Context, ct *CareerTest) error {
	query := `INSERT INTO career_tests (user_id)
		VALUES ($1)
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := cts.db.QueryRowxContext(ctx, query, ct.UserId).Scan(&ct.ID)
	if err != nil {
		return fmt.Errorf("error inserting career test: %w", err)
	}

	return nil
}

func (cts *careerTestStore) Update(ctx context.Context, ct *CareerTest) error {
	query := `UPDATE career_tests SET
		full_conversation = $1,
		conversation_summary = $2, 
		status = $3,
		skills = $4,
		last_question = $5,	
		last_answer = $6,
		ai_questions = $7
		WHERE id = $8
		RETURNING version
	`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	fcData, err := json.Marshal(ct.FullConversation)
	if err != nil {
		return fmt.Errorf("error inserting career test: %w", err)
	}

	err = cts.db.QueryRowxContext(
		ctx,
		query,
		fcData,
		ct.ConversationSummary,
		ct.Status,
		pq.Array(ct.Skills),
		ct.LastQuestion,
		ct.LastAnswer,
		pq.Array(ct.AIQuestions),
		ct.ID,
	).Scan(&ct.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return fmt.Errorf("error updating career test: %w", ErrEditConflict)
		default:
			return fmt.Errorf("error updating career test: %w", err)
		}
	}

	return nil
}

func (cts *careerTestStore) Get(ctx context.Context, careerTestId int64) (*CareerTest, error) {
	query := `SELECT id, user_id, full_conversation, conversation_summary, 
		status, skills,  version, last_question, last_answer, ai_questions
		FROM career_tests WHERE id = $1
		`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var careerTest CareerTest
	var fullConversation []byte
	var conversationSummary sql.NullString
	var lastQuestion sql.NullString
	var lastAnswer sql.NullString

	err := cts.db.QueryRowxContext(ctx, query, careerTestId).Scan(
		&careerTest.ID,
		&careerTest.UserId,
		&fullConversation,
		&conversationSummary,
		&careerTest.Status,
		pq.Array(&careerTest.Skills),
		&careerTest.Version,
		&lastQuestion,
		&lastAnswer,
		pq.Array(&careerTest.AIQuestions),
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("%w (id %v)", ErrRecordNotFound, careerTestId)
		default:
			return nil, fmt.Errorf("error fetching career test (id: %v): %w", careerTestId, err)
		}
	}

	if len(fullConversation) > 0 {
		if err := json.Unmarshal(fullConversation, &careerTest.FullConversation); err != nil {
			return nil, fmt.Errorf("error unmarshaling full_conversation for career test (id: %v): %w", careerTestId, err)
		}
	}

	if conversationSummary.Valid {
		careerTest.ConversationSummary = conversationSummary.String
	}

	if lastQuestion.Valid {
		careerTest.LastQuestion = lastQuestion.String
	}

	if lastAnswer.Valid {
		careerTest.LastAnswer = lastAnswer.String
	}

	return &careerTest, nil
}

func (cts *careerTestStore) GetActive(ctx context.Context, userId int64) (*CareerTest, error) {
	query := `SELECT id, user_id, full_conversation, conversation_summary, 
		status, skills,  version, last_question, last_answer, ai_questions
		FROM career_tests
		WHERE status != 'Completed'
		AND user_id = $1
		ORDER BY id DESC
		LIMIT 1
		`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var careerTest CareerTest
	var fullConversation []byte
	var conversationSummary sql.NullString
	var lastQuestion sql.NullString
	var lastAnswer sql.NullString

	err := cts.db.QueryRowxContext(ctx, query, userId).Scan(
		&careerTest.ID,
		&careerTest.UserId,
		&fullConversation,
		&conversationSummary,
		&careerTest.Status,
		pq.Array(&careerTest.Skills),
		&careerTest.Version,
		&lastQuestion,
		&lastAnswer,
		pq.Array(&careerTest.AIQuestions),
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("%w (not active career test found)", ErrRecordNotFound)
		default:
			return nil, fmt.Errorf("error fetching career test: %w", err)
		}
	}

	if len(fullConversation) > 0 {
		if err := json.Unmarshal(fullConversation, &careerTest.FullConversation); err != nil {
			return nil, fmt.Errorf("error unmarshaling full_conversation for career test (id: %v): %w", careerTest.ID, err)
		}
	}

	if conversationSummary.Valid {
		careerTest.ConversationSummary = conversationSummary.String
	}

	if lastQuestion.Valid {
		careerTest.LastQuestion = lastQuestion.String
	}

	if lastAnswer.Valid {
		careerTest.LastAnswer = lastAnswer.String
	}

	return &careerTest, nil
}
