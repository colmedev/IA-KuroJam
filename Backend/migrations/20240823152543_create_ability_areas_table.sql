-- +goose Up
-- +goose StatementBegin
CREATE TABLE ability_areas (
    id SERIAL PRIMARY KEY,
    category_id INT REFERENCES ability_categories(id),
    area_name VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ability_areas;
-- +goose StatementEnd
