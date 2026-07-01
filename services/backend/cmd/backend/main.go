package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/session-manager/backend/internal/api"
)

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	projectsRoot := filepath.Join(os.Getenv("HOME"), ".claude", "projects")
	http.Handle(
		"/api/claude-code/sessions",
		withCORS(api.NewClaudeCodeSessionsHandler(api.NewClaudeCodeSessionTreeProvider(projectsRoot))),
	)

	log.Println("Backend starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
