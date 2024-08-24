-- name: CreateTodoItem :one
INSERT INTO todo_items (title, description, status, due_date)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTodoItemByID :one
SELECT * 
FROM todo_items
WHERE id = $1;

-- name: ListTodoItem :many
SELECT * 
FROM todo_items
LIMIT $1 OFFSET $2;

-- name: CountTodoItems :one
SELECT COUNT(*) 
FROM todo_items;

-- name: UpdateTodoItem :exec
UPDATE todo_items
SET status = $1, updated_at = NOW(), title = $2, description = $3, due_date = $4
WHERE id = $5;