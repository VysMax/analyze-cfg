package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/VysMax/analyze-cfg/analysis"
	"github.com/VysMax/analyze-cfg/models"
)

func main() {
	isSilent := flag.Bool("s", false, "не выходить с ошибкой при наличии проблем")
	flag.BoolVar(isSilent, "silent", false, "не выходить с ошибкой при наличии проблем")

	flag.Parse()
	cfgPath := flag.Arg(0)

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file %s: %v\n", cfgPath, err)
		os.Exit(1)
	}

	var cfg models.Config

	format := filepath.Ext(cfgPath)
	switch {
	case format == ".json":
		err = json.Unmarshal(data, &cfg)
	case format == ".yaml" || format == ".yml":
		err = yaml.Unmarshal(data, &cfg)
	default:
		fmt.Fprintf(os.Stderr, "неподдерживаемый формат: %s\n", format)
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка извлечения данных из файла конфигурации")
		os.Exit(1)
	}

	problems := analysis.AnalyseCfg(&cfg)

	if len(problems) == 0 {
		fmt.Println("No problems")
		os.Exit(0)
	}

	switch {
	case len(problems)%10 == 1 && len(problems)%100 != 11:
		fmt.Printf("Обнаружена %d потенциально опасная настройка:\n\n", len(problems))
	case len(problems)%10 >= 2 && len(problems)%10 <= 4 && (len(problems)%100)/10 != 1:
		fmt.Printf("Обнаружено %d потенциально опасные настройки:\n\n", len(problems))
	default:
		fmt.Printf("Обнаружено %d потенциально опасных настроек:\n\n", len(problems))
	}

	for i, problem := range problems {
		fmt.Printf("%d. %s: %s\n\n", i+1, problem.Severity, problem.Description)
	}

	if *isSilent {
		os.Exit(0)
	}

	os.Exit(1)
}
