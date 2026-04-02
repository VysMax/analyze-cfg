package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VysMax/analyze-cfg/analysis"
	"github.com/VysMax/analyze-cfg/handlers"
	"github.com/VysMax/analyze-cfg/models"
	"github.com/joho/godotenv"
)

func main() {
	apiMode := flag.Bool("api", false, "Запуск в качестве REST API")
	isSilent := flag.Bool("s", false, "не выходить с ошибкой при наличии проблем")
	flag.BoolVar(isSilent, "silent", false, "не выходить с ошибкой при наличии проблем")
	isStdin := flag.Bool("stdin", false, "прочитать конфигурацию из стандартного потока ввода вместо файла")

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	if *apiMode {
		http.HandleFunc("/analyse", handlers.AnalyseHandler)

		server := &http.Server{
			Addr: os.Getenv("PORT"),
		}

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			log.Printf("Запуск сервера на порту %s", server.Addr)
			err := server.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				log.Fatalf("ошибка запуска сервера: %v", err)
			}
		}()

		<-quit
		log.Println("Получен сигнал завершения, выключение...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Ошибка завершения работы сервера: %v", err)
		}

		log.Println("Сервер завершил работу.")
		return

	}

	var (
		cfg     models.Config
		message string
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
