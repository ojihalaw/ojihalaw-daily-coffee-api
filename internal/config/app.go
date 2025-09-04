package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ojihalawa/daily-coffee-api.git/internal/delivery/http"
	"github.com/ojihalawa/daily-coffee-api.git/internal/delivery/http/route"
	"github.com/ojihalawa/daily-coffee-api.git/internal/repository"
	"github.com/ojihalawa/daily-coffee-api.git/internal/usecase"
	"github.com/ojihalawa/daily-coffee-api.git/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB        *gorm.DB
	App       *fiber.App
	Log       *logrus.Logger
	Validator *utils.Validator
	Config    *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	customerRepository := repository.NewCustomerRepository(config.Log)
	customerUseCase := usecase.NewCustomerUseCase(config.DB, config.Log, config.Validator, customerRepository)
	customerController := http.NewCustomerController(customerUseCase, config.Log)

	userRepository := repository.NewUserRepository(config.Log)
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validator, userRepository)
	userController := http.NewUserController(userUseCase, config.Log)

	categoryRepository := repository.NewCategoryRepository(config.Log)
	categoryUseCase := usecase.NewCategoryUseCase(config.DB, config.Log, config.Validator, categoryRepository)
	categoryController := http.NewCategoryController(categoryUseCase, config.Log)

	routeConfig := route.RouteConfig{
		App:                config.App,
		CustomerController: customerController,
		UserController:     userController,
		CategoryController: categoryController,
		// ContactController: contactController,
		// AddressController: addressController,
		// AuthMiddleware:    authMiddleware,
	}
	routeConfig.Setup()
}
