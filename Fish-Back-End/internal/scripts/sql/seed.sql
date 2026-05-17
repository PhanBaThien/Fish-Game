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
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS fishs (
    id                INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name              VARCHAR(100) NOT NULL,
    health            INT          NOT NULL,
    reward_multiplier INT          NOT NULL,
    speed             FLOAT        NOT NULL DEFAULT 1.0,
    asset_path        TEXT         NOT NULL,
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

INSERT INTO rooms (name, min_bet, max_players) VALUES ('Sảnh Tân Thủ', 10, 4), ('Đại Dương Đại Gia', 1000, 4);
INSERT INTO fishs (name, health, reward_multiplier, asset_path) VALUES ('Cá Con', 1, 2, '/assets/fish/small_fish.glb'), ('Cá Mập Boss', 500, 100, '/assets/fish/shark_boss.glb');
INSERT INTO roles (role_name) VALUES ('player'), ('admin') ON CONFLICT DO NOTHING;