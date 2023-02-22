package repositories

import (
	"github.com/Lenstack/lensaas-app/internal/core/entities"
	"github.com/Masterminds/squirrel"
)

type IUserRepository interface {
	Create(user entities.User) (userId string, err error)
	FindByEmail(email string) (user entities.User, err error)
}

type UserRepository struct {
	Database squirrel.StatementBuilderType
}

// Create TODO: 1. Create user, 2. Return user id
func (ur *UserRepository) Create(user entities.User) (userId string, err error) {
	qb := ur.Database.Insert(entities.UserTableName).
		Columns("Id", "Name", "Email", "Password", "Verified", "Code", "SendExpiresAt").
		Values(user.Id, user.Name, user.Email, user.Password, user.Verified, user.Code, user.SendExpiresAt).
		Suffix("RETURNING Id")
	err = qb.QueryRow().Scan(&userId)
	if err != nil {
		return "", err
	}
	return userId, nil
}

// FindByEmail TODO: 1. Find user by email, 2. Return user
func (ur *UserRepository) FindByEmail(email string) (user entities.User, err error) {
	err = ur.Database.Select("Id", "Name", "Email", "Password", "Verified", "Code", "SendExpiresAt", "CreatedAt", "UpdatedAt").
		From(entities.UserTableName).
		Where(squirrel.Eq{"email": email}).
		QueryRow().
		Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Verified, &user.Code, &user.SendExpiresAt, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}
