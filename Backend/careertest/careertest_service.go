package careertest

import (
	"context"
	"errors"
	"fmt"

	"github.com/colmedev/IA-KuroJam/Backend/llm"
	"github.com/jmoiron/sqlx"
)

type Service interface {
	StartTest(ctx context.Context, ct *CareerTest) error
	GetQuestion(ctx context.Context, careerTestId int64, userId int64) (*Message, error)
	PostAnswer(ctx context.Context, lastAnswer string, careerTestId int64, userId int64) (*Message, error)
	GetResultsEmbedding(ctx context.Context, userId int64) ([]float32, error)
	GetActiveTest(ctx context.Context, userId int64) (*CareerTest, error)
}

type CareerTestService struct {
	store      store
	llmService llm.Service
}

var (
	ErrValidation    = errors.New("missing fields")
	ErrNotPermission = errors.New("no permission to edit")
)

func NewService(db *sqlx.DB, llmService llm.Service) *CareerTestService {
	cts := newCareerTestStore(db)

	return &CareerTestService{
		store:      &cts,
		llmService: llmService,
	}
}

func (cts *CareerTestService) StartTest(ctx context.Context, ct *CareerTest) error {

	if ct.UserId == 0 {
		return fmt.Errorf("career test service: %w: missing user id", ErrValidation)
	}

	return cts.store.Insert(ctx, ct)
}

func (cts *CareerTestService) PostAnswer(ctx context.Context, lastAnswer string, careerTestId int64, userId int64) (*Message, error) {

	careerTest, err := cts.store.Get(ctx, careerTestId)
	if err != nil {
		return nil, fmt.Errorf("CareerTestService PostAnswer: %w", err)
	}

	if careerTest.UserId != userId {
		return nil, fmt.Errorf("CareerTestService PostAnswer: %w", ErrNotPermission)
	}

	if careerTest.Status == "Completed" {
		return &Message{
			Sender:  SenderIA,
			Content: "Ha terminado la entrevista. Puedes proceder a ver los resultados",
		}, nil
	}

	// Update last answer
	careerTest.LastAnswer = lastAnswer

	// Generate summary
	summary, err := cts.llmService.UpdateSummary(ctx, careerTest.LastQuestion, careerTest.LastAnswer, careerTest.ConversationSummary)
	if err != nil {
		return nil, fmt.Errorf("CareerTestService PostAnswer: %w", err)
	}

	careerTest.ConversationSummary = summary

	// Generate new skills list
	skills, err := cts.llmService.UpdateSkills(ctx, careerTest.LastQuestion, lastAnswer, careerTest.Skills)
	if err != nil {
		return nil, fmt.Errorf("CareerTestService PostAnswer: %w", err)
	}

	careerTest.Skills = skills

	// Return new question
	newQuestion, err := cts.llmService.GetQuestion(ctx, careerTest.LastQuestion, careerTest.LastAnswer, careerTest.ConversationSummary, careerTest.Skills, careerTest.AIQuestions)
	if err != nil {
		return nil, fmt.Errorf("CareerTestService PostAnswer: %w", err)
	}

	careerTest.LastQuestion = newQuestion

	// Add messages to full conversation
	usrMsg := &Message{
		Sender:  SenderUser,
		Content: lastAnswer,
	}

	iaMsg := &Message{
		Sender:  SenderIA,
		Content: newQuestion,
	}

	careerTest.FullConversation = append(careerTest.FullConversation, *usrMsg)
	careerTest.FullConversation = append(careerTest.FullConversation, *iaMsg)
	careerTest.AIQuestions = append(careerTest.AIQuestions, newQuestion)

	if len(careerTest.AIQuestions) >= 10 {
		careerTest.Status = "Completed"
	} else {
		careerTest.Status = "In Process"
	}

	err = cts.store.Update(ctx, careerTest)
	if err != nil {
		return nil, fmt.Errorf("CareerTestService PostAnswer: %w", err)
	}

	return iaMsg, nil
}

func (cts *CareerTestService) GetQuestion(ctx context.Context, careerTestId int64, userId int64) (*Message, error) {
	careerTest, err := cts.store.Get(ctx, careerTestId)
	if err != nil {
		return nil, fmt.Errorf("CareerTestService GetQuestion: %w", err)
	}

	if careerTest.UserId != userId {
		return nil, fmt.Errorf("CareerTestService GetQuestion: %w", ErrNotPermission)
	}

	// Return new question
	newQuestion, err := cts.llmService.GetQuestion(ctx, careerTest.LastQuestion, careerTest.LastAnswer, careerTest.ConversationSummary, careerTest.Skills, careerTest.AIQuestions)
	if err != nil {
		return nil, fmt.Errorf("CareerTestService GetQuestion: %w", err)
	}

	careerTest.LastQuestion = newQuestion

	iaMsg := &Message{
		Sender:  SenderIA,
		Content: newQuestion,
	}

	careerTest.FullConversation = append(careerTest.FullConversation, *iaMsg)

	err = cts.store.Update(ctx, careerTest)
	if err != nil {
		return nil, fmt.Errorf("CareerTestService GetQuestion: %w", err)
	}

	return iaMsg, nil
}

func (cts *CareerTestService) GetResultsEmbedding(ctx context.Context, userId int64) ([]float32, error) {
	careerTest, err := cts.store.GetLastCompleted(ctx, userId)
	if err != nil {
		return []float32{}, fmt.Errorf("CareerTestService GetResultsString: %w", err)
	}

	if careerTest.UserId != userId {
		return []float32{}, fmt.Errorf("CareerTestService GetResultsString: %w", ErrNotPermission)
	}

	str := fmt.Sprintf(
		"%s\n%s\n",
		careerTest.ConversationSummary,
		careerTest.Skills,
	)

	emb, err := cts.llmService.GetEmbeddings(ctx, str)
	if err != nil {
		return []float32{}, fmt.Errorf("CareerTestService GetresultsString: %w", err)
	}

	return emb, nil
}

func (cts *CareerTestService) GetActiveTest(ctx context.Context, userId int64) (*CareerTest, error) {
	careerTest, err := cts.store.GetActive(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("CareerTestService GetActiveTest: %w", err)
	}

	return careerTest, err
}
