package main

import (
	"fmt"
	"log"

	"github.com/thedevscott/blogaggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %v\n", cfg)

	err = cfg.SetUser("Scott")

	if err != nil {
		log.Fatalf("Failed to set user: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	fmt.Printf("Read config: %v\n", cfg)
}
