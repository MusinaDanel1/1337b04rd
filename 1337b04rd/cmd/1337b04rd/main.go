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
	"path/filepath"
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

	// Статические файлы (CSS, JS, изображения и т.п.)
	// Правильный путь для статических файлов
	// mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	// Главная страница, которая открывает catalog.html
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Указываем правильный путь к шаблону
		http.ServeFile(w, r, filepath.Join("web", "templates", "catalog.html"))
	})

	// Пример других страниц
	mux.HandleFunc("/archive-post", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "templates", "archive-post.html"))
	})
	mux.HandleFunc("/archive", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "templates", "archive.html"))
	})
	mux.HandleFunc("/create-post", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "templates", "create-post.html"))
	})
	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "templates", "error.html"))
	})
	mux.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "templates", "post.html"))
	})

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
