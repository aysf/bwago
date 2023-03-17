package dbrepo

import (
	"database/sql"

	"github.com/aysf/bwago/internal/config"
	"github.com/aysf/bwago/internal/repository"
)

// PostgresDBRepo is type of postgresSQL connection
// if you want to add another, i.e. MariaDB, you can
// simply create an another one like MariaDBRepo
type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
