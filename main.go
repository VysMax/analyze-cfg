package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/VysMax/analyze-cfg/analysis"
	pb "github.com/VysMax/analyze-cfg/gen/proto"
	grpchandle "github.com/VysMax/analyze-cfg/grpc"
	"github.com/VysMax/analyze-cfg/models"
	rest "github.com/VysMax/analyze-cfg/rest"
	"github.com/VysMax/analyze-cfg/usecase"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	apiMode := flag.Bool("api", false, "Запуск в качестве REST API")
	grpcMode := flag.Bool("grpc", false, "Запуск в качестве gRPC")
	isSilent := flag.Bool("s", false, "не выходить с ошибкой при наличии проблем")
	flag.BoolVar(isSilent, "silent", false, "не выходить с ошибкой при наличии проблем")
	isStdin := flag.Bool("stdin", false, "прочитать конфигурацию из стандартного потока ввода вместо файла")

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	if *apiMode {
		analyzer := analysis.NewAnalyzer()
		usecase := usecase.New(analyzer)
		handler := rest.NewHandler(usecase)

		http.HandleFunc("/analyze", handler.Analyze)

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

	if *grpcMode {
		analyzer := analysis.NewAnalyzer()
		usecase := usecase.New(analyzer)
		handler := grpchandle.NewHandler(usecase)

		lis, err := net.Listen("tcp", os.Getenv("PORT"))
		if err != nil {
			log.Fatalf("ошибка запуска сервера: %v", err)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterCfgAnalyzerServer(grpcServer, handler)

		reflection.Register(grpcServer)

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		go func() {
			log.Printf("gRPC-сервер запущен на порту %s", os.Getenv("PORT"))
			if err = grpcServer.Serve(lis); err != nil {
				log.Fatalf("ошибка запуска сервера: %v", err)
			}
		}()

		<-quit
		log.Println("Получен сигнал завершения, выключение...")
		grpcServer.GracefulStop()
		log.Println("Сервер завершил работу.")
		return
	}

	var (
		cfg     models.Config
		message string
	)

	switch *isStdin {
	case true:

		problems, err := analysis.AnalyzeStdin(os.Stdin, &cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка десериализации из стандартного ввода: %v\n", err)
			os.Exit(1)
		}

		message = analysis.MessageBuilder("", problems)

		if len(problems) == 0 {
			fmt.Println(message)
			os.Exit(0)
		}

	case false:
		info, err := os.Stat(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка получения информации о файле %v\n", err)
			os.Exit(1)
		}

		if info.IsDir() {
			dirName := flag.Arg(0)
			multFileProblems, err := analysis.AnalyzeDirectory(dirName, &cfg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ошибка анализа папки %v\n", err)
				os.Exit(1)
			}

			messages := make([]string, 1)
			switch dirName {
			case ".":
				messages[0] = "Анализ текущей папки:"
			default:
				messages[0] = fmt.Sprintf("Анализ папки %s:\n", dirName)
			}

			if len(multFileProblems) == 0 {
				fmt.Println(messages[0])
				fmt.Println("Директория не содержит файлов конфигурации")
				os.Exit(0)
			}

			for _, problems := range multFileProblems {
				if len(problems) > 0 {
					message = analysis.MessageBuilder(problems[0].Filename, problems)
					messages = append(messages, message)
				}

			}

			message = strings.Join(messages, "\n")
		} else {
			cfg.File = flag.Arg(0)

			problems, err := analysis.AnalyzeFile(&cfg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ошибка анализа файла: %v\n", err)
				os.Exit(1)
			}

			message = analysis.MessageBuilder(cfg.File, problems)

			if len(problems) == 0 {
				fmt.Println(message)
				os.Exit(0)
			}
		}
	}
	fmt.Println(message)

	if *isSilent {
		os.Exit(0)
	}

	err = errors.New("потенциально опасные настройки обнаружены")
	fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
	os.Exit(1)
}
