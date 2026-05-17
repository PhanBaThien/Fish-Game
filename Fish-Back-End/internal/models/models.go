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
	ID          int64              `db:"id"          json:"id"`
	Name        string             `db:"name"        json:"name"`
	MinBet      int64              `db:"min_bet"     json:"min_bet"`
	MaxPlayers  int32              `db:"max_players" json:"max_players"`
	Description pgtype.Text        `db:"description" json:"description"`
	CreatedAt   pgtype.Timestamptz `db:"created_at"  json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `db:"updated_at"  json:"updated_at"`
}

type Fish struct {
	ID               int32              `db:"id"                json:"id"`
	Name             string             `db:"name"              json:"name"`
	Health           int32              `db:"health"            json:"health"`
	RewardMultiplier int32              `db:"reward_multiplier" json:"reward_multiplier"`
	Speed            float64            `db:"speed"             json:"speed"`
	AssetPath        string             `db:"asset_path"        json:"asset_path"`
	CreatedAt        pgtype.Timestamptz `db:"created_at"        json:"created_at"`
	UpdatedAt        pgtype.Timestamptz `db:"updated_at"        json:"updated_at"`
}