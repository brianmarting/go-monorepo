package psql

import (
	"fmt"
	"os"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	once       = sync.Once{}
	dbInstance *sqlx.DB
	err        error
)

func GetDB() *sqlx.DB {
	once.Do(func() {
		var (
			username = os.Getenv("DB_USERNAME")
			password = os.Getenv("DB_PASSWORD")
			host     = os.Getenv("DB_HOST")
			port     = os.Getenv("DB_PORT")
		)

		dbInstance, err = sqlx.Open(
			"postgres",
			fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable", username, password, host, port),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to open connection to db")
		}
	})
	return dbInstance
}
