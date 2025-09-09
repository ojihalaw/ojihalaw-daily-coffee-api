package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ojihalawa/daily-coffee-api.git/internal/delivery/http"
)

type RouteConfig struct {
	App                *fiber.App
	AuthController     *http.AuthController
	CustomerController *http.CustomerController
	UserController     *http.UserController
	CategoryController *http.CategoryController
	ProductController  *http.ProductController
	OrderController    *http.OrderController
	AuthMiddleware     fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupAuthRoute()
	c.SetupGuestRoute()
	c.SetupCMSRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	api := c.App.Group("/api/v1")
	guest := api.Group("/guest")

	order := guest.Group("/orders")
	order.Post("", c.OrderController.Create)
}

func (c *RouteConfig) SetupAuthRoute() {
	api := c.App.Group("/api/v1")
	guest := api.Group("/guest")

	auth := guest.Group("/auth")
	auth.Post("/register", c.CustomerController.Register)
	auth.Post("/login", c.AuthController.LoginClientWithEmail)
}

func (c *RouteConfig) SetupCMSRoute() {
	api := c.App.Group("/api/v1")
	cms := api.Group("/cms")

	auth := cms.Group("/auth")
	auth.Post("/login", c.AuthController.Login)

	c.App.Use(c.AuthMiddleware)
	user := cms.Group("/users")
	user.Post("", c.UserController.Create)
	user.Get("", c.UserController.FindAll)
	user.Get(":id", c.UserController.FindByID)
	user.Put(":id", c.UserController.Update)
	user.Delete(":id", c.UserController.Delete)

	category := cms.Group("/categories")
	category.Post("", c.CategoryController.Create)
	category.Get("", c.CategoryController.FindAll)
	category.Get(":id", c.CategoryController.FindByID)
	category.Put(":id", c.CategoryController.Update)
	category.Delete(":id", c.CategoryController.Delete)

	product := cms.Group("/products")
	product.Post("", c.ProductController.Create)
	product.Get("", c.ProductController.FindAll)
	product.Get(":id", c.ProductController.FindByID)
	product.Put(":id", c.ProductController.Update)
	product.Delete(":id", c.ProductController.Delete)
	product.Get("/product/special", c.ProductController.FindSpecialProduct)
	product.Put("/product/special", c.ProductController.UpdateSpecialProduct)

	customer := cms.Group("/customers")
	customer.Post("", c.CustomerController.Register)
	customer.Get("", c.CustomerController.FindAll)
	customer.Get(":id", c.CustomerController.FindByID)
	customer.Put(":id", c.CustomerController.Update)
	customer.Delete(":id", c.CustomerController.Delete)
}
