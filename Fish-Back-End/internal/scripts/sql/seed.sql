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
    min_bet       BIGINT       NOT NULL DEFAULT 0,
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

CREATE TABLE IF NOT EXISTS transactions (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount      BIGINT      NOT NULL,
    type        VARCHAR(20) NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Tự động tạo wallet khi user đăng ký
CREATE OR REPLACE FUNCTION fn_create_wallet()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO wallets (user_id) VALUES (NEW.id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_create_wallet
    AFTER INSERT ON users
    FOR EACH ROW EXECUTE FUNCTION fn_create_wallet();

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  VARCHAR(64) NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO rooms (name, min_bet, max_players, rtp) VALUES ('Sảnh Tân Thủ', 10, 4, 0.95), ('Đại Dương Đại Gia', 1000, 4, 0.90);
INSERT INTO fishes (name, health, reward_multiplier, asset_path) VALUES ('Cá Con', 1, 2, '/assets/fish/small_fish.glb'), ('Cá Mập Boss', 500, 100, '/assets/fish/shark_boss.glb');
INSERT INTO roles (role_name) VALUES ('player'), ('admin') ON CONFLICT DO NOTHING;