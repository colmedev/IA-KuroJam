package careertest

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type CareerTest struct {
	ID                  string       `db:"id" json:"id"`
	UserId              int64        `db:"user_id" json:"userId"`
	FullConversation    Conversation `db:"full_conversation" json:"fullConversation"`
	ConversationSummary string       `db:"conversation_summary" json:"conversationSummary"`
	Status              string       `db:"status" json:"status"`
	Skills              []string     `db:"skills" json:"skills"`
	LastQuestion        string       `db:"last_question" json:"lastQuestion"`
	LastAnswer          string       `db:"last_answer" json:"lastAnswer"`
	AIQuestions         []string     `db:"ai_questions" json:"-"`
	Version             int          `db:"version" json:"version"`
}

type Conversation []Message

type Message struct {
	Sender  Sender `db:"sender" json:"sender"`
	Content string `db:"content" json:"content"`
}

type Sender string

const (
	SenderIA   Sender = "IA"
	SenderUser Sender = "User"
)

func (c Conversation) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Conversation) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &c)
}

func (m Message) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *Message) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &m)
}
