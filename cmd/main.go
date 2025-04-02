package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := flag.Int("port", 8080, "Port number for server")
	help := flag.Bool("help", false, "Show help message")

	flag.Parse()

	if *help {
		fmt.Println("hacker board")
		fmt.Println("Usage:")
		fmt.Println("  1337b04rd [--port <N>]")
		fmt.Println("  1337b04rd --help")
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	http.HandleFunc("/set-cookie", setCookieHandler)
	http.HandleFunc("/get-cookie", getCookieHandler)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Server is starting on host %s...", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Starting server error: %v", err)
	}
}

// установка кук
func setCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    "123456",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
	w.Write([]byte("Cookie set!"))
}

// чтение кук
func getCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusNotFound)
		return
	}
	w.Write([]byte("Cookie value: " + cookie.Value))
}
