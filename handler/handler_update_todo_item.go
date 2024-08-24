package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jikei25/todo/internal/database"
	"github.com/jikei25/todo/response"
)

func (apiCfg *ApiConfig) HandlerUpdateTodoItem(w http.ResponseWriter, r *http.Request) {
	var todoItemUpdate struct {
        Title       *string `json:"title,omitempty"`
        Description *string `json:"description,omitempty"`
        Status      *string `json:"status,omitempty"`
        DueDate     *string `json:"due_date,omitempty"`
    }

	err := json.NewDecoder(r.Body).Decode(&todoItemUpdate)
	if err != nil {
		response.RespondWithError(w, 400, fmt.Sprintf("Can't parse JSON: %v", err))
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.RespondWithError(w, 400, "Invalid ID")
		return
	}
	var (
        status      database.NullStatusEnum
        title       *string
        description sql.NullString
        dueDate     sql.NullTime
    )

	curTodoItem, err := apiCfg.DB.GetTodoItemByID(r.Context(), int32(id))
	if err != nil {
		response.RespondWithError(w, 400, "Can't find item with the given ID")
		return
	}

	status = database.NullStatusEnum{
		StatusEnum: curTodoItem.Status.StatusEnum,
		Valid: true,
	}
	if todoItemUpdate.Status != nil {
		switch *todoItemUpdate.Status {
			case string(database.StatusEnumCompleted),
			string(database.StatusEnumDeleted), 
			string(database.StatusEnumInProgress), 
			string(database.StatusEnumPending):
				status.StatusEnum = database.StatusEnum(*todoItemUpdate.Status)
		default:
			response.RespondWithError(w, 400, "Invalid status type")
		}
	}

	if curTodoItem.Description.Valid == true {
		description.String = curTodoItem.Description.String
		description.Valid = true
	}
	if todoItemUpdate.Description != nil {
		description.String = *todoItemUpdate.Description
	}

	if todoItemUpdate.Title != nil {
		title = todoItemUpdate.Title
	}

	if curTodoItem.DueDate.Valid == true {
		dueDate.Time = curTodoItem.DueDate.Time
		dueDate.Valid = true
	}
	if todoItemUpdate.DueDate != nil {
		parseDate, err := time.Parse("02/01/2006", *todoItemUpdate.DueDate)
		if err != nil {
			response.RespondWithError(w, 400, "Invalid date format (Correct format: dd/mm/yyyy)")
			return
		}
		dueDate.Time = parseDate
	}

	title = &curTodoItem.Title
	if todoItemUpdate.Title != nil {
		title = todoItemUpdate.Title
	}

	err = apiCfg.DB.UpdateTodoItem(r.Context(), database.UpdateTodoItemParams{
		Status: status,
		Title: *title,
		Description: description,
		DueDate: dueDate,
		ID: int32(id),
	})
	
	if err != nil {
		log.Println(err)
		response.RespondWithError(w, 500, "Can't update todo item")
		return
	}

	response.RespondWithJSON(w, 200, struct{
		Status string `json:"status"`
	}{
		Status: "succesful",
	})
}