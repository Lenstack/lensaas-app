package infrastructure

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Postgres struct {
	Database squirrel.StatementBuilderType
}

func NewPostgres(host, port, user, password, dbName string, logger *zap.Logger) *Postgres {
	datasource := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Bogota",
		host, user, password, dbName, port)

	database, err := sql.Open("postgres", datasource)
	if err != nil {
		logger.Sugar().Error(err)
		return nil
	}

	err = database.Ping()
	if err != nil {
		logger.Sugar().Error(err)
		return nil
	}

	logger.Sugar().Info("Database connection successful")
	return &Postgres{Database: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(database)}
}
