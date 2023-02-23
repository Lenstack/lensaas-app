package main

import (
	"github.com/Lenstack/lensaas-app/internal/core/applications"
	"github.com/Lenstack/lensaas-app/internal/core/services"
	"github.com/Lenstack/lensaas-app/internal/infrastructure"
	"github.com/spf13/viper"
)

func main() {
	infrastructure.NewLoadEnv()

	var (
		AppEnvironment = viper.Get("APP_ENVIRONMENT").(string)
		AppPort        = viper.Get("APP_PORT").(string)
		DBHost         = viper.Get("DB_HOST").(string)
		DBPort         = viper.Get("DB_PORT").(string)
		DBUser         = viper.Get("DB_USER").(string)
		DBPassword     = viper.Get("DB_PASSWORD").(string)
		DBName         = viper.Get("DB_NAME").(string)
		MailHost       = viper.Get("MAIL_HOST").(string)
		MailPort       = viper.Get("MAIL_PORT").(string)
		MailEmail      = viper.Get("MAIL_EMAIL").(string)
		MailPass       = viper.Get("MAIL_PASSWORD").(string)
		JwtSecret      = viper.Get("JWT_SECRET").(string)
		JwtExpiration  = viper.Get("JWT_EXPIRATION").(string)
	)

	logger := infrastructure.NewLogger(AppEnvironment)
	postgres := infrastructure.NewPostgres(DBHost, DBPort, DBUser, DBPassword, DBName, logger.Log)

	// Register common services
	emailService := services.NewEmailService(MailHost, MailPort, MailEmail, MailPass)
	tokenService := services.NewTokenService(JwtSecret, JwtExpiration)
	// Register all services
	userService := services.NewUserService(postgres.Database, *tokenService, *emailService)
	// Register all applications
	microservice := applications.NewMicroservice(*emailService, *tokenService, *userService)

	routes := infrastructure.NewRoutes(*microservice)
	infrastructure.NewHttpServer(AppPort, routes.Handlers, logger.Log)
}
