-- +goose Up
-- +goose StatementBegin
CREATE TABLE careers (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    personality_description TEXT,
    education VARCHAR(255),
    average_salary VARCHAR(255),
    lower_salary VARCHAR(255),
    highest_salary VARCHAR(255),
    embedding VECTOR(1536)
);
-- -- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE careers;
-- +goose StatementEnd
