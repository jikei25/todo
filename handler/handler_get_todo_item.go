package handler

import (
	"net/http"
	"strconv"

	"github.com/jikei25/todo/response"
)

func (apiCfg *ApiConfig) HandlerGetTodoItemByID(w http.ResponseWriter, r *http.Request) {
	todoItemIDStr := r.URL.Query().Get("id")
	todoItemID, err := strconv.Atoi(todoItemIDStr)
	if err != nil {
		response.RespondWithError(w, 400, "Invalid ID type")
		return
	}
	todoItem, err := apiCfg.DB.GetTodoItemByID(r.Context(), int32(todoItemID))
	if err != nil {
		response.RespondWithError(w, 400, "Can't find item with the given ID")
		return
	}
	response.RespondWithJSON(w, 200, response.ConvertTodoItem(todoItem))
}