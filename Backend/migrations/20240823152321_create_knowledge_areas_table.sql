-- +goose Up
-- +goose StatementBegin
CREATE TABLE knowledge_areas (
    id SERIAL PRIMARY KEY,
    category_id INT REFERENCES knowledge_categories(id),
    area_name VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE knowledge_areas;
-- +goose StatementEnd
