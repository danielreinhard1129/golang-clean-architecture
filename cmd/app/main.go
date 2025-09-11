package main

import (
	"github.com/danielreinhard1129/fiber-clean-arch/configs"
	"github.com/danielreinhard1129/fiber-clean-arch/internal/delivery/http/handler"
	"github.com/danielreinhard1129/fiber-clean-arch/internal/repository"
	"github.com/danielreinhard1129/fiber-clean-arch/internal/usecase"
	"github.com/danielreinhard1129/fiber-clean-arch/pkg/exception"
	"github.com/danielreinhard1129/fiber-clean-arch/pkg/mail"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// configurations
	config := configs.New()
	db := configs.NewDatabase(config)
	e := configs.NewEmail(config)

	// services
	mailService := mail.NewMailService(e.Host, e.Port, e.Username, e.Password, e.From)

	// user
	userRepository := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(&userRepository, mailService)
	userHandler := handler.NewUserHandler(&userUsecase)

	// auth
	authUsecase := usecase.NewAuthUsecase(&userRepository, config)
	authHandler := handler.NewAuthHandler(&authUsecase)

	// fiber
	app := fiber.New(configs.NewFiberConfiguration())
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover.New())

	// routes
	userHandler.Route(app)
	authHandler.Route(app)

	err := app.Listen(":" + config.Get("APP_PORT"))
	exception.PanicLogging(err)
}
