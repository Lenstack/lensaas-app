package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/Lenstack/lensaas-app/internal/core/entities"
	"github.com/Masterminds/squirrel"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type IUserRepository interface {
	Create(user entities.User) (userId string, err error)
	FindById(userId string) (user entities.User, err error)
	FindByEmail(email string) (user entities.User, err error)
	FindByRefreshToken(refreshToken string) (user entities.User, err error)
	UpdateVerified(email string, verified bool) (message string, err error)
	UpdateVerificationCode(email string, code string, sendExpiresAt time.Time) (message string, err error)

	FindRefreshToken(userId string) (refreshToken []entities.TokenList, err error)
	SaveRefreshToken(userId string, refreshToken string, expiresIn time.Duration) (message string, err error)
	DeleteRefreshToken(tokenId string) (message string, err error)
	BlockRefreshToken(tokenId string) (message string, err error)
}

type UserRepository struct {
	Database squirrel.StatementBuilderType
	Redis    *redis.Client
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

// FindByRefreshToken TODO: 1. Find user by refresh token, 2. Return user
func (ur *UserRepository) FindByRefreshToken(refreshToken string) (user entities.User, err error) {
	err = ur.Database.Select("Id", "Name", "Email", "Password",
		"Verified", "Code", "Token", "SendExpiresAt", "CreatedAt", "UpdatedAt").
		From(entities.UserTableName).
		Where(squirrel.Eq{"token": refreshToken}).
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

// FindRefreshToken TODO: 1. Find refresh token by user id, 2. Return refresh token
func (ur *UserRepository) FindRefreshToken(userId string) (refreshToken []entities.TokenList, err error) {
	userKey := fmt.Sprintf("refresh_token:%s", userId)
	keys, err := ur.Redis.SMembers(context.Background(), userKey).Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		tokenData, err := ur.Redis.HGetAll(context.Background(), key).Result()
		if err != nil {
			return nil, err
		}

		// Convert string to bool and int64
		blocked, err := strconv.ParseBool(tokenData["Blocked"])
		if err != nil {
			return nil, err
		}
		expiration, err := strconv.ParseInt(tokenData["Expiration"], 10, 64)
		if err != nil {
			return nil, err
		}

		tokenList := &entities.TokenList{
			Type:       tokenData["Type"],
			Token:      tokenData["Token"],
			UserId:     tokenData["UserId"],
			Blocked:    blocked,
			Expiration: expiration,
		}

		refreshToken = append(refreshToken, *tokenList)
	}
	return refreshToken, nil
}

// SaveRefreshToken TODO: 1. Save refresh token to redis, 2. Return success message
func (ur *UserRepository) SaveRefreshToken(userId string, refreshToken string, expiresIn time.Duration) (message string, err error) {
	key := fmt.Sprintf("refresh_token:%s:%s", userId, refreshToken)
	expiresInTime := time.Now().Add(expiresIn).Unix()

	tokenList := &entities.TokenList{
		Type:       "Refresh_Token",
		Token:      refreshToken,
		UserId:     userId,
		Blocked:    false,
		Expiration: expiresInTime,
	}

	err = ur.Redis.HMSet(context.Background(), key, map[string]interface{}{
		"Type":       tokenList.Type,
		"Token":      tokenList.Token,
		"UserId":     tokenList.UserId,
		"Blocked":    tokenList.Blocked,
		"Expiration": tokenList.Expiration,
	}).Err()
	if err != nil {
		return "", err
	}

	err = ur.Redis.Expire(context.Background(), key, expiresIn).Err()
	if err != nil {
		return "", err
	}

	userKey := fmt.Sprintf("refresh_token:%s", userId)
	err = ur.Redis.SAdd(context.Background(), userKey, key).Err()
	if err != nil {
		return "", err
	}

	return "success", nil
}

// DeleteRefreshToken TODO: 1. Delete refresh token from redis, 2. Return success message
func (ur *UserRepository) DeleteRefreshToken(tokenId string) (message string, err error) {
	values, err := ur.Redis.HMGet(context.Background(), tokenId, "user_id", "refresh_token").Result()
	if err != nil {
		return "", err
	}

	userId, ok1 := values[0].(string)
	refreshToken, ok2 := values[1].(string)
	if !ok1 || !ok2 {
		return "", errors.New("failed to get user_id or refresh_token from Redis")
	}

	deletedCount, err := ur.Redis.Del(context.Background(), tokenId).Result()
	if err != nil {
		return "", err
	}
	if deletedCount != 1 {
		return "", errors.New("failed to delete refresh token from Redis")
	}

	tokensKey := fmt.Sprintf("user_refresh_tokens:%s", userId)
	deletedCount, err = ur.Redis.ZRem(context.Background(), tokensKey, refreshToken).Result()
	if err != nil {
		return "", err
	}

	if deletedCount != 1 {
		return "", errors.New("failed to delete refresh token from user's refresh token list")
	}
	return "success", nil
}

// BlockRefreshToken TODO: 1. Block refresh token from redis, 2. Return success message
func (ur *UserRepository) BlockRefreshToken(userId, tokenId string) (message string, err error) {
	userKey := fmt.Sprintf("refresh_token:%s", userId)

	// Get all refresh tokens for the user.
	keys, err := ur.Redis.SMembers(context.Background(), userKey).Result()
	if err != nil {
		return "", err
	}

	// Check each refresh token for a match with the provided token ID.
	for _, key := range keys {
		tokenData, err := ur.Redis.HGetAll(context.Background(), key).Result()
		if err != nil {
			return "", err
		}

		// tokenData["Blocked"] to bool type and check if the token is already blocked. If so, return an error.
		blocked, err := strconv.ParseBool(tokenData["Blocked"])
		if err != nil {
			return "", err
		}

		if blocked {
			return "", errors.New("token already blocked")
		}

		// If the token ID matches, delete the token and return a success message.
		if tokenData["Token"] == tokenId {
			// Update the token's "blocked" status to true.
			err = ur.Redis.HSet(context.Background(), key, "Blocked", true).Err()
			if err != nil {
				return "", err
			}
			return "refresh token blocked successfully", nil
		}
	}

	// If no matching token ID is found, return an error.
	return "", errors.New("token not found")
}
