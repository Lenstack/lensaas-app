package services

import (
	"errors"
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
	SignOut(userId string) (message string, err error)
	SendVerificationEmail(email string) (message string, err error)
}

type UserService struct {
	UserRepository repositories.UserRepository
	Jwt            *utils.Jwt
	Email          *utils.Email
	Bcrypt         *utils.Bcrypt
}

func NewUserService(database squirrel.StatementBuilderType, jwt *utils.Jwt, email *utils.Email) *UserService {
	return &UserService{
		UserRepository: repositories.UserRepository{
			Database: database,
		},
		Jwt:   jwt,
		Email: email,
	}
}

// SignIn TODO: 1. Check if user exists, 2. If user exists, check if password is correct, 3. If password is correct, generate token, 4. Return token
func (us *UserService) SignIn(email string, password string) (token string, err error) {
	return "", nil
}

// SignUp TODO: 1. Check if user already exists, 2. If user does not exist, create user, 3. Send email to user, 4. Return success message
func (us *UserService) SignUp(user entities.User) (message string, err error) {
	isFounded, err := us.UserRepository.FindByEmail(user.Email)
	if isFounded.Id != "" {
		return "", errors.New("user already exists")
	}

	hashedPassword, err := us.Bcrypt.HashPassword(user.Password)
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
func (us *UserService) SignOut(userId string) (string, error) {
	return "", nil
}

// SendVerificationEmail TODO: 1. Check if user exists, 2. If user exists, send email to user, 3. Return success message
func (us *UserService) SendVerificationEmail(name, email string) (message string, err error) {
	code := utils.NewCode()
	sendExpiresAt := time.Now().Add(time.Minute * 5)
	message, err = us.UserRepository.UpdateVerificationCode(email, code, sendExpiresAt)
	if err != nil {
		return "", err
	}

	err = us.Email.Send("internal/templates/verification_template.html", []string{email},
		"Verification Code", templates.Verification{Name: strings.ToTitle(name), Code: code}, []string{})
	if err != nil {
		return "", err
	}
	return message, nil
}
