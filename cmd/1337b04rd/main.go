package main

import (
	"1337b04rd/internal/Infrastructure/database"
	postrepo "1337b04rd/internal/Infrastructure/database/post"
	"1337b04rd/internal/domain/core"
	"1337b04rd/internal/domain/ports"
	httproutes "1337b04rd/internal/ui/http"
	"1337b04rd/internal/ui/http/post"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Подключение к БД
	dbConn, err := database.NewConnection()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Репозиторий
	postRepo := postrepo.NewPostgresPostRepository(dbConn)

	// Сервис
	var postService ports.PostService = core.NewPostService(postRepo)

	// HTTP-хендлер
	postHandler := post.NewPostHandler(postService)

	// Роутинг
	mux := http.NewServeMux()
	httproutes.RegisterPostRoutes(mux, postHandler)

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
