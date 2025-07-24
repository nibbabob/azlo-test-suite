package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"azlo-test-suite/dashboard" // <-- Replace
	"azlo-test-suite/handlers"  // <-- Replace

	"github.com/gorilla/mux"
)

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

	// Static file server for frontend assets
	staticFileServer := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFileServer))

	// Root handler to serve the index.html
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
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
