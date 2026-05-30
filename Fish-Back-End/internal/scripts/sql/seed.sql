CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE TABLE IF NOT EXISTS roles (
    id        INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    role_name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS users (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username      VARCHAR(50)  NOT NULL UNIQUE,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password      VARCHAR(255) NOT NULL,
    role_id       INT       NOT NULL DEFAULT 1 REFERENCES roles(id),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS rooms (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name          VARCHAR(100) NOT NULL,
    max_players   INT          NOT NULL DEFAULT 4,
    description   TEXT,
    rtp           FLOAT        NOT NULL DEFAULT 0.95,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS fishes (
    id                INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name              VARCHAR(100) NOT NULL,
    health            INT          NOT NULL,
    reward_multiplier INT          NOT NULL,
    speed             FLOAT        NOT NULL DEFAULT 1.0,
    asset_path        TEXT         NOT NULL,
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS wallets (
    user_id    BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    balance    BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Session chơi game: mỗi lần vào phòng = 1 session
CREATE TABLE IF NOT EXISTS game_sessions (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    room_id     BIGINT      NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    shots_fired INT         NOT NULL DEFAULT 0,
    fish_killed INT         NOT NULL DEFAULT 0,
    total_spend BIGINT      NOT NULL DEFAULT 0,  -- tổng vàng đã chi (đạn)
    total_earn  BIGINT      NOT NULL DEFAULT 0,  -- tổng vàng đã nhận (cá)
    status      VARCHAR(20) NOT NULL DEFAULT 'active', -- active | finished
    started_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at    TIMESTAMPTZ
);

-- Lịch sử giao dịch:
--   type = 'play'     → chơi game, có session_id, amount = earn - spend (có thể âm)
--   type = 'deposit'  → nạp vàng, session_id NULL, amount > 0
--   type = 'withdraw' → rút vàng, session_id NULL, amount < 0
CREATE TABLE IF NOT EXISTS transactions (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id  BIGINT      REFERENCES game_sessions(id),
    amount      BIGINT      NOT NULL,
    type        VARCHAR(20) NOT NULL CHECK (type IN ('play', 'deposit', 'withdraw')),
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS refresh_tokens (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  VARCHAR(64) NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO rooms (name, max_players, rtp) VALUES ('Sảnh Tân Thủ', 4, 0.95), ('Đại Dương', 4, 0.90);
INSERT INTO fishes (name, health, reward_multiplier, asset_path) VALUES ('Cá Con', 1, 2, '/assets/fish/small_fish.glb'), ('Cá Mập Boss', 500, 100, '/assets/fish/shark_boss.glb');
INSERT INTO roles (role_name) VALUES ('player'), ('admin') ON CONFLICT DO NOTHING;
