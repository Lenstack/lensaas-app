package services

import (
	"github.com/Lenstack/lensaas-app/internal/core/models"
	"github.com/Lenstack/lensaas-app/internal/core/repositories"
	"github.com/Lenstack/lensaas-app/internal/utils"
	"github.com/Masterminds/squirrel"
)

type IUserService interface {
	SignIn(email string, password string) (string, error)
	SignUp(user models.User) (string, error)
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
func (us *UserService) SignIn(email string, password string) (string, error) {
	return "", nil
}

// SignUp TODO: 1. Check if user already exists, 2. If user does not exist, create user, 3. Send email to user, 4. Return success message
func (us *UserService) SignUp(user models.User) (string, error) {
	return "", nil
}
