package psql

import (
	"go-monorepo/internal-backend/persistence/db"
	"go-monorepo/internal/model"

	"github.com/jmoiron/sqlx"
)

type mineralStore struct {
	*sqlx.DB
}

func NewMineralStore() db.MineralStore {
	return &mineralStore{
		DB: GetDB(),
	}
}

func (m mineralStore) GetByName(name string) (model.Mineral, error) {
	var mineral model.Mineral

	if err := m.Get(&mineral, "SELECT * FROM mineral WHERE name = $1", name); err != nil {
		return model.Mineral{}, err
	}

	return mineral, nil
}

func (m mineralStore) UpdateAmount(name string, amount int) error {
	_, err := m.Exec("UPDATE mineral SET amount = $1 WHERE name = $2", amount, name)

	return err
}
