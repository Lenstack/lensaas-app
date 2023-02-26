package services

import (
	"errors"
	"fmt"
	"github.com/Lenstack/lensaas-app/internal/core/entities"
	"github.com/Lenstack/lensaas-app/internal/core/repositories"
	"github.com/Lenstack/lensaas-app/internal/templates"
	"github.com/Lenstack/lensaas-app/internal/utils"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

type IUserService interface {
	SignIn(email string, password string) (token string, refreshToken string, err error)
	SignUp(user entities.User) (message string, err error)
	SignOut(token string) (message string, err error)
	SendVerificationCode(name, email string) (message string, err error)
	SendVerificationEmail(name, email string) (message string, err error)
	VerifyEmail(token string) (message string, err error)
	VerifyCode(email string, code string) (message string, err error)
	RefreshToken(refreshToken string) (token string, err error)
}

type UserService struct {
	UserRepository repositories.UserRepository
	TokenService   TokenService
	EmailService   EmailService
	bcrypt         *utils.Bcrypt
}

func NewUserService(database squirrel.StatementBuilderType, redis *redis.Client, tokenService TokenService, emailService EmailService) *UserService {
	return &UserService{
		UserRepository: repositories.UserRepository{
			Database: database,
			Redis:    redis,
		},
		TokenService: tokenService,
		EmailService: emailService,
	}
}

// SignIn TODO: 1. Check if user exists, 2. If user exists, check if password is correct, 3. If password is correct, generate token, 4. Return token
func (us *UserService) SignIn(email string, password string) (token string, refreshToken string, expiresIn int64, err error) {
	user, err := us.UserRepository.FindByEmail(email)
	if err != nil {
		return "", "", 0, err
	}

	if !user.Verified {
		return "", "", 0, errors.New("user is not verified")
	}

	err = us.bcrypt.ComparePassword(user.Password, password)
	if err != nil {
		return "", "", 0, err
	}

	token, err = us.TokenService.GenerateToken(user.Id, us.TokenService.ExpirationTime)
	if err != nil {
		return "", "", 0, err
	}

	refreshToken, err = us.TokenService.NewRefreshToken()
	if err != nil {
		return "", "", 0, err
	}

	_, err = us.UserRepository.UpdateRefreshToken(user.Id, refreshToken)
	if err != nil {
		return "", "", 0, err
	}

	return token, refreshToken, time.Now().Add(us.TokenService.ExpirationTime).Unix(), nil
}

// SignUp TODO: 1. Check if user already exists, 2. If user does not exist, create user, 3. Send email to user, 4. Return success message
func (us *UserService) SignUp(user entities.User) (message string, err error) {
	isFounded, err := us.UserRepository.FindByEmail(user.Email)
	if isFounded.Id != "" {
		return "", errors.New("user already exists")
	}

	hashedPassword, err := us.bcrypt.HashPassword(user.Password)
	if err != nil {
		return "", err
	}

	newUser := entities.User{
		Id:       uuid.New().String(),
		Email:    user.Email,
		Name:     user.Name,
		Password: hashedPassword,
	}

	userId, err := us.UserRepository.Create(newUser)
	if err != nil {
		return "", err
	}

	_, err = us.SendVerificationEmail(user.Name, user.Email)
	if err != nil {
		return "", err
	}

	return userId, nil
}

// SignOut TODO: 1. Check if user exists, 2. If user exists, delete token, 3. Return success message
func (us *UserService) SignOut(token string) (string, error) {
	userId, err := us.TokenService.ValidateToken(token)
	if err != nil {
		return "", err
	}

	fmt.Println("revoke token for user: ", userId)
	return "", nil
}

// SendVerificationCode TODO: 1. Generate verification code, 2. Send email to user, 3. Return success message
func (us *UserService) SendVerificationCode(name, email string) (message string, err error) {
	code := utils.NewCode()
	sendExpiresAt := time.Now().Add(time.Minute * 5)

	message, err = us.UserRepository.UpdateVerificationCode(email, code, sendExpiresAt)
	if err != nil {
		return "", err
	}

	mail, err := us.EmailService.Create("internal/templates/verification_code_template.html", []string{email},
		"Verification Code", templates.VerificationCode{Name: strings.ToTitle(name), Code: code}, []string{})
	if err != nil {
		return "", err
	}

	err = us.EmailService.Send(mail)
	if err != nil {
		return "", err
	}

	return message, nil
}

// SendVerificationEmail TODO: 1. Generate verification token, 2. Send email to user, 3. Return success message
func (us *UserService) SendVerificationEmail(name, email string) (message string, err error) {
	inFiveMinutes := time.Now().Add(time.Minute * 5)
	expiration := inFiveMinutes.Sub(time.Now())

	user, err := us.UserRepository.FindByEmail(email)
	if err != nil {
		return "", err
	}

	token, err := us.TokenService.GenerateToken(user.Id, expiration)
	if err != nil {
		return "", err
	}

	mail, err := us.EmailService.Create("internal/templates/verification_email_template.html", []string{email},
		"Verification Email", templates.VerificationEmail{Name: strings.ToTitle(name), Token: token}, []string{})
	if err != nil {
		return "", err
	}

	err = us.EmailService.Send(mail)
	if err != nil {
		return "", err
	}

	return "success", nil
}

// VerifyEmail TODO: 1. Get token from request, 2. Validate token, 3. Call EmailVerification method from UserService, 4. Return success message
func (us *UserService) VerifyEmail(token string) (message string, err error) {
	userId, err := us.TokenService.ValidateToken(token)
	if err != nil {
		return "", err
	}

	user, err := us.UserRepository.FindById(userId)
	if err != nil {
		return "", err
	}

	if user.Verified {
		return "", errors.New("user is already verified")
	}

	_, err = us.UserRepository.UpdateVerified(user.Email, true)
	if err != nil {
		return "", err
	}

	return "your email has been verified successfully", nil
}

// VerifyCode TODO: 1. Get code from request, 2. Validate code, 3. Call EmailVerification method from UserService, 4. Return success message
func (us *UserService) VerifyCode(email string, code string) (message string, err error) {
	return "", errors.New("not implemented")
}

// RefreshToken TODO: 1. Get refresh token from request, 2. Validate refresh token, 3. Generate new token, 4. Return new token
func (us *UserService) RefreshToken(refreshToken string) (token string, expiresIn int64, err error) {
	user, err := us.UserRepository.FindByRefreshToken(refreshToken)
	if err != nil {
		return "", 0, errors.New("invalid refresh token")
	}

	token, err = us.TokenService.GenerateToken(user.Id, us.TokenService.ExpirationTime)
	if err != nil {
		return "", 0, err
	}

	message, err := us.UserRepository.SaveAccessToken(user.Id, token, us.TokenService.ExpirationTime)
	if err != nil {
		return "", 0, err
	}
	fmt.Println(message)

	return token, time.Now().Add(us.TokenService.ExpirationTime).Unix(), nil
}
