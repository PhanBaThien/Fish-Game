-- ============================================================
-- Fish Game Management System - Full Database Initialization
-- PostgreSQL 15+ (Generated from Go Models)
-- ============================================================

-- Kích hoạt các Extension cần thiết
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- Hỗ trợ tìm kiếm nhanh tên người chơi

-- 1. BẢNG ADMIN (Từ struct Admin)
CREATE TABLE IF NOT EXISTS admin_users (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username      VARCHAR(50)  NOT NULL UNIQUE,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,
    role          VARCHAR(20)  NOT NULL DEFAULT 'admin' CHECK (role IN ('super_admin', 'admin', 'moderator')),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

---tại sạo lại set role là varchar ?  password_hash là text ? 

-- 2. BẢNG NGƯỜI CHƠI (Từ struct Player)
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

CREATE INDEX IF NOT EXISTS idx_players_status ON players(status);
CREATE INDEX IF NOT EXISTS idx_players_username ON players USING gin(username gin_trgm_ops);

-- 3. BẢNG VÍ (Từ struct Wallet)
CREATE TABLE IF NOT EXISTS wallets (
    player_id  UUID PRIMARY KEY REFERENCES players(id) ON DELETE CASCADE,
    balance    BIGINT       NOT NULL DEFAULT 0,
    version    INT          NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- 4. BẢNG CẤU HÌNH CÁ (Từ struct Fish)
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

-- 5. BẢNG PHÒNG CHƠI (Từ struct Room)
CREATE TABLE IF NOT EXISTS rooms (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(100) NOT NULL,
    type        VARCHAR(20)  NOT NULL CHECK (type IN ('beginner', 'advanced', 'expert', 'vip', 'boss')),
    bet_amount  BIGINT       NOT NULL CHECK (bet_amount > 0),
    max_players INT          NOT NULL DEFAULT 4 CHECK (max_players BETWEEN 1 AND 8),
    status      VARCHAR(20)  NOT NULL DEFAULT 'waiting' CHECK (status IN ('waiting', 'playing', 'closed')),
    base_rtp    NUMERIC(5,2) NOT NULL DEFAULT 90.00 CHECK (base_rtp BETWEEN 50 AND 100),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- 6. BẢNG CHỖ NGỒI TRONG PHÒNG (Từ struct RoomSeat)
CREATE TABLE IF NOT EXISTS room_seats (
    room_id     UUID    REFERENCES rooms(id) ON DELETE CASCADE,
    seat_index  INT     NOT NULL CHECK (seat_index BETWEEN 0 AND 7),
    occupied_by UUID    REFERENCES players(id) ON DELETE SET NULL,
    joined_at   TIMESTAMPTZ,
    PRIMARY KEY (room_id, seat_index)
);

-- 7. BẢNG PHIÊN CHƠI (Từ struct GameSession)
CREATE TABLE IF NOT EXISTS game_sessions (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id     UUID        NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    player_id   UUID        NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    gold_start  BIGINT      NOT NULL,
    gold_end    BIGINT,
    fish_caught JSONB       NOT NULL DEFAULT '[]',
    started_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at    TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_sessions_player ON game_sessions(player_id);
CREATE INDEX IF NOT EXISTS idx_sessions_room   ON game_sessions(room_id);

-- 8. BẢNG ĐẠN BẮN (Từ struct Shot)
CREATE TABLE IF NOT EXISTS shots (
    id              BIGSERIAL PRIMARY KEY,
    session_id      UUID         NOT NULL REFERENCES game_sessions(id) ON DELETE CASCADE,
    player_id       UUID         NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    room_id         UUID         NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    bet_amount      BIGINT       NOT NULL,
    angle           NUMERIC(5,2) NOT NULL,
    shot_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    idempotency_key VARCHAR(100) NOT NULL UNIQUE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- 9. BẢNG LỊCH SỬ DIỆT CÁ (Từ struct FishKill)
CREATE TABLE IF NOT EXISTS fish_kills (
    id         BIGSERIAL PRIMARY KEY,
    shot_id    BIGINT       NOT NULL REFERENCES shots(id) ON DELETE CASCADE,
    fish_id    VARCHAR(20)  NOT NULL REFERENCES fish(id),
    payout     BIGINT       NOT NULL,
    rng_seed   BIGINT       NOT NULL,
    prob_used  NUMERIC(6,3) NOT NULL,
    killed_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- 10. BẢNG LỊCH SỬ GIAO DỊCH (Từ struct Transaction)
CREATE TABLE IF NOT EXISTS transactions (
    id              BIGSERIAL PRIMARY KEY,
    player_id       UUID         NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    type            VARCHAR(20)  NOT NULL CHECK (type IN ('shot', 'win', 'gift', 'deposit', 'withdraw', 'adjust')),
    amount          BIGINT       NOT NULL,
    balance_after   BIGINT       NOT NULL,
    ref_shot_id     BIGINT       REFERENCES shots(id) ON DELETE SET NULL,
    ref_kill_id     BIGINT       REFERENCES fish_kills(id) ON DELETE SET NULL,
    idempotency_key VARCHAR(100) NOT NULL UNIQUE,
    metadata        JSONB        NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- 11. BẢNG NHẬT KÝ HỆ THỐNG (Từ struct AuditLog)
CREATE TABLE IF NOT EXISTS audit_logs (
    id         BIGSERIAL PRIMARY KEY,
    actor_type VARCHAR(20) NOT NULL CHECK (actor_type IN ('admin', 'player', 'system')),
    actor_id   UUID        NOT NULL, -- Lưu ID của admin hoặc player tương ứng
    action     VARCHAR(100) NOT NULL,
    target     VARCHAR(100) NOT NULL,
    payload    JSONB        NOT NULL DEFAULT '{}',
    ip         VARCHAR(45)  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- 12. BẢNG CÀI ĐẶT CHUNG (Từ struct Setting)
CREATE TABLE IF NOT EXISTS settings (
    key        VARCHAR(100) PRIMARY KEY,
    value      JSONB        NOT NULL DEFAULT '{}',
    updated_by UUID         NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ============================================================
-- DỮ LIỆU KHỞI TẠO ĐỂ TEST (SEED DATA)
-- ============================================================

-- Thêm các loại cá mặc định
INSERT INTO fish (id, name, multiplier, base_prob, speed, role) VALUES
    ('F01', 'Cá Nhỏ (Xanh)',    2,    50.0, 'fast',      'common'),
    ('F02', 'Cá Đuối',          15,   6.6,  'medium',    'common'),
    ('F03', 'Cá Mập',           100,  1.0,  'slow',      'mid'),
    ('F04', 'Cá Heo',           30,   3.5,  'medium',    'common'),
    ('F05', 'Cá Voi',           200,  0.5,  'slow',      'mid'),
    ('B01', 'Tiên Cá (Boss)',   1000, 0.1,  'very_slow', 'boss'),
    ('B02', 'Rồng Biển',        5000, 0.01, 'very_slow', 'boss')
ON CONFLICT (id) DO NOTHING;

-- Thêm phòng chơi mặc định
INSERT INTO rooms (name, type, bet_amount, max_players, base_rtp) VALUES
    ('Biển Tân Thủ 1',       'beginner', 10,    4, 90.00),
    ('Biển Đại Dương 1',     'advanced', 100,   4, 85.00),
    ('Vịnh Thử Thách',       'expert',   1000,  4, 80.00),
    ('Đảo Kho Báu (VIP)',    'vip',      10000, 4, 88.00),
    ('Bão Táp Biển Sâu',     'boss',     5000,  4, 75.00)
ON CONFLICT DO NOTHING;

-- Tạo tài khoản Admin tối cao mặc định (Mật khẩu: admin123)
INSERT INTO admin_users (username, email, password_hash, role) VALUES
    ('admin', 'admin@fishgame.local', '$2a$12$6R8gTbe8a9wVw1/fQpY1beC1G3ClyQY1E77Hl63K8BvVjRWhXgKDG', 'super_admin')
ON CONFLICT (username) DO NOTHING;