package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"azlo-test-suite/dashboard" // <-- IMPORTANT: Replace with your module name

	"github.com/gorilla/mux"
)

// Handler holds dependencies for the handlers
type Handler struct {
	Dashboard *dashboard.TestDashboard
}

// ProjectPathRequest represents the request body for setting project path
type ProjectPathRequest struct {
	Path string `json:"path"`
}

// ProjectPathResponse represents the response for project path operations
type ProjectPathResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Path    string `json:"path,omitempty"`
}

func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := h.Dashboard.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	h.Dashboard.Clients[conn] = true
	defer delete(h.Dashboard.Clients, conn)

	// Send current data to new client
	conn.WriteJSON(h.Dashboard.Data)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (h *Handler) HandleRunTests(w http.ResponseWriter, r *http.Request) {
	go h.Dashboard.RunTests()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tests started"))
}

func (h *Handler) ServeCoverageData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	packageName := vars["package"]

	for _, result := range h.Dashboard.Data.Results {
		if result.Package == packageName {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result.Files)
			return
		}
	}
	http.Error(w, "Package not found", http.StatusNotFound)
}

// HandleSetProjectPath sets the project root directory
func (h *Handler) HandleSetProjectPath(w http.ResponseWriter, r *http.Request) {
	var req ProjectPathRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Path == "" {
		http.Error(w, "Path is required", http.StatusBadRequest)
		return
	}

	err := h.Dashboard.SetProjectPath(req.Path)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		response := ProjectPathResponse{
			Success: false,
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := ProjectPathResponse{
		Success: true,
		Message: "Project path updated successfully",
		Path:    req.Path,
	}
	json.NewEncoder(w).Encode(response)
}

// HandleGetProjectInfo returns current project information
func (h *Handler) HandleGetProjectInfo(w http.ResponseWriter, r *http.Request) {
	info := h.Dashboard.GetProjectInfo()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

// HandleHTMLCoverage serves HTML coverage reports with custom styling
func (h *Handler) HandleHTMLCoverage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	if filename == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}

	htmlContent, err := h.Dashboard.GetHTMLCoverage(filename)
	if err != nil {
		log.Printf("Error serving HTML coverage: %v", err)
		http.Error(w, "Coverage file not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(htmlContent))
}
