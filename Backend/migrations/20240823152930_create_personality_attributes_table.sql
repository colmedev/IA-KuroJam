-- +goose Up
-- +goose StatementBegin
CREATE TABLE personality_attributes (
    id SERIAL PRIMARY KEY,
    career_id INT REFERENCES careers(id),
    attribute_name VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE personality_attributes;
-- +goose StatementEnd
