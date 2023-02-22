package infrastructure

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

type Postgres struct {
	Database squirrel.StatementBuilderType
}

func NewPostgres(host, port, user, password, dbName string, logger *Logger) *Postgres {
	datasource := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Bogota",
		host, user, password, dbName, port)

	database, err := sql.Open("postgres", datasource)
	if err != nil {
		logger.Log.Sugar().Error(err)
		return nil
	}

	err = database.Ping()
	if err != nil {
		logger.Log.Sugar().Error(err)
		return nil
	}

	logger.Log.Sugar().Info("Database connection successful")
	return &Postgres{Database: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(database)}
}
