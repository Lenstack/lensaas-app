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
	"strings"
	"time"
)

type IUserService interface {
	SignIn(email string, password string) (token string, err error)
	SignUp(user entities.User) (message string, err error)
	SignOut(token string) (message string, err error)
	SendVerificationEmail(name, email string) (message string, err error)
}

type UserService struct {
	UserRepository repositories.UserRepository
	TokenService   TokenService
	EmailService   EmailService
	bcrypt         *utils.Bcrypt
}

func NewUserService(database squirrel.StatementBuilderType, emailService EmailService) *UserService {
	return &UserService{
		UserRepository: repositories.UserRepository{
			Database: database,
		},
		EmailService: emailService,
	}
}

// SignIn TODO: 1. Check if user exists, 2. If user exists, check if password is correct, 3. If password is correct, generate token, 4. Return token
func (us *UserService) SignIn(email string, password string) (token string, err error) {
	user, err := us.UserRepository.FindByEmail(email)
	if err != nil {
		return "", err
	}

	if !user.Verified {
		return "", errors.New("user is not verified")
	}

	err = us.bcrypt.ComparePassword(user.Password, password)
	if err != nil {
		return "", err
	}

	token, err = us.TokenService.GenerateToken(user.Id, us.TokenService.ExpirationTime)
	if err != nil {
		return "", err
	}

	refreshToken, err := us.TokenService.NewRefreshToken()
	if err != nil {
		return "", err
	}
	fmt.Printf("refreshToken: %s", refreshToken)
	return token, nil
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

// SendVerificationEmail TODO: 1. Generate verification code, 2. Send email to user, 3. Return success message
func (us *UserService) SendVerificationEmail(name, email string) (message string, err error) {
	code := utils.NewCode()
	sendExpiresAt := time.Now().Add(time.Minute * 5)

	message, err = us.UserRepository.UpdateVerificationCode(email, code, sendExpiresAt)
	if err != nil {
		return "", err
	}

	mail, err := us.EmailService.Create("internal/templates/verification_template.html", []string{email},
		"Verification Code", templates.Verification{Name: strings.ToTitle(name), Code: code}, []string{})
	if err != nil {
		return "", err
	}

	err = us.EmailService.Send(mail)
	if err != nil {
		return "", err
	}

	return message, nil
}
