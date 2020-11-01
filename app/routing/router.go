package routing

import (
	"simple-mpesa/app"
	"simple-mpesa/app/registry"
	"simple-mpesa/app/routing/account_handlers"
	"simple-mpesa/app/routing/error_handlers"
	"simple-mpesa/app/routing/middleware"
	"simple-mpesa/app/routing/user_handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Router(domain *registry.Domain, config app.Config) *fiber.App {

	srv := fiber.New(
		fiber.Config{ErrorHandler: error_handlers.ErrorHandler},
	)

	apiGroup := srv.Group("/api")
	apiGroup.Use(logger.New())

	apiRouteGroup(apiGroup, domain, config)

	return srv
}

func apiRouteGroup(api fiber.Router, domain *registry.Domain, config app.Config) {

	api.Post("/login/:user_type", user_handlers.Authenticate(domain, config))
	api.Post("/user/:user_type", user_handlers.Register(domain))

	// create group at /api/account
	account := api.Group("/account", middleware.AuthByBearerToken(config.Secret))
	account.Get("/balance", account_handlers.BalanceEnquiry(domain.Account))
	account.Get("/statement", account_handlers.MiniStatement(domain.Transaction))

	// create group at /api/transaction
	transaction := api.Group("/transaction", middleware.AuthByBearerToken(config.Secret))
	transaction.Post("/deposit", account_handlers.Deposit(domain.Account))
	// transaction.Post("/transfer", account_handlers.Withdraw(domain.Account))
	transaction.Post("/withdraw", account_handlers.Withdraw(domain.Account))
}
