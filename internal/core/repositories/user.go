package repositories

import (
	"github.com/Lenstack/lensaas-app/internal/core/entities"
	"github.com/Masterminds/squirrel"
	"time"
)

type IUserRepository interface {
	Create(user entities.User) (userId string, err error)
	FindById(userId string) (user entities.User, err error)
	FindByEmail(email string) (user entities.User, err error)
	UpdateVerified(email string, verified bool) (message string, err error)
	UpdateVerificationCode(email string, code string, sendExpiresAt time.Time) (message string, err error)
	UpdateRefreshToken(userId string, refreshToken string) (message string, err error)
}

type UserRepository struct {
	Database squirrel.StatementBuilderType
}

// Create TODO: 1. Create user, 2. Return user id
func (ur *UserRepository) Create(user entities.User) (userId string, err error) {
	qb := ur.Database.Insert(entities.UserTableName).
		Columns("Id", "Name", "Email", "Password", "Verified", "Code", "Token", "SendExpiresAt").
		Values(user.Id, user.Name, user.Email, user.Password, user.Verified, user.Code, user.Token, user.SendExpiresAt).
		Suffix("RETURNING Id")
	err = qb.QueryRow().Scan(&userId)
	if err != nil {
		return "", err
	}
	return userId, nil
}

// FindById TODO: 1. Find user by id, 2. Return user
func (ur *UserRepository) FindById(userId string) (user entities.User, err error) {
	err = ur.Database.Select("Id", "Name", "Email", "Password",
		"Verified", "Code", "Token", "SendExpiresAt", "CreatedAt", "UpdatedAt").
		From(entities.UserTableName).
		Where(squirrel.Eq{"id": userId}).
		QueryRow().
		Scan(&user.Id, &user.Name, &user.Email, &user.Password,
			&user.Verified, &user.Code, &user.Token, &user.SendExpiresAt,
			&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

// FindByEmail TODO: 1. Find user by email, 2. Return user
func (ur *UserRepository) FindByEmail(email string) (user entities.User, err error) {
	err = ur.Database.Select("Id", "Name", "Email", "Password",
		"Verified", "Code", "Token", "SendExpiresAt", "CreatedAt", "UpdatedAt").
		From(entities.UserTableName).
		Where(squirrel.Eq{"email": email}).
		QueryRow().
		Scan(&user.Id, &user.Name, &user.Email, &user.Password,
			&user.Verified, &user.Code, &user.Token, &user.SendExpiresAt,
			&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

// UpdateVerified TODO: 1. Update verified, 2. Return success message
func (ur *UserRepository) UpdateVerified(email string, verified bool) (message string, err error) {
	qb := ur.Database.Update(entities.UserTableName).
		Set("Verified", verified).
		Where(squirrel.Eq{"Email": email}).
		Suffix("RETURNING Id")

	err = qb.QueryRow().Scan(&message)
	if err != nil {
		return "", err
	}
	return message, nil
}

// UpdateVerificationCode TODO: 1. Update verification code,SendExpiresAt by email 2. Return success message
func (ur *UserRepository) UpdateVerificationCode(email string, code string, sendExpiresAt time.Time) (message string, err error) {
	qb := ur.Database.Update(entities.UserTableName).
		Set("Code", code).
		Set("SendExpiresAt", sendExpiresAt).
		Where(squirrel.Eq{"Email": email}).
		Suffix("RETURNING Id")

	err = qb.QueryRow().Scan(&message)
	if err != nil {
		return "", err
	}
	return message, nil
}

// UpdateRefreshToken TODO: 1. Update refresh token by user id, 2. Return success message
func (ur *UserRepository) UpdateRefreshToken(userId string, refreshToken string) (message string, err error) {
	qb := ur.Database.Update(entities.UserTableName).
		Set("Token", refreshToken).
		Where(squirrel.Eq{"Id": userId}).
		Suffix("RETURNING Id")

	err = qb.QueryRow().Scan(&message)
	if err != nil {
		return "", err
	}
	return message, nil
}
