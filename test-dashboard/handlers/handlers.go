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

func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := h.Dashboard.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	h.Dashboard.Clients[conn] = true
	defer delete(h.Dashboard.Clients, conn)

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
