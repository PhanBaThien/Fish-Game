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

// InitializeApp là hàm mẫu, Wire sẽ đọc hàm này và đẻ ra code thật
func InitializeApp(db *pgxpool.Pool, hasher utils.PasswordHasher, tokenMaker utils.TokenMaker) (authHttp.Handlers, error) {
	wire.Build(
		// 1. Tiêm Repositories
		repository.NewUserRepository,

		// 2. Tiêm Usecases
		usecase.NewAuthUsecase,

		// 3. Tiêm Handlers
		authHttp.NewAuthHandler,

		// 4. Tự động nhét tất cả Handlers vào struct authHttp.Handlers
		wire.Struct(new(authHttp.Handlers), "*"),
	)

	// Giá trị return này chỉ là dummy (giả), Wire sẽ tự ghi đè
	return authHttp.Handlers{}, nil
}