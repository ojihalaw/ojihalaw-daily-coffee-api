package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ojihalawa/daily-coffee-api.git/internal/delivery/http"
)

type RouteConfig struct {
	App                *fiber.App
	CustomerController *http.CustomerController
	UserController     *http.UserController
	CategoryController *http.CategoryController
	AuthController     *http.AuthController
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

	guest.Post("/users/register", c.CustomerController.Register)
	// c.App.Post("/api/users/_login", c.UserController.Login)
}

func (c *RouteConfig) SetupAuthRoute() {
	// c.App.Use(c.AuthMiddleware)
	// api := c.App.Group("/api/v1")
	// auth := api.Group("", c.AuthMiddleware)
	// c.App.Delete("/api/users", c.UserController.Logout)
	// c.App.Patch("/api/users/_current", c.UserController.Update)
	// c.App.Get("/api/users/_current", c.UserController.Current)

}

func (c *RouteConfig) SetupCMSRoute() {
	api := c.App.Group("/api/v1")
	cms := api.Group("/cms")

	auth := cms.Group("/auth")
	auth.Post("/login", c.AuthController.Login)

	c.App.Use(c.AuthMiddleware)
	user := cms.Group("/users")
	user.Post("", c.UserController.Create)

	category := cms.Group("/categories")
	category.Post("", c.CategoryController.Create)
	category.Get("", c.CategoryController.FindAll)
	category.Get(":id", c.CategoryController.FindByID)
	category.Put(":id", c.CategoryController.Update)
	category.Delete(":id", c.CategoryController.Delete)
}
