package analysis

import (
	"fmt"

	"github.com/VysMax/analyze-cfg/models"
)

func AnalyseCfg(cfg *models.Config) []models.Problem {
	var problems []models.Problem

	if cfg.Log.Level == "debug" {
		problems = append(problems, models.Problem{
			Path:        "add later",
			Description: "логирование в debug-режиме. Поменяйте режим на более избирательный (info+).",
			Severity:    "LOW",
		})
	}

	if cfg.DigestAlgorithm == "md5" {
		problems = append(problems, models.Problem{
			Path:        "add later",
			Description: fmt.Sprintf("слишком слабый алгоритм - %s. Замените его на более безопасный.", cfg.DigestAlgorithm),
			Severity:    "HIGH",
		})
	}

	return problems
}
