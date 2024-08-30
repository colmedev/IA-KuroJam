-- +goose Up
-- +goose StatementBegin
CREATE TABLE skill_areas(
    id SERIAL PRIMARY KEY,
    category_id INT REFERENCES skill_categories(id),
    category_name VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE skill_areas;
-- +goose StatementEnd
