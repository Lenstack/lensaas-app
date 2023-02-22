package repositories

import (
	"fmt"
	"github.com/Lenstack/lensaas-app/internal/core/models"
	"github.com/Masterminds/squirrel"
)

type IUserRepository interface {
	Create(user models.User) (string, error)
	FindByEmail(email string) (models.User, error)
}

type UserRepository struct {
	Database squirrel.StatementBuilderType
}

// Create TODO: 1. Create user, 2. Return user id
func (ur *UserRepository) Create(user models.User) (string, error) {
	_, err := ur.Database.Insert(models.UserTableName).
		Columns("Id", "Name", "Email", "Password", "Verified", "Code", "SendExpiresAt").
		Values(user.Id, user.Name, user.Email, user.Password, user.Verified, user.Code, user.SendExpiresAt).
		Exec()
	if err != nil {
		return "", err
	}
	return user.Id, nil
}

// FindByEmail TODO: 1. Find user by email, 2. Return user
func (ur *UserRepository) FindByEmail(email string) (user models.User, err error) {
	err = ur.Database.Select("*").
		From(models.UserTableName).
		Where(squirrel.Eq{"email": email}).
		QueryRow().
		Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		fmt.Println(err)
		return models.User{}, err
	}
	fmt.Println(user)
	return user, nil
}
