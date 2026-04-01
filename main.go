package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"gopkg.in/yaml.v3"

	"github.com/VysMax/analyze-cfg/analysis"
	"github.com/VysMax/analyze-cfg/models"
)

func main() {
	isSilent := flag.Bool("s", false, "не выходить с ошибкой при наличии проблем")
	flag.BoolVar(isSilent, "silent", false, "не выходить с ошибкой при наличии проблем")
	isStdin := flag.Bool("stdin", false, "прочитать конфигурацию из стандартного потока ввода вместо файла")

	flag.Parse()

	var (
		r   io.Reader
		cfg models.Config
		err error
	)

	switch *isStdin {
	case true:
		r = os.Stdin
		err = analysis.SetReader(&r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка чтения из стандартного ввода: %v\n", err)
			os.Exit(1)
		}

		err = analysis.DecodeFromStdin(r, &cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка десериализации из стандартного ввода: %v\n", err)
			os.Exit(1)
		}

	case false:
		cfgPath := flag.Arg(0)

		file, err := os.Open(cfgPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка открытия файла %s: %v\n", cfgPath, err)
			os.Exit(1)
		}
		defer file.Close()

		r = file
		err = analysis.SetReader(&r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка чтения из стандартного ввода: %v\n", err)
			os.Exit(1)
		}

		format := path.Ext(cfgPath)

		switch format {
		case ".json":
			err = json.NewDecoder(r).Decode(&cfg)
		case ".yaml", ".yml":
			err = yaml.NewDecoder(r).Decode(&cfg)
		default:
			fmt.Fprintf(os.Stderr, "неподдерживаемый формат: %s\n", format)
			os.Exit(1)
		}
	}
	if err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "ошибка парсинга конфигурации: %v\n", err)
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
