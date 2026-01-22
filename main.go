package main

import (
	"fmt"

	"www.github.com/a-fleming/blog-aggregator/internal/config"
)

func main() {
	var cfg config.Config
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Config data: %+v\n", cfg)
	}

	cfg.SetUser("aaron")

	cfg, err = config.Read()

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Config data: %+v\n", cfg)
	}

}
