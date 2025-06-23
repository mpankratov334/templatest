package api

import (
	"Templatest/internal/api/middleware"
	"Templatest/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Routers struct {
	Service service.Service
}

func NewRouters(r *Routers, token string) *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowMethods:  "GET, POST, PUT, DELETE",
		AllowHeaders:  "Accept, Authorization, Content-Type, X-CSRF-Token, X-REQUEST-ID",
		ExposeHeaders: "Link",
		MaxAge:        300,
	}))

	apiGroup := app.Group("/v1", middleware.Authorization(token))

	apiGroup.Post("/tasks", r.Service.Create)
	apiGroup.Get("/tasks", r.Service.ReadAll)
	apiGroup.Get("/tasks/:id", r.Service.Read)
	apiGroup.Put("/tasks/:id", r.Service.Update)
	apiGroup.Delete("/tasks/:id", r.Service.Delete)

	return app
}
