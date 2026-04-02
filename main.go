package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/VysMax/analyze-cfg/analysis"
	"github.com/VysMax/analyze-cfg/handlers"
	"github.com/VysMax/analyze-cfg/models"
)

func main() {
	apiMode := flag.Bool("api", false, "Запуск в качестве REST API")

	isSilent := flag.Bool("s", false, "не выходить с ошибкой при наличии проблем")
	flag.BoolVar(isSilent, "silent", false, "не выходить с ошибкой при наличии проблем")
	isStdin := flag.Bool("stdin", false, "прочитать конфигурацию из стандартного потока ввода вместо файла")

	flag.Parse()

	if *apiMode {
		http.HandleFunc("/analyse", handlers.AnalyseHandler)

		port := "8080"

		log.Printf("Запуск сервера на порту %s", port)

		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("ошибка запуска сервера: %v", err)
		}
	}

	var (
		cfg     models.Config
		message string
		err     error
	)

	switch *isStdin {
	case true:

		message, err = analysis.AnalysisStdin(os.Stdin, &cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка десериализации из стандартного ввода: %v\n", err)
			os.Exit(1)
		}

	case false:
		info, err := os.Stat(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка получения информации о файле %v\n", err)
			os.Exit(1)
		}

		if info.IsDir() {
			message, err = analysis.AnalyseDirectory(flag.Arg(0), &cfg)
		} else {
			cfg.ConfigPath = flag.Arg(0)

			message, err = analysis.AnalyseFile(&cfg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ошибка анализа файла: %v\n", err)
				os.Exit(1)
			}
		}
	}

	fmt.Println(message)

	if *isSilent {
		os.Exit(0)
	}

	os.Exit(1)
}
