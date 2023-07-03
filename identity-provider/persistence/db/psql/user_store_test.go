package psql

import (
	"errors"
	"go-monorepo/internal/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

var user = model.User{
	Id:           32,
	ExternalId:   "dc9b5f7e-61e6-480e-91bb-269f7a929622",
	Name:         "John",
	Password:     "somepwd",
	TokenVersion: 1,
}

func Test_userStore_GetByExternalId(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	store := userStore{
		DB: sqlxDB,
	}
	tests := []struct {
		name    string
		arg     string
		mock    func()
		wantErr bool
	}{
		{
			name: "Should get by id",
			arg:  user.ExternalId,
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "external_id", "name", "password", "token_version"}).
					AddRow(user.Id, user.ExternalId, user.Name, user.Password, user.TokenVersion)
				mock.ExpectQuery("SELECT (.+) FROM users WHERE external_id = (.+)").
					WithArgs(user.ExternalId).
					WillReturnRows(rows)
			},
		},
		{
			name: "Should return err",
			arg:  user.ExternalId,
			mock: func() {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE external_id = (.+)").
					WithArgs(user.ExternalId).
					WillReturnError(errors.New("failed to exec query"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			result, err := store.GetByExternalId(tt.arg)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetByExternalId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && result != user {
				t.Errorf("Bodies dont match")
			}
		})
	}
}

func Test_userStore_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	store := userStore{
		DB: sqlxDB,
	}

	tests := []struct {
		name    string
		arg     model.User
		mock    func()
		wantErr bool
	}{
		{
			name: "Should create user",
			arg:  user,
			mock: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(sqlmock.AnyArg(), user.Name, sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "Should return err",
			arg:  model.User{},
			mock: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(sqlmock.AnyArg(), user.Name, sqlmock.AnyArg(), 1).
					WillReturnError(errors.New("failed to exec query"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			if err := store.Create(tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
