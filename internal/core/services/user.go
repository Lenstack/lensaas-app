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
	SignIn(email string, password string) (accessToken string, refreshToken string, expiresIn time.Duration, err error)
	SignUp(user entities.User) (message string, err error)
	SignOut(token string) (message string, err error)
	SendVerificationCode(name, email string) (message string, err error)
	SendVerificationEmail(name, email string) (message string, err error)
	VerifyEmail(token string) (message string, err error)
	VerifyCode(token string, code string) (message string, err error)
	RefreshToken(refreshToken string) (token string, expiresIn time.Duration, err error)
	RevokeToken(token string) (message string, err error)
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
func (us *UserService) SignIn(email string, password string) (accessToken string, refreshToken string, expiresIn time.Duration, err error) {
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

	accessToken, err = us.TokenService.GenerateToken(user.Id, us.TokenService.ExpirationTimeAccess)
	if err != nil {
		return "", "", 0, err
	}

	refreshToken, err = us.TokenService.GenerateToken(user.Id, us.TokenService.ExpirationTimeRefresh)
	if err != nil {
		return "", "", 0, err
	}

	_, err = us.UserRepository.SaveRefreshToken(user.Id, refreshToken, us.TokenService.ExpirationTimeRefresh)
	if err != nil {
		return "", "", 0, err
	}

	return accessToken, refreshToken, us.TokenService.ExpirationTimeAccess, nil
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
func (us *UserService) SignOut(token string) (message string, err error) {
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
		return "", errors.New("invalid token")
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
func (us *UserService) VerifyCode(token string, code string) (message string, err error) {
	userId, err := us.TokenService.ValidateToken(token)
	if err != nil {
		return "", errors.New("invalid token")
	}

	user, err := us.UserRepository.FindById(userId)
	if err != nil {
		return "", err
	}

	if user.Code != code {
		return "", errors.New("invalid code")
	}

	if user.SendExpiresAt.Before(time.Now()) {
		return "", errors.New("code has been expired")
	}

	// TODO: 1. Authorize user, 2. , 3. Return success message
	return "your email has been verified successfully", nil
}

// RefreshToken TODO: 1. Get refresh token from request, 2. Validate refresh token, 3. Generate new token, 4. Return new token
func (us *UserService) RefreshToken(refreshToken string) (token string, expiresIn time.Duration, err error) {
	userId, err := us.TokenService.ValidateToken(refreshToken)
	if err != nil {
		return "", 0, errors.New("invalid refresh token")
	}

	refreshTokenList, err := us.UserRepository.FindRefreshToken(userId)
	if err != nil {
		return "", 0, err
	}

	var tokenExist *entities.TokenList
	for _, tok := range refreshTokenList {
		if tok.Token == refreshToken {
			tokenExist = &tok
			if tok.Blocked {
				return "", 0, errors.New("refresh token is blocked")
			}
			break
		}
	}

	if tokenExist == nil {
		return "", 0, errors.New("invalid refresh token")
	}

	token, err = us.TokenService.GenerateToken(userId, us.TokenService.ExpirationTimeAccess)
	if err != nil {
		return "", 0, err
	}

	return token, us.TokenService.ExpirationTimeAccess, nil
}

// RevokeToken TODO: 1. Get refresh token from request, 2. Validate refresh token, 3. Block refresh token, 4. Return success message
func (us *UserService) RevokeToken(refreshToken string) (message string, err error) {
	userId, err := us.TokenService.ValidateToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	message, err = us.UserRepository.BlockRefreshToken(userId, refreshToken)
	if err != nil {
		return "", err
	}

	return message, nil
}
