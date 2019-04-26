.PHONY: build, run, test

GO_CMD=go

BIN_NAME=matchmaker
BIN_PATH=bin

DB_PATH=.data

setup:
	mkdir -p $(BIN_PATH)
	mkdir -p $(DB_PATH)

build: setup
	$(GO_CMD) build -o bin/$(BIN_NAME) 

run: build
	./$(BIN_PATH)/$(BIN_NAME) --port 8080

run-verbose: build
	./$(BIN_PATH)/$(BIN_NAME) --port 8080 --verbose

test:
	go test -v ./...
clean:
	rm -r $(BIN_PATH)/*
	rm -r $(DB_PATH)
	
