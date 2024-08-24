package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jikei25/todo/internal/database"
	"github.com/jikei25/todo/response"
)

type ApiConfig struct {
	DB *database.Queries
}

func (apiCfg * ApiConfig) HandlerCreateTodoItem(w http.ResponseWriter, r *http.Request) {
	var todo_item struct {
		Title 		string `json:"title"`
		Description string `json:"description,omitempty"`
		Status 		string `json:"status"`
		Due_date 	string `json:"due_date,omitempty"`

	}
	err := json.NewDecoder(r.Body).Decode(&todo_item)
	if err != nil {
		response.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	description := sql.NullString{}
	if todo_item.Description != "" {
		description.String = todo_item.Description
		description.Valid = true
	}

	dueDate := sql.NullTime{}
	if todo_item.Due_date != "" {
		parseDueDate, err := time.Parse("02/01/2006", todo_item.Due_date)
		if err != nil {
			response.RespondWithError(w, 400, "Invalid date format (Correct format: dd/mm/yyyy)")
			return
		}
		dueDate.Time = parseDueDate
		dueDate.Valid = true
	}
	
	status := database.NullStatusEnum{}
	if todo_item.Status != "" {
		switch todo_item.Status {
			case string(database.StatusEnumCompleted), 
			string(database.StatusEnumDeleted), 
			string(database.StatusEnumInProgress), 
			string(database.StatusEnumPending):
				status.StatusEnum = database.StatusEnum(todo_item.Status)
				status.Valid = true
		default:
			response.RespondWithError(w, 400, "Invalid status type")
		}
	}

	return_todo_item, err := apiCfg.DB.CreateTodoItem(r.Context(), database.CreateTodoItemParams{
		Title: todo_item.Title,
		Description: description,
		Status: status,
		DueDate: dueDate,
	})

	response.RespondWithJSON(w, 201, response.ConvertTodoItem(return_todo_item))
}