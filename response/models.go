package response

import "github.com/jikei25/todo/internal/database"

type TodoItem struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status"`
	Due_date    string `json:"due_date,omitempty"`
	Created_at  string `json:"created_at"`
	Updated_at  string `json:"updated_at"`
}

func ConvertTodoItem(todoItem database.TodoItem) TodoItem {
	description := ""
	if todoItem.Description.Valid == true {
		description = todoItem.Description.String
	}

	dueDate := ""
	if todoItem.DueDate.Valid == true {
		dueDate = todoItem.DueDate.Time.Format("02/01/2006")
	}

	return TodoItem{
		ID: int(todoItem.ID),
		Title: todoItem.Title,
		Description: description,
		Status: string(todoItem.Status.StatusEnum),
		Due_date: dueDate,
		Created_at: todoItem.CreatedAt.Time.Format("02/01/2006"),
		Updated_at: todoItem.UpdatedAt.Time.Format("02/01/2006"),
	}
}

func ConvertTodoItems(dbTodoItems []database.TodoItem) []TodoItem {
	todoItems := []TodoItem{}
	for _, dbTodoItem := range dbTodoItems {
		todoItems = append(todoItems, ConvertTodoItem(dbTodoItem))
	}
	return todoItems
}