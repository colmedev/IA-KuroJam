-- +goose Up
-- +goose StatementBegin
CREATE TABLE technology_areas (
    id SERIAL PRIMARY KEY,
    category_id INT REFERENCES technology_categories(id),
    area_name VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE technology_areas;
-- +goose StatementEnd
