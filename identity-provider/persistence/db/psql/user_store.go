package psql

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"go-monorepo/identity-provider/persistence/db"
	"go-monorepo/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type userStore struct {
	*sqlx.DB
}

func NewUserStore() db.UserStore {
	return &userStore{
		DB: GetDB(),
	}
}

func (u userStore) GetByExternalId(id string) (model.User, error) {
	var user model.User

	if err := u.Get(&user, "SELECT * FROM users WHERE external_id = $1", id); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u userStore) Create(user model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	log.Info().Msg(fmt.Sprintf("hash pwd: %s", hashedPassword))
	if err != nil {
		return err
	}

	_, err = u.Exec(
		"INSERT INTO users VALUES ($1, $2, $3, $4)",
		uuid.New().String(),
		user.Name,
		hashedPassword,
		1,
	)

	return err
}
