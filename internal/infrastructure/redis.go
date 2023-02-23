package infrastructure

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis(host, port, password string, db string, logger *zap.Logger) *Redis {
	dbName, err := strconv.Atoi(db)
	if err != nil {
		logger.Sugar().Error(err)
		return nil
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       dbName,
	})

	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.Sugar().Error(err)
		return nil
	}
	logger.Sugar().Info("Redis connection successful")
	return &Redis{Client: redisClient}
}
