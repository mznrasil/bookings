package dbrepo

import (
	"database/sql"

	"github.com/mznrasil/bookings/internal/config"
	"github.com/mznrasil/bookings/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(appConfig *config.AppConfig, db *sql.DB) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: appConfig,
		DB:  db,
	}
}
