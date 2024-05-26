package main

import (
	"log"

	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/app"
	"github.com/DimaGitHahahab/ozon-fintech-posts/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	a := app.New(cfg)

	a.Run()
}
