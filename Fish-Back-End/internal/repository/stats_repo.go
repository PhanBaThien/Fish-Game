package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/domain"
)

// StatsRepository định nghĩa interface lấy thống kê tổng hợp
type StatsRepository interface {
	Overview(ctx context.Context) (*domain.StatsOverviewResponse, error)
	Timeseries(ctx context.Context, from, to time.Time) ([]domain.TimeseriesPoint, error)
}

type statsPgRepo struct {
	db *sql.DB
}

func NewStatsRepository(db *sql.DB) StatsRepository {
	return &statsPgRepo{db: db}
}

func (r *statsPgRepo) Overview(ctx context.Context) (*domain.StatsOverviewResponse, error) {
	var resp domain.StatsOverviewResponse

	// Thống kê players
	err := r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE status = 'active'),
			COUNT(*) FILTER (WHERE status = 'banned')
		FROM players`).
		Scan(&resp.TotalPlayers, &resp.ActivePlayers, &resp.BannedPlayers)
	if err != nil {
		return nil, fmt.Errorf("statsRepo.Overview players: %w", err)
	}

	// Thống kê rooms
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*), COUNT(*) FILTER (WHERE status = 'playing')
		FROM rooms`).
		Scan(&resp.TotalRooms, &resp.ActiveRooms)
	if err != nil {
		return nil, fmt.Errorf("statsRepo.Overview rooms: %w", err)
	}

	// Thống kê transactions
	err = r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*),
			COALESCE(SUM(amount) FILTER (WHERE type IN ('deposit','gift')), 0),
			COALESCE(SUM(amount) FILTER (WHERE type IN ('withdraw','win')), 0)
		FROM transactions`).
		Scan(&resp.TotalTransactions, &resp.TotalGoldIn, &resp.TotalGoldOut)
	if err != nil {
		return nil, fmt.Errorf("statsRepo.Overview transactions: %w", err)
	}

	return &resp, nil
}

func (r *statsPgRepo) Timeseries(ctx context.Context, from, to time.Time) ([]domain.TimeseriesPoint, error) {
	query := `
		WITH dates AS (
			SELECT generate_series($1::date, $2::date, '1 day'::interval)::date AS d
		),
		tx_agg AS (
			SELECT
				created_at::date AS d,
				COALESCE(SUM(amount) FILTER (WHERE type IN ('deposit','gift')), 0)    AS gold_in,
				COALESCE(SUM(amount) FILTER (WHERE type IN ('withdraw','win')), 0)    AS gold_out
			FROM transactions
			WHERE created_at BETWEEN $1 AND $2
			GROUP BY 1
		),
		player_agg AS (
			SELECT created_at::date AS d, COUNT(*) AS new_players
			FROM players
			WHERE created_at BETWEEN $1 AND $2
			GROUP BY 1
		),
		shot_agg AS (
			SELECT shot_at::date AS d, COUNT(*) AS shots
			FROM shots
			WHERE shot_at BETWEEN $1 AND $2
			GROUP BY 1
		)
		SELECT
			dates.d::text,
			COALESCE(tx_agg.gold_in, 0),
			COALESCE(tx_agg.gold_out, 0),
			COALESCE(player_agg.new_players, 0),
			COALESCE(shot_agg.shots, 0)
		FROM dates
		LEFT JOIN tx_agg     ON tx_agg.d = dates.d
		LEFT JOIN player_agg ON player_agg.d = dates.d
		LEFT JOIN shot_agg   ON shot_agg.d = dates.d
		ORDER BY dates.d`

	rows, err := r.db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, fmt.Errorf("statsRepo.Timeseries: %w", err)
	}
	defer rows.Close()

	var points []domain.TimeseriesPoint
	for rows.Next() {
		var p domain.TimeseriesPoint
		if err := rows.Scan(&p.Date, &p.GoldIn, &p.GoldOut, &p.NewPlayers, &p.Shots); err != nil {
			return nil, fmt.Errorf("statsRepo.Timeseries scan: %w", err)
		}
		points = append(points, p)
	}
	return points, nil
}
