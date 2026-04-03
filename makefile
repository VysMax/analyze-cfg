.PHONY: build

APP_NAME = Analyze-Cfg


build:
	@echo "Building $(APP_NAME)"
	go build -o $(APP_NAME)

analyse-current-directory:
	./$(APP_NAME) .

grpc-gen:
	protoc --go_out=./gen --go-grpc_out=./gen proto/analyze-cfg.proto
