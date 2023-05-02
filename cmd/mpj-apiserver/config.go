package main

import (
	"mpj/internal/models"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

func LoadConfig(path string) (*models.Config, error) {
	v := validator.New()
	cfg := &models.Config{}

	viper.SetConfigFile(path)
	viper.ReadInConfig()

	err := viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	err = v.Struct(cfg)
	if err != nil {
		return nil, err
	}

	// os.Setenv("INTRA_CLIENT", cfg.Intra.ClientID)
	// os.Setenv("INTRA_SECRET", cfg.Intra.SecretID)
	// os.Setenv("INTRA_RPS", strconv.FormatInt(int64(cfg.Intra.RPS), 10))

	return cfg, err
}
