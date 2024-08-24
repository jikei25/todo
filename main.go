package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jikei25/todo/handler"
	"github.com/jikei25/todo/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)



func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	connection, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}
	queries := database.New(connection)

	apiCfg := handler.ApiConfig {
		DB: queries,
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
    	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    	AllowedHeaders:   []string{"*"},
    	ExposedHeaders:   []string{"Link"},
    	AllowCredentials: false,
    	MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handler.HandlerReadiness)
	v1Router.Post("/todo_item", apiCfg.HandlerCreateTodoItem)
	v1Router.Get("/todo_item", apiCfg.HandlerGetTodoItemByID)
	v1Router.Get("/todo_items", apiCfg.HandlerListTodoItem)
	v1Router.Patch("/todo_item", apiCfg.HandlerUpdateTodoItem)
	
	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}

	log.Println("Server starting on port", portString)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}