package model

type User struct {
	Id           int    `db:"id" json:"id"`
	ExternalId   string `db:"external_id" json:"externalId"`
	Name         string `db:"name" json:"name"`
	Password     string `db:"password" json:"password"`
	TokenVersion int    `db:"token_version" json:"tokenVersion"`
}
