package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        int64     `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"`
	RoleID    int32     `db:"role_id" json:"role_id"`
	CreatedAt pgtype.Timestamptz `db:"created_at" json:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at" json:"updated_at"`
}

type Role struct {
	ID       int32  `db:"id" json:"id"`
	RoleName string `db:"role_name" json:"role_name"`
}

type Room struct {
    ID         int64
    Name       string
    MinBet     int64
    MaxPlayers int32
    Description pgtype.Text
    CreatedAt  pgtype.Timestamptz
    UpdatedAt  pgtype.Timestamptz
}

type Fish struct {
    ID               int64
    Name             string
    Health           int32
    RewardMultiplier int32
    Speed            float64
    AssetPath        string
    CreatedAt        pgtype.Timestamptz
    UpdatedAt        pgtype.Timestamptz
}