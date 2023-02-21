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
		MailHost       = viper.Get("MAIL_HOST").(string)
		MailPort       = viper.Get("MAIL_PORT").(string)
		MailUser       = viper.Get("MAIL_USER").(string)
		MailPass       = viper.Get("MAIL_PASSWORD").(string)
		JwtSecret      = viper.Get("JWT_SECRET").(string)
		JwtExpiration  = viper.Get("JWT_EXPIRATION").(string)
	)

	logger := infrastructure.NewLogger(AppEnvironment)
	utils.NewJwt(JwtSecret, JwtExpiration)

	email := infrastructure.NewEmail(MailHost, MailPort, MailUser, MailPass, logger)
	err := email.Send("internal/templates/email_template.html", "asesinblood@gmail.com", []string{"asesinblood@gmail.com"}, "Test", []string{"Test", "Code"}, []string{})
	if err != nil {
		logger.Log.Sugar().Errorf("Failed to send email: %v", err)
	}

	routes := infrastructure.NewRoutes()
	infrastructure.NewHttpServer("localhost", AppPort, routes.Handlers, logger)
}
