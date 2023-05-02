package main

import (
	"context"
	"fmt"
	"log"

	"mpj/internal/models"
	"mpj/pkg/ent"

	_ "github.com/lib/pq"
)

func NewEntClient(ctx context.Context, cfg *models.ConfigPostgres) (*ent.Client, error) {

	connstring := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Username,
		cfg.Password,
		cfg.SSLMode,
	)

	postgres, err := ent.Open("postgres", connstring)
	if err != nil {
		log.Fatal(err)
	}

	return postgres, nil
}
