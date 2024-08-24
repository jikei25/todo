package handler

import (
	"net/http"
	"github.com/jikei25/todo/response"
)

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	response.RespondWithJSON(w, 200, struct{
		Status string `json:"status"`
	}{
		Status: "success",
	})
}