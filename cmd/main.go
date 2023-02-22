package main

import (
	"github.com/Lenstack/lensaas-app/internal/core/applications"
	"github.com/Lenstack/lensaas-app/internal/core/services"
	"github.com/Lenstack/lensaas-app/internal/infrastructure"
	"github.com/Lenstack/lensaas-app/internal/utils"
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
		MailUser       = viper.Get("MAIL_USER").(string)
		MailPass       = viper.Get("MAIL_PASSWORD").(string)
		JwtSecret      = viper.Get("JWT_SECRET").(string)
		JwtExpiration  = viper.Get("JWT_EXPIRATION").(string)
	)

	logger := infrastructure.NewLogger(AppEnvironment)
	postgres := infrastructure.NewPostgres(DBHost, DBPort, DBUser, DBPassword, DBName, logger)
	email := utils.NewEmail(MailHost, MailPort, MailUser, MailPass)
	jwt := utils.NewJwt(JwtSecret, JwtExpiration)

	userService := services.NewUserService(postgres.Database, jwt, email)
	microservice := applications.NewMicroservice(*userService)

	routes := infrastructure.NewRoutes(*microservice)
	infrastructure.NewHttpServer("localhost", AppPort, routes.Handlers, logger)
}
