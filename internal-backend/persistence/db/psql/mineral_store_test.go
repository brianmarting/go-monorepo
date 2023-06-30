package psql

import (
	"errors"
	"go-monorepo/internal/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

var mineral = model.Mineral{
	Id:     43,
	Name:   "coal",
	Amount: 33,
}

func Test_mineralStore_GetByName(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	store := mineralStore{
		DB: sqlxDB,
	}

	tests := []struct {
		name    string
		args    model.Mineral
		mock    func()
		wantErr bool
	}{
		{
			name: "Should get by name",
			args: mineral,
			mock: func() {
				rows := sqlmock.
					NewRows([]string{"id", "name", "amount"}).
					AddRow(mineral.Id, mineral.Name, mineral.Amount)
				mock.ExpectQuery("SELECT (.+) FROM mineral WHERE name = (.+)").
					WithArgs(mineral.Name).
					WillReturnRows(rows)
			},
		},
		{
			name:    "Should return err",
			args:    mineral,
			wantErr: true,
			mock: func() {
				mock.ExpectQuery("SELECT (.+) FROM mineral WHERE name = (.+)").
					WithArgs(mineral.Name).
					WillReturnError(errors.New("failed to exec query"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			result, err := store.GetByName(tt.args.Name)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && result != mineral {
				t.Errorf("Bodies dont match")
			}
		})
	}
}

func Test_mineralStore_Update(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	store := mineralStore{
		DB: sqlxDB,
	}

	tests := []struct {
		name    string
		args    model.Mineral
		mock    func()
		wantErr bool
	}{
		{
			name: "Should update mineral",
			args: mineral,
			mock: func() {
				mock.ExpectExec("UPDATE mineral SET amount = (.+) WHERE name = (.+)").
					WithArgs(mineral.Amount, mineral.Name).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "Should return err",
			args: mineral,
			mock: func() {
				mock.ExpectExec("UPDATE mineral SET amount = (.+) WHERE name = (.+)").
					WithArgs(mineral.Amount, mineral.Name).
					WillReturnError(errors.New("failed to exec query"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			if err := store.UpdateAmount(tt.args.Name, tt.args.Amount); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
