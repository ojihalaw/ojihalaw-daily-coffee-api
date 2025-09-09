package config

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/ojihalawa/daily-coffee-api.git/internal/delivery/http"
	"github.com/ojihalawa/daily-coffee-api.git/internal/delivery/http/middleware"
	"github.com/ojihalawa/daily-coffee-api.git/internal/delivery/http/route"
	"github.com/ojihalawa/daily-coffee-api.git/internal/repository"
	"github.com/ojihalawa/daily-coffee-api.git/internal/service"
	"github.com/ojihalawa/daily-coffee-api.git/internal/usecase"
	"github.com/ojihalawa/daily-coffee-api.git/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB         *gorm.DB
	App        *fiber.App
	Log        *logrus.Logger
	Validator  *utils.Validator
	Config     *viper.Viper
	JWTMaker   *utils.JWTMaker
	Cloudinary *cloudinary.Cloudinary
	Midtrans   *service.MidtransService
}

func Bootstrap(config *BootstrapConfig) {
	customerRepository := repository.NewCustomerRepository(config.Log)
	customerUseCase := usecase.NewCustomerUseCase(config.DB, config.Log, config.Validator, customerRepository)
	customerController := http.NewCustomerController(customerUseCase, config.Log)

	userRepository := repository.NewUserRepository(config.Log)
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validator, userRepository)
	userController := http.NewUserController(userUseCase, config.Log)

	sessionRepository := repository.NewRefreshRepository(config.Log)
	authUseCase := usecase.NewAuthUseCase(config.DB, config.Log, config.Validator, config.JWTMaker, userRepository, sessionRepository, customerRepository)
	authController := http.NewAuthController(authUseCase, config.Log)

	categoryRepository := repository.NewCategoryRepository(config.Log)
	categoryUseCase := usecase.NewCategoryUseCase(config.DB, config.Log, config.Validator, categoryRepository)
	categoryController := http.NewCategoryController(categoryUseCase, config.Log)

	productRepository := repository.NewProductRepository(config.Log)
	productUseCase := usecase.NewProductUseCase(config.DB, config.Log, config.Validator, config.Cloudinary, productRepository)
	productController := http.NewProductController(productUseCase, config.Log)

	orderRepository := repository.NewOrderRepository(config.Log)
	orderUseCase := usecase.NewOrderUseCase(config.DB, config.Log, config.Validator, orderRepository, config.Midtrans)
	orderController := http.NewOrderController(orderUseCase, config.Log)

	authMiddleware := middleware.AuthMiddleware(config.JWTMaker)

	routeConfig := route.RouteConfig{
		App:                config.App,
		AuthController:     authController,
		AuthMiddleware:     authMiddleware,
		CustomerController: customerController,
		UserController:     userController,
		CategoryController: categoryController,
		ProductController:  productController,
		OrderController:    orderController,
	}
	routeConfig.Setup()
}
