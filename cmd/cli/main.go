package main

import (
	"log"

	"github.com/0gl04q/go-link-checker/internal/cli"
	"github.com/0gl04q/go-link-checker/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if err := cli.New(cfg).Run(); err != nil {
		log.Fatalf("cli: %v", err)
	}
}
