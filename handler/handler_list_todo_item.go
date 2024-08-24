package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/jikei25/todo/internal/database"
	"github.com/jikei25/todo/response"
)

func (apiCfg *ApiConfig) HandlerListTodoItem(w http.ResponseWriter, r *http.Request) {
	limit := 5
	pageNumber := 1
	
	if r.URL.Query().Get("limit") != "" {
		var err error
		limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			response.RespondWithError(w, 400, "Invalid limit")
			return
		}
	}

	if r.URL.Query().Get("page_number") != "" {
		var err error
		pageNumber, err = strconv.Atoi(r.URL.Query().Get("page_number"))
		if err != nil {
			response.RespondWithError(w, 400, "Invalid page number")
			return
		}
	}

	total, err := apiCfg.DB.CountTodoItems(r.Context())
	if err != nil {
		log.Printf("Can't get number of todo items: %v", err)
		response.RespondWithError(w, 400, "Can't get number of todo items")
	}
	
	if total < int64(pageNumber) * int64(limit) {
		if total % int64(limit) == 0 {
			pageNumber = int(total) / limit
		} else {
			pageNumber = int(total) / limit + 1
		}

	}

	offSet := (pageNumber - 1) * limit

	dbTodoItems, err := apiCfg.DB.ListTodoItem(r.Context(), database.ListTodoItemParams{
		Limit: int32(limit),
		Offset: int32(offSet),
	})
	if err != nil {
		response.RespondWithError(w, 400, "Can't get todo items")
		return
	}
	
	todoItem := response.ConvertTodoItems(dbTodoItems)

	type Paging struct {
		Limit int `json:"limit"`
		PageNumber int `json:"page_number"`
		Total int64 `json:"total"`
	}
	
	type Data struct {
		Items []response.TodoItem `json:"todo_items"`
		PagingResponse Paging `json:"paging"`
	}

	response.RespondWithJSON(w, 200, Data{
		Items: todoItem,
		PagingResponse: Paging{
			Limit: limit,
			PageNumber: pageNumber,
			Total: total,
		},
	})
}