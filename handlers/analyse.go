package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VysMax/analyze-cfg/analysis"
	"github.com/VysMax/analyze-cfg/models"
)

func AnalyseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req *models.Config

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	var problems analysis.Problems
	problems.AnalyseCfg(req)

	resp := analysis.MessageBuilder("", problems)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp))
}
