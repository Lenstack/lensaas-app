package main

import (
	"github.com/Lenstack/lensaas-app/internal/infrastructure"
	"github.com/Lenstack/lensaas-app/internal/utils"
	"github.com/spf13/viper"
)

func main() {
	infrastructure.NewLoadEnv()

	var (
		AppEnvironment = viper.Get("APP_ENVIRONMENT").(string)
		AppPort        = viper.Get("APP_PORT").(string)
		JwtSecret      = viper.Get("JWT_SECRET").(string)
		JwtExpiration  = viper.Get("JWT_EXPIRATION").(string)
	)

	logger := infrastructure.NewLogger(AppEnvironment)
	utils.NewJwt(JwtSecret, JwtExpiration)

	routes := infrastructure.NewRoutes()
	infrastructure.NewHttpServer("localhost", AppPort, routes.Handlers, logger)
}
