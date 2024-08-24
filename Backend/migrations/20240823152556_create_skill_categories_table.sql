-- +goose Up
-- +goose StatementBegin
CREATE TABLE skill_categories (
    id SERIAL PRIMARY KEY,
    career_id INT REFERENCES careers(id),
    category_name VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE skill_categories;
-- +goose StatementEnd
