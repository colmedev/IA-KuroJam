package llm

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

type Service interface {
	GetQuestion(ctx context.Context, lastQuestion string, lastAnswer string, conversationSummary string, skills []string, previousQuestions []string) (string, error)
	UpdateSummary(ctx context.Context, lastQuestion string, lastAnswer string, conversationSummary string) (string, error)
	UpdateSkills(ctx context.Context, lastQuestion string, lastAnswer string, skills []string) ([]string, error)
	GetEmbeddings(ctx context.Context, str string) ([]float32, error)
}

type OpenAIService struct {
	client *openai.Client
}

func NewOpenAIService(llmApiKey string) *OpenAIService {
	return &OpenAIService{
		client: openai.NewClient(llmApiKey),
	}
}

var (
	ErrLlmCall = errors.New("llm api call failed")
)

func (ai *OpenAIService) GetQuestion(ctx context.Context, lastQuestion string, lastAnswer string, conversationSummary string, skills []string, previousQuestions []string) (string, error) {

	userMessage, err := questionPlaceholder.fill(lastQuestion, lastAnswer, strings.Join(skills, ", "), conversationSummary, strings.Join(previousQuestions, ", "))
	if err != nil {
		return "", fmt.Errorf("LLM Service: Get Question %w", err)
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	fmt.Println(userMessage)

	resp, err := ai.client.CreateChatCompletion(
		ctxTimeout,
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: questionPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userMessage,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("LLM Service: Get Question %w: %w", ErrLlmCall, err)
	}

	return strings.TrimPrefix(strings.TrimSpace(resp.Choices[0].Message.Content), "NEXT_QUESTION: "), nil
}

func (ai *OpenAIService) UpdateSummary(ctx context.Context, lastQuestion string, lastAnswer string, conversationSummary string) (string, error) {

	userMessage, err := summaryPlaceholder.fill(lastQuestion, lastAnswer, conversationSummary)
	if err != nil {
		return "", fmt.Errorf("LLM Service: Get Question %w", err)
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := ai.client.CreateChatCompletion(
		ctxTimeout,
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: summaryPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userMessage,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("LLM Service: Get Question %w: %w", ErrLlmCall, err)
	}

	return strings.Split(resp.Choices[0].Message.Content, ": ")[1], nil
}

func (ai *OpenAIService) UpdateSkills(ctx context.Context, lastQuestion string, lastAnswer string, skills []string) ([]string, error) {

	userMessage, err := skillsPlaceholder.fill(lastQuestion, lastAnswer, strings.Join(skills, ", "))
	if err != nil {
		return []string{}, fmt.Errorf("LLM Service: Get Question %w", err)
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := ai.client.CreateChatCompletion(
		ctxTimeout,
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: skillsPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userMessage,
				},
			},
		},
	)
	if err != nil {
		return []string{}, fmt.Errorf("LLM Service: Get Question %w: %w", ErrLlmCall, err)
	}

	skillsStr := strings.TrimPrefix(resp.Choices[0].Message.Content, "UPDATED_SKILLS: ")
	newSkills := strings.Split(skillsStr, ", ")

	return newSkills, nil
}

func (ai *OpenAIService) GetEmbeddings(ctx context.Context, str string) ([]float32, error) {
	queryReq := openai.EmbeddingRequest{
		Input: []string{str},
		Model: openai.LargeEmbedding3,
	}

	queryResponse, err := ai.client.CreateEmbeddings(ctx, queryReq)
	if err != nil {
		return []float32{}, fmt.Errorf("LLM Service: Get Question %w: %w", ErrLlmCall, err)
	}

	return queryResponse.Data[0].Embedding, nil
}
