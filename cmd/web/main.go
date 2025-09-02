package main

import (
	"fmt"

	"github.com/ojihalawa/daily-coffee-api.git/internal/config"
	"github.com/ojihalawa/daily-coffee-api.git/internal/migration"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)

	migration.Run(db, log)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
	})

	webPort := viperConfig.GetInt("APP_PORT")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
