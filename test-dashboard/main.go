package main

import (
	"bytes" // Import the bytes package
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time" // Import the time package

	"azlo-test-suite/dashboard" // <-- Replace with your module path
	"azlo-test-suite/handlers"  // <-- Replace with your module path

	"github.com/gorilla/mux"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	// 1. Initialize the core application
	dash := dashboard.NewTestDashboard()
	go dash.BroadcastUpdates()

	// 2. Initialize the handlers with the dashboard instance
	h := &handlers.Handler{Dashboard: dash}

	// 3. Set up the router
	r := mux.NewRouter()

	// API and WebSocket routes
	r.HandleFunc("/ws", h.HandleWebSocket)
	r.HandleFunc("/run-tests", h.HandleRunTests).Methods("POST")
	r.HandleFunc("/coverage/{package}", h.ServeCoverageData)
	r.HandleFunc("/html-coverage/{filename}", h.HandleHTMLCoverage).Methods("GET")

	// New project path management routes
	r.HandleFunc("/set-project-path", h.HandleSetProjectPath).Methods("POST")
	r.HandleFunc("/project-info", h.HandleGetProjectInfo).Methods("GET")

	// --- ADDED: Handle favicon requests to prevent 404 errors ---
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Create a sub-filesystem for the static directory
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}

	// Static file server for frontend assets from embedded files
	staticFileServer := http.FileServer(http.FS(staticFS))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFileServer))

	// Root handler to serve the index.html from embedded files
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexHTML, err := fs.ReadFile(staticFS, "index.html")
		if err != nil {
			log.Printf("Error reading embedded index.html: %v", err)
			http.Error(w, "index.html not found", http.StatusInternalServerError)
			return
		}

		reader := bytes.NewReader(indexHTML)
		http.ServeContent(w, r, "index.html", time.Time{}, reader)
	})

	// 4. Start the server
	port := "8484"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	fmt.Printf("🧪 Go Test Dashboard starting on http://localhost:%s\n", port)
	fmt.Printf("📊 Open in your browser to see live test results and coverage\n")
	log.Fatal(http.ListenAndServe(":"+port, r))
}
