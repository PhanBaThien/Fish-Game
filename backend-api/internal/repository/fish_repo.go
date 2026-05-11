package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/yourname/fish-game-backend/internal/models"
	"github.com/yourname/fish-game-backend/pkg/database"
)

// FishRepository defines the interface for fish config data access.
type FishRepository interface {
	FindAll(ctx context.Context) ([]models.Fish, error)
	FindByID(ctx context.Context, id string) (*models.Fish, error)
	Create(ctx context.Context, req models.CreateFishRequest) (*models.Fish, error)
	Update(ctx context.Context, id string, req models.UpdateFishRequest) error
	Delete(ctx context.Context, id string) error
}

type fishSQLServerRepo struct {
	db *database.DB
}

// NewFishRepository creates a new SQL Server–backed FishRepository.
func NewFishRepository(db *database.DB) FishRepository {
	return &fishSQLServerRepo{db: db}
}

func (r *fishSQLServerRepo) FindAll(ctx context.Context) ([]models.Fish, error) {
	const query = `
		SELECT id, name, multiplier, base_prob, speed, role, is_active
		FROM fish_configs
		ORDER BY multiplier ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fishRepo.FindAll: %w", err)
	}
	defer rows.Close()

	var list []models.Fish
	for rows.Next() {
		var f models.Fish
		if err := rows.Scan(&f.ID, &f.Name, &f.Multiplier, &f.BaseProb, &f.Speed, &f.Role, &f.IsActive); err != nil {
			return nil, fmt.Errorf("fishRepo.FindAll scan: %w", err)
		}
		list = append(list, f)
	}
	return list, rows.Err()
}

func (r *fishSQLServerRepo) FindByID(ctx context.Context, id string) (*models.Fish, error) {
	const query = `
		SELECT id, name, multiplier, base_prob, speed, role, is_active
		FROM fish_configs WHERE id = @p1`

	var f models.Fish
	err := r.db.QueryRowContext(ctx, query, sql.Named("p1", id)).
		Scan(&f.ID, &f.Name, &f.Multiplier, &f.BaseProb, &f.Speed, &f.Role, &f.IsActive)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("fishRepo.FindByID: fish %q not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("fishRepo.FindByID: %w", err)
	}
	return &f, nil
}

func (r *fishSQLServerRepo) Create(ctx context.Context, req models.CreateFishRequest) (*models.Fish, error) {
	const query = `
		INSERT INTO fish_configs (name, multiplier, base_prob, speed, role, is_active)
		OUTPUT INSERTED.id, INSERTED.name, INSERTED.multiplier, INSERTED.base_prob,
		       INSERTED.speed, INSERTED.role, INSERTED.is_active
		VALUES (@p1, @p2, @p3, @p4, @p5, 1)`

	var f models.Fish
	err := r.db.QueryRowContext(ctx, query,
		sql.Named("p1", req.Name),
		sql.Named("p2", req.Multiplier),
		sql.Named("p3", req.BaseProb),
		sql.Named("p4", req.Speed),
		sql.Named("p5", req.Role),
	).Scan(&f.ID, &f.Name, &f.Multiplier, &f.BaseProb, &f.Speed, &f.Role, &f.IsActive)
	if err != nil {
		return nil, fmt.Errorf("fishRepo.Create: %w", err)
	}
	return &f, nil
}

func (r *fishSQLServerRepo) Update(ctx context.Context, id string, req models.UpdateFishRequest) error {
	setClauses := []string{}
	args := []interface{}{}
	paramIdx := 1

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = @p%d", paramIdx))
		args = append(args, sql.Named(fmt.Sprintf("p%d", paramIdx), *req.Name))
		paramIdx++
	}
	if req.Multiplier != nil {
		setClauses = append(setClauses, fmt.Sprintf("multiplier = @p%d", paramIdx))
		args = append(args, sql.Named(fmt.Sprintf("p%d", paramIdx), *req.Multiplier))
		paramIdx++
	}
	if req.BaseProb != nil {
		setClauses = append(setClauses, fmt.Sprintf("base_prob = @p%d", paramIdx))
		args = append(args, sql.Named(fmt.Sprintf("p%d", paramIdx), *req.BaseProb))
		paramIdx++
	}
	if req.Speed != nil {
		setClauses = append(setClauses, fmt.Sprintf("speed = @p%d", paramIdx))
		args = append(args, sql.Named(fmt.Sprintf("p%d", paramIdx), *req.Speed))
		paramIdx++
	}
	if req.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = @p%d", paramIdx))
		args = append(args, sql.Named(fmt.Sprintf("p%d", paramIdx), *req.IsActive))
		paramIdx++
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, sql.Named(fmt.Sprintf("p%d", paramIdx), id))
	query := fmt.Sprintf(
		"UPDATE fish_configs SET %s WHERE id = @p%d",
		strings.Join(setClauses, ", "), paramIdx,
	)
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("fishRepo.Update: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("fishRepo.Update: fish %q not found", id)
	}
	return nil
}

func (r *fishSQLServerRepo) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM fish_configs WHERE id = @p1`
	result, err := r.db.ExecContext(ctx, query, sql.Named("p1", id))
	if err != nil {
		return fmt.Errorf("fishRepo.Delete: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("fishRepo.Delete: fish %q not found", id)
	}
	return nil
}
