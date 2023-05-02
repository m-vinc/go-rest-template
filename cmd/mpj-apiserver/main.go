package main

import (
	"context"
	"log"

	application "mpj/internal/application"
	"mpj/internal/services"

	"github.com/spf13/cobra"
)

var configPath *string

var VERSION = "dev"

func init() {
	configPath = rootCmd.PersistentFlags().StringP("config", "c", "~/.config/idcheck.yml", "path to yaml file of idcheck configurations")
}

var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "Run the mpj-apiserver",
	Long:    ``,
	Version: VERSION,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(log.Ldate | log.Ltime)
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		cfg, err := LoadConfig(*configPath)
		if err != nil {
			log.Fatal(err)
		}

		database, err := NewEntClient(ctx, cfg.Postgres)
		if err != nil {
			log.Fatal(err)
		}

		if err := database.Schema.Create(cmd.Context()); err != nil {
			log.Fatal(err)
		}

		gatewayService := services.NewGatewayService(services.NewLoggerService("Service", "Live"))
		usersService := services.NewUsersService(database, services.NewLoggerService("Service", "Users"))

		usersController := application.NewUsersController(
			usersService,
			gatewayService,
			cfg.Application,
			services.NewLoggerService("Controller", "Users"),
		)
		liveController := application.NewLiveController(gatewayService, cfg.Application, services.NewLoggerService("Controller", "Live"))

		app := application.New(
			cfg.Application,
			services.NewLoggerService("Controller", "Application"),
			liveController,
			usersController,
		)

		app.Run(cmd.Context())
	},
}

var rootCmd = &cobra.Command{
	Use:   "mpj-apiserver",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func main() {
	rootCmd.AddCommand(runCmd)

	rootCmd.Execute()
}
