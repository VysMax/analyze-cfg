package rest

import (
	"encoding/json"
	"net/http"

	"github.com/VysMax/analyze-cfg/analysis"
	"github.com/VysMax/analyze-cfg/models"
)

func AnalyzeREST(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req *models.Config

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	var problems analysis.Problems
	if err := problems.AnalyzeCfg(req); err != nil {
		http.Error(w, "Failed to analyse config", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(problems); err != nil {
		http.Error(w, "Failed to encode to JSON", http.StatusInternalServerError)
	}

}
