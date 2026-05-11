-- ================================================================
-- Fish Game Database - Full Schema Migration
-- Run: sqlcmd -S GIGABYTE -E -C -i migrations\init_schema.sql
-- ================================================================

USE FishGameDB;
GO

-- ─── Table: players ──────────────────────────────────────────────
IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'players')
BEGIN
    CREATE TABLE players (
        id            NVARCHAR(50)   NOT NULL DEFAULT CAST(NEWID() AS NVARCHAR(50)) PRIMARY KEY,
        username      NVARCHAR(100)  NOT NULL UNIQUE,
        email         NVARCHAR(255)  NOT NULL UNIQUE,
        password_hash NVARCHAR(255)  NOT NULL DEFAULT '',
        gold_balance  BIGINT         NOT NULL DEFAULT 0,
        status        NVARCHAR(20)   NOT NULL DEFAULT 'active'
            CHECK (status IN ('active', 'banned', 'suspended')),
        win_rate      FLOAT          NOT NULL DEFAULT 0.0,
        created_at    DATETIME2      NOT NULL DEFAULT GETDATE(),
        last_login_at DATETIME2      NOT NULL DEFAULT GETDATE()
    );
    PRINT 'Created table: players';

    INSERT INTO players (id, username, email, gold_balance, status, win_rate)
    VALUES
        ('USR001', N'HaiKute',     'hai@fishgame.vn',   45000,  'active',  95.0),
        ('USR002', N'SharkHunter', 'shark@fishgame.vn', 1200,   'banned',  120.0),
        ('USR003', N'CaVangNo1',   'ca@fishgame.vn',    900000, 'active',  80.0);
    PRINT 'Seeded table: players (3 rows)';
END
ELSE
    PRINT 'Table players already exists, skipping.';
GO

-- ─── Table: fish_configs ─────────────────────────────────────────
IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'fish_configs')
BEGIN
    CREATE TABLE fish_configs (
        id          NVARCHAR(50)   NOT NULL DEFAULT CAST(NEWID() AS NVARCHAR(50)) PRIMARY KEY,
        name        NVARCHAR(100)  NOT NULL,
        multiplier  INT            NOT NULL CHECK (multiplier >= 1),
        base_prob   FLOAT          NOT NULL CHECK (base_prob >= 0 AND base_prob <= 100),
        speed       NVARCHAR(20)   NOT NULL
            CHECK (speed IN ('fast', 'medium', 'slow', 'very_slow')),
        role        NVARCHAR(20)   NOT NULL
            CHECK (role IN ('common', 'mid', 'boss')),
        is_active   BIT            NOT NULL DEFAULT 1,
        created_at  DATETIME2      NOT NULL DEFAULT GETDATE()
    );
    PRINT 'Created table: fish_configs';

    INSERT INTO fish_configs (id, name, multiplier, base_prob, speed, role, is_active)
    VALUES
        ('F01', N'Cá Nhỏ (Xanh)',  2,    50.0, 'fast',      'common', 1),
        ('F02', N'Cá Đuối',        15,   6.6,  'medium',    'common', 1),
        ('F03', N'Cá Mập',         100,  1.0,  'slow',      'mid',    1),
        ('B01', N'Tiên Cá (Boss)', 1000, 0.1,  'very_slow', 'boss',   1);
    PRINT 'Seeded table: fish_configs (4 rows)';
END
ELSE
    PRINT 'Table fish_configs already exists, skipping.';
GO

-- ─── Table: rooms ────────────────────────────────────────────────
IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'rooms')
BEGIN
    CREATE TABLE rooms (
        id          NVARCHAR(50)   NOT NULL DEFAULT CAST(NEWID() AS NVARCHAR(50)) PRIMARY KEY,
        name        NVARCHAR(100)  NOT NULL,
        type        NVARCHAR(20)   NOT NULL
            CHECK (type IN ('beginner', 'advanced', 'expert', 'vip', 'boss')),
        bet_amount  BIGINT         NOT NULL CHECK (bet_amount >= 1),
        players     INT            NOT NULL DEFAULT 0,
        max_players INT            NOT NULL CHECK (max_players BETWEEN 1 AND 8),
        status      NVARCHAR(20)   NOT NULL DEFAULT 'waiting'
            CHECK (status IN ('waiting', 'playing', 'closed')),
        base_rtp    FLOAT          NOT NULL CHECK (base_rtp BETWEEN 50 AND 100),
        created_at  DATETIME2      NOT NULL DEFAULT GETDATE()
    );
    PRINT 'Created table: rooms';

    INSERT INTO rooms (id, name, type, bet_amount, players, max_players, status, base_rtp)
    VALUES
        ('R-001', N'Biển Tân Thủ 1',    'beginner', 10,    4, 4, 'playing', 90.0),
        ('R-002', N'Biển Đại Dương 1',  'advanced', 100,   2, 4, 'waiting', 85.0),
        ('R-003', N'Vịnh Thử Thách',    'expert',   1000,  4, 4, 'playing', 80.0),
        ('R-004', N'Đảo Kho Báu (VIP)', 'vip',      10000, 1, 4, 'waiting', 88.0);
    PRINT 'Seeded table: rooms (4 rows)';
END
ELSE
    PRINT 'Table rooms already exists, skipping.';
GO

-- ─── Table: admins ───────────────────────────────────────────────
IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'admins')
BEGIN
    CREATE TABLE admins (
        id            NVARCHAR(50)   NOT NULL DEFAULT CAST(NEWID() AS NVARCHAR(50)) PRIMARY KEY,
        username      NVARCHAR(100)  NOT NULL UNIQUE,
        password_hash NVARCHAR(255)  NOT NULL,
        role          NVARCHAR(20)   NOT NULL DEFAULT 'admin'
            CHECK (role IN ('admin', 'superadmin')),
        created_at    DATETIME2      NOT NULL DEFAULT GETDATE()
    );
    PRINT 'Created table: admins';

    -- Default admin account: username=admin / password=admin123
    -- Hash generated by: go build -o hash_pw.exe ./tools/hash_password; .\hash_pw.exe -password admin123
    -- Algorithm: bcrypt, cost=10 — same mechanism used by POST /api/v1/auth/login
    INSERT INTO admins (username, password_hash, role)
    VALUES (
        'admin',
        '$2a$10$04zVpQ2qZoU2JFjER4DRpuRSjYQgWb9lCG9x//r.GNiG5OwClqWC2',
        'superadmin'
    );
    PRINT 'Seeded admin: admin / admin123 (bcrypt cost=10)';
END
ELSE
    PRINT 'Table admins already exists, skipping.';
GO

PRINT '==============================';
PRINT 'Migration completed successfully.';
PRINT '==============================';
