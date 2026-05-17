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
		// Repositories
		repository.NewUserRepository,
		repository.NewRoomRepository,
		repository.NewFishRepository,

		// Usecases
		usecase.NewAuthUsecase,
		usecase.NewRoomUsecase,
		usecase.NewFishUsecase,

		// Handlers
		authHttp.NewAuthHandler,
		authHttp.NewRoomHandler,
		authHttp.NewFishHandler,

		wire.Struct(new(authHttp.Handlers), "*"),
	)

	return authHttp.Handlers{}, nil
}
