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
			username = os.Getenv("IPS_USERNAME")
			password = os.Getenv("IPS_PASSWORD")
			host     = os.Getenv("IPS_HOST")
		)

		dbInstance, err = sqlx.Open(
			"postgres",
			fmt.Sprintf("postgres://%s:%s@%s:5432/postgres?sslmode=disable", username, password, host),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to open connection to db")
		}
	})
	return dbInstance
}
