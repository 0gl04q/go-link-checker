.PHONY: run
run:
	go run ./cmd/cli/ $(ARGS)

.PHONY: build
build:
	go build -o bin/worker ./cmd/cli/