-- ============================================================
-- Fish Game Management System - Database Initialization
-- PostgreSQL 15+
-- ============================================================

-- Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- For fuzzy text search
-- change the name table Users:
-- ─── Players Table ────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS players (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username      VARCHAR(50)  NOT NULL UNIQUE,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,
    gold_balance  BIGINT       NOT NULL DEFAULT 0 CHECK (gold_balance >= 0),
    status        VARCHAR(20)  NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'banned', 'suspended')),
    win_rate      NUMERIC(6,2) NOT NULL DEFAULT 100.00,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMPTZ
);

CREATE INDEX idx_players_status     ON players(status);
CREATE INDEX idx_players_username   ON players USING gin(username gin_trgm_ops);

-- ─── Fish Config Table ────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS fish (
    id          VARCHAR(20)  PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    multiplier  INT          NOT NULL CHECK (multiplier > 0),
    base_prob   NUMERIC(6,3) NOT NULL CHECK (base_prob >= 0 AND base_prob <= 100),
    speed       VARCHAR(20)  NOT NULL CHECK (speed IN ('fast', 'medium', 'slow', 'very_slow')),
    role        VARCHAR(20)  NOT NULL CHECK (role IN ('common', 'mid', 'boss')),
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ─── Rooms Table ─────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS rooms (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(100) NOT NULL,
    type        VARCHAR(20)  NOT NULL CHECK (type IN ('beginner', 'advanced', 'expert', 'vip', 'boss')),
    bet_amount  BIGINT       NOT NULL CHECK (bet_amount > 0),
    max_players INT          NOT NULL DEFAULT 4 CHECK (max_players BETWEEN 1 AND 8),
    status      VARCHAR(20)  NOT NULL DEFAULT 'waiting' CHECK (status IN ('waiting', 'playing', 'closed')),
    base_rtp    NUMERIC(5,2) NOT NULL DEFAULT 90.00 CHECK (base_rtp BETWEEN 50 AND 100),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);


-- ─── Game Sessions Table ──────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS game_sessions (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id       UUID         NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    player_id     UUID         NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    gold_start    BIGINT       NOT NULL,
    gold_end      BIGINT,
    fish_caught   JSONB        NOT NULL DEFAULT '[]',
    started_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    ended_at      TIMESTAMPTZ
);

CREATE INDEX idx_sessions_player ON game_sessions(player_id);
CREATE INDEX idx_sessions_room   ON game_sessions(room_id);

--change import row UserType( if value = 1 data type char(1))
-- ─── Admin Users Table ────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS admin_users (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username      VARCHAR(50) NOT NULL UNIQUE,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT        NOT NULL,
    role          VARCHAR(20) NOT NULL DEFAULT 'admin' CHECK (role IN ('super_admin', 'admin', 'moderator')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ─── Seed Data ────────────────────────────────────────────────────────────────

-- Default fish types
INSERT INTO fish (id, name, multiplier, base_prob, speed, role) VALUES
    ('F01', 'Cá Nhỏ (Xanh)',    2,    50.0, 'fast',      'common'),
    ('F02', 'Cá Đuối',          15,   6.6,  'medium',    'common'),
    ('F03', 'Cá Mập',           100,  1.0,  'slow',      'mid'),
    ('F04', 'Cá Heo',           30,   3.5,  'medium',    'common'),
    ('F05', 'Cá Voi',           200,  0.5,  'slow',      'mid'),
    ('B01', 'Tiên Cá (Boss)',   1000, 0.1,  'very_slow', 'boss'),
    ('B02', 'Rồng Biển',        5000, 0.01, 'very_slow', 'boss')
ON CONFLICT (id) DO NOTHING;

-- Default game rooms
INSERT INTO rooms (name, type, bet_amount, max_players, base_rtp) VALUES
    ('Biển Tân Thủ 1',       'beginner', 10,    4, 90.00),
    ('Biển Đại Dương 1',     'advanced', 100,   4, 85.00),
    ('Vịnh Thử Thách',       'expert',   1000,  4, 80.00),
    ('Đảo Kho Báu (VIP)',    'vip',      10000, 4, 88.00),
    ('Bão Táp Biển Sâu',     'boss',     5000,  4, 75.00)
ON CONFLICT DO NOTHING;

-- Default admin user (password: admin123 - CHANGE IN PRODUCTION)
INSERT INTO admin_users (username, email, password_hash, role) VALUES
    ('admin', 'admin@fishgame.local', '$2a$12$placeholder_bcrypt_hash_change_me', 'super_admin')
ON CONFLICT (username) DO NOTHING;
