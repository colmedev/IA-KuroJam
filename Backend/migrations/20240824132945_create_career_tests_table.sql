-- +goose Up
-- +goose StatementBegin
CREATE TABLE career_tests (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    full_conversation JSONB,
    conversation_summary TEXT,
    skills text[],
    status VARCHAR(50) NOT NULL DEFAULT 'New',
    last_question TEXT,
    last_answer TEXT,
    ai_questions text[], 
    version INT NOT NULL DEFAULT 1
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE career_tests;
-- +goose StatementEnd
