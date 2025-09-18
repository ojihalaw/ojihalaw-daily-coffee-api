package main

import (
	"fmt"

	"github.com/ojihalawa/daily-coffee-api.git/internal/config"
	"github.com/ojihalawa/daily-coffee-api.git/internal/migration"
	"github.com/ojihalawa/daily-coffee-api.git/internal/service"
	"github.com/ojihalawa/daily-coffee-api.git/internal/utils"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validator := utils.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)
	jwtMaker := utils.NewJWTMaker(viperConfig)
	cloudinary := config.NewCloudinary(viperConfig)
	midClient := config.NewMidtransClient(viperConfig)
	midtransService := service.NewMidtransService(midClient)
	redisClient := config.NewRedisClient(viperConfig, log)

	migration.Run(db, log)

	config.Bootstrap(&config.BootstrapConfig{
		DB:          db,
		App:         app,
		Log:         log,
		Validator:   validator,
		Config:      viperConfig,
		JWTMaker:    jwtMaker,
		Cloudinary:  cloudinary,
		Midtrans:    midtransService,
		RedisClient: redisClient,
	})

	webPort := viperConfig.GetInt("APP_PORT")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
