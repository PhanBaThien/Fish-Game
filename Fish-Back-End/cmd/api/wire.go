//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/repository"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/usecase"
	authHttp "github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/transport/http"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/utils"
)

func InitializeApp(db *pgxpool.Pool, hasher utils.PasswordHasher, tokenMaker utils.TokenMaker) (authHttp.Handlers, error) {
	wire.Build(
		repository.NewUserRepository,
		repository.NewRefreshTokenRepository,
		repository.NewRoomRepository,
		repository.NewFishRepository,

		usecase.NewAuthUsecase,
		usecase.NewRoomUsecase,
		usecase.NewFishUsecase,

		authHttp.NewAuthHandler,
		authHttp.NewRoomHandler,
		authHttp.NewFishHandler,

		wire.Struct(new(authHttp.Handlers), "*"),
	)
	return authHttp.Handlers{}, nil
}
