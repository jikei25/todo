-- +goose Up
CREATE TYPE status_enum AS ENUM ('pending', 'in_progress', 'completed', 'deleted');
CREATE TABLE todo_items (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status status_enum DEFAULT 'pending',
    due_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose Down
DROP TABLE todo_items;
DROP TYPE IF EXISTS status_enum;
