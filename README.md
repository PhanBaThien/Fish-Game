# 🐟 Fish Game

Multiplayer 2D fish shooting game với kiến trúc full-stack: **Go/Gin** REST API + **React 18** frontend render bằng **Canvas 2D API**.

---

## Tech Stack

| Layer | Công nghệ |
|---|---|
| **Backend language** | Go 1.25 |
| **HTTP framework** | Gin v1.9 |
| **Database** | PostgreSQL 17 |
| **DB driver** | pgx/v5 + pgxpool |
| **Code generation** | sqlc v2 |
| **Dependency injection** | Google Wire |
| **Authentication** | JWT (golang-jwt/jwt v5) |
| **Frontend** | React 18 + Vite 5 |
| **State management** | Zustand v4 |
| **Data fetching** | TanStack Query v5 |
| **HTTP client** | Axios v1 |
| **Styling** | Tailwind CSS v3 |
| **Game engine** | HTML5 Canvas 2D API |
| **Container** | Docker + Docker Compose |
| **Reverse proxy** | nginx (production) |

---

## Cấu trúc project

```
Fish-Game/
├── Fish-Back-End/
│   ├── cmd/api/
│   │   ├── main.go              # Entry point
│   │   └── wire_gen.go          # Wire-generated DI
│   ├── internal/
│   │   ├── database/
│   │   │   └── postgre.go       # Pool config + auto-create DB + seed
│   │   ├── domain/              # Request/Response DTOs, role constants
│   │   ├── models/              # DB model structs
│   │   ├── repository/
│   │   │   ├── dbgen/           # sqlc-generated Go code (KHÔNG sửa tay)
│   │   │   ├── auth_repo.go
│   │   │   ├── room_repo.go
│   │   │   ├── fish_repo.go
│   │   │   └── refresh_token_repo.go
│   │   ├── transport/http/
│   │   │   ├── middleware/      # Auth, CORS, Logger, RequireRoles
│   │   │   ├── auth_handler.go
│   │   │   ├── room_handler.go
│   │   │   ├── fish_handler.go
│   │   │   ├── response.go      # Success() / Fail() helpers
│   │   │   └── router.go
│   │   ├── usecase/             # Business logic layer
│   │   └── scripts/sql/
│   │       └── seed.sql         # Schema + dữ liệu mẫu (chạy lúc init)
│   ├── db/query/                # sqlc query files (.sql)
│   ├── pkg/
│   │   ├── apperror/            # AppError + InternalError + error vars
│   │   └── utils/               # TokenMaker, PasswordHasher, ToInt64
│   └── sqlc.yaml
│
├── Fish-Front-End/
│   ├── src/
│   │   ├── api/
│   │   │   ├── client.ts        # axios instance + auto-refresh interceptor
│   │   │   ├── auth.ts
│   │   │   ├── rooms.ts
│   │   │   └── fish.ts
│   │   ├── components/
│   │   │   ├── Navbar.tsx
│   │   │   ├── ProtectedRoute.tsx
│   │   │   └── RoomCard.tsx
│   │   ├── game/
│   │   │   ├── entities/
│   │   │   │   ├── FishEntity.ts   # Cá 2D (Canvas)
│   │   │   │   └── BulletEntity.ts # Đạn 2D (Canvas)
│   │   │   ├── scenes/
│   │   │   │   └── GameScene.ts    # Game loop + render + collision
│   │   │   └── GameCanvas.tsx      # React wrapper cho canvas
│   │   ├── pages/
│   │   │   ├── LoginPage.tsx
│   │   │   ├── LobbyPage.tsx
│   │   │   └── GamePage.tsx
│   │   ├── stores/
│   │   │   ├── authStore.ts     # User + access token (Zustand + persist)
│   │   │   └── gameStore.ts     # Coins + score
│   │   └── types/index.ts       # TypeScript interfaces
│   ├── nginx.conf               # Reverse proxy + SPA fallback
│   └── Dockerfile
│
├── docker-compose.yml
├── .env.example                 # Template biến môi trường
└── .gitignore
```

---

## Cài đặt và chạy

### Phương án 1 — Docker (khuyến nghị)

> Không cần cài Go, Node, hay PostgreSQL. Chỉ cần Docker Desktop.

**Yêu cầu:** [Docker Desktop](https://www.docker.com/products/docker-desktop/) đang chạy.

#### Bước 1 — Clone repo

```bash
git clone <repo-url>
cd Fish-Game
```

#### Bước 2 — Tạo file `.env`

```bash
cp .env.example .env
```

Mở `.env` và thay thế bằng secret thực:

```env
ACCESS_TOKEN_KEY=thay_bang_chuoi_bat_ky_dai_hon_32_ky_tu_!!
REFRESH_TOKEN_KEY=thay_bang_chuoi_khac_access_key_dai_hon_32_!!
```

> ⚠️ **Bắt buộc:** Hai key phải **khác nhau**, độ dài **≥ 32 ký tự**. Không để giá trị mặc định khi deploy.

#### Bước 3 — Build và khởi động

```bash
docker compose up --build
```

Lần đầu chạy mất khoảng 3–5 phút (tải image, compile Go, build React).  
Các lần sau không cần `--build` trừ khi thay đổi code:

```bash
docker compose up
```

#### Bước 4 — Truy cập

| Service | URL |
|---|---|
| 🎮 Frontend (game) | http://localhost:3000 |
| 🔌 Backend API | http://localhost:8080/api/v1 |
| 🗄️ PostgreSQL (từ máy host) | `localhost:5433` |

> Database được tạo tự động khi backend khởi động lần đầu. Không cần chạy migration thủ công.

#### Dừng services

```bash
docker compose down          # dừng, giữ nguyên dữ liệu DB
docker compose down -v       # dừng + xóa sạch DB (reset hoàn toàn)
```

---

### Phương án 2 — Chạy thủ công (Development)

#### Yêu cầu

| Tool | Version tối thiểu | Link |
|---|---|---|
| Go | 1.22 | https://go.dev/dl |
| Node.js | 20 LTS | https://nodejs.org |
| PostgreSQL | 15 | https://www.postgresql.org/download |

---

#### Thiết lập Backend

**1. Tạo database**

Mở `psql` hoặc pgAdmin, chạy:

```sql
CREATE DATABASE fish_game;
```

Sau đó kết nối vào `fish_game` và chạy schema:

```bash
psql -U postgres -d fish_game -f Fish-Back-End/internal/scripts/sql/seed.sql
```

**2. Tạo file `.env`**

```bash
cd Fish-Back-End
cp .env.example .env
```

Chỉnh nội dung `.env`:

```env
SERVER_PORT=8080
DATABASE_URL=postgres://postgres:your_password@localhost:5432/fish_game?sslmode=disable

ACCESS_TOKEN_KEY=your_access_token_secret_min_32_chars!!
ACCESS_TOKEN_EXPIRY=15m

REFRESH_TOKEN_KEY=your_refresh_token_secret_different!!
REFRESH_TOKEN_EXPIRY=168h

ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

**3. Tải dependencies và chạy**

```bash
go mod download
go run ./cmd/api
```

Server sẽ khởi động tại `http://localhost:8080`.

---

#### Thiết lập Frontend

**1. Cài dependencies**

```bash
cd Fish-Front-End
npm install
```

**2. Chạy dev server**

```bash
npm run dev
```

Mở trình duyệt tại `http://localhost:5173`.

> Dev server đã cấu hình proxy: mọi request đến `/api/*` sẽ tự động chuyển tiếp đến `http://localhost:8080`.

---

## Biến môi trường

### Root `.env` (dùng cho Docker Compose)

| Biến | Bắt buộc | Mô tả |
|---|---|---|
| `ACCESS_TOKEN_KEY` | ✅ | Secret ký JWT access token, ≥ 32 ký tự |
| `REFRESH_TOKEN_KEY` | ✅ | Secret ký JWT refresh token, **khác** access key, ≥ 32 ký tự |
| `ACCESS_TOKEN_EXPIRY` | ❌ | Thời hạn access token (default: `15m`) |
| `REFRESH_TOKEN_EXPIRY` | ❌ | Thời hạn refresh token (default: `168h`) |

### `Fish-Back-End/.env` (dùng khi chạy thủ công)

| Biến | Mô tả |
|---|---|
| `SERVER_PORT` | Port server lắng nghe (default: `8080`) |
| `DATABASE_URL` | PostgreSQL connection string |
| `ACCESS_TOKEN_KEY` | Như trên |
| `ACCESS_TOKEN_EXPIRY` | Như trên |
| `REFRESH_TOKEN_KEY` | Như trên |
| `REFRESH_TOKEN_EXPIRY` | Như trên |
| `ALLOWED_ORIGINS` | Danh sách origin cho CORS, phân cách bằng dấu phẩy |

---

## API Reference

Tất cả response theo format chuẩn:

```json
// Thành công
{ "data": { ... }, "error": null }

// Lỗi
{ "data": null, "error": { "code": "ERROR_CODE", "message": "mô tả lỗi" } }
```

### Auth

| Method | Endpoint | Auth | Mô tả |
|---|---|---|---|
| `POST` | `/api/v1/auth/register` | — | Đăng ký tài khoản mới |
| `POST` | `/api/v1/auth/login` | — | Đăng nhập, trả về access token + set HttpOnly cookie |
| `POST` | `/api/v1/auth/refresh` | cookie | Làm mới access token (refresh token đọc từ cookie) |
| `GET` | `/api/v1/auth/me` | Bearer | Thông tin user hiện tại |
| `POST` | `/api/v1/auth/logout` | Bearer | Đăng xuất, xóa refresh token |

**Ví dụ — Login request:**
```json
POST /api/v1/auth/login
{
  "username": "player1",
  "password": "password123"
}
```

**Login response:**
```json
{
  "data": {
    "access_token": "eyJhbGci...",
    "access_token_expires_at": 1700000000,
    "user": {
      "id": 1, "username": "player1",
      "email": "player1@example.com", "role_id": 1
    }
  },
  "error": null
}
```
> Refresh token được set tự động trong HttpOnly cookie `refresh_token`, không xuất hiện trong JSON response.

### Rooms

| Method | Endpoint | Auth | Mô tả |
|---|---|---|---|
| `GET` | `/api/v1/rooms` | Bearer | Danh sách tất cả phòng |
| `GET` | `/api/v1/rooms/:id` | Bearer | Chi tiết một phòng |
| `POST` | `/api/v1/rooms` | Bearer + Admin | Tạo phòng mới |
| `PUT` | `/api/v1/rooms/:id` | Bearer + Admin | Cập nhật phòng |
| `DELETE` | `/api/v1/rooms/:id` | Bearer + Admin | Xóa phòng |

### Fishes

| Method | Endpoint | Auth | Mô tả |
|---|---|---|---|
| `GET` | `/api/v1/fishes` | Bearer | Danh sách loại cá |
| `GET` | `/api/v1/fishes/:id` | Bearer | Chi tiết một loại cá |
| `POST` | `/api/v1/fishes` | Bearer + Admin | Tạo loại cá mới |
| `PUT` | `/api/v1/fishes/:id` | Bearer + Admin | Cập nhật loại cá |
| `DELETE` | `/api/v1/fishes/:id` | Bearer + Admin | Xóa loại cá |

### Role ID

| role_id | Vai trò | Quyền hạn |
|---|---|---|
| `1` | Player | Xem rooms, xem fishes, chơi game |
| `2` | Admin | Player + tạo/sửa/xóa rooms và fishes |
| `3` | Super Admin | Toàn quyền |

---

## Kết nối pgAdmin với Docker PostgreSQL

Khi chạy Docker, PostgreSQL chạy bên trong container, expose ra cổng **5433** (để tránh conflict với PostgreSQL local cổng 5432).

1. Mở pgAdmin → chuột phải **Servers** → **Register → Server...**
2. Tab **General**: đặt `Name` tuỳ ý (vd: `Fish Docker`)
3. Tab **Connection**:
   ```
   Host name/address : 127.0.0.1
   Port              : 5433
   Maintenance DB    : postgres
   Username          : postgres
   Password          : postgres
   ```
4. Bấm **Save**

---

## Công cụ phát triển

### sqlc — Tái generate code từ SQL

Khi thêm/sửa file trong `db/query/*.sql` hoặc `internal/scripts/sql/seed.sql`, chạy lại:

```bash
cd Fish-Back-End
sqlc generate
```

> Cài sqlc: https://docs.sqlc.dev/en/latest/overview/install.html

File được generate ra nằm trong `internal/repository/dbgen/` — **không sửa tay**.

### Wire — Tái generate Dependency Injection

Khi thêm dependency mới vào `cmd/api/wire.go`:

```bash
cd Fish-Back-End/cmd/api
wire
```

> Cài Wire: `go install github.com/google/wire/cmd/wire@latest`

---

## Bảo mật

| Cơ chế | Chi tiết |
|---|---|
| **Refresh token lưu DB** | Lưu dưới dạng SHA-256 hash, không bao giờ lưu plain text |
| **HttpOnly Cookie** | Refresh token truyền qua cookie, không thể đọc bằng JavaScript (chống XSS) |
| **Token Rotation** | Mỗi lần refresh sẽ tạo refresh token mới, invalidate token cũ |
| **Tách secret** | Access token và refresh token dùng 2 secret key hoàn toàn khác nhau |
| **SameSite=Lax** | Cookie chỉ gửi khi navigation từ cùng site (giảm thiểu CSRF) |
| **Role-based access** | Endpoint admin kiểm tra `role_id` từ JWT claim |

---

## Xử lý sự cố thường gặp

**`docker compose up` lỗi "ACCESS_TOKEN_KEY is required"**
→ Chưa tạo file `.env`. Chạy `cp .env.example .env` rồi điền giá trị.

**Database trống sau khi chạy Docker**
→ Volume cũ còn tồn tại. Chạy `docker compose down -v` để xóa sạch rồi `docker compose up --build` lại.

**Frontend build lỗi exit code 2**
→ Thường là lỗi TypeScript. Chạy `cd Fish-Front-End && npm run build` để xem chi tiết lỗi.

**Không kết nối được pgAdmin vào Docker**
→ Dùng port `5433` (không phải `5432`) và password `postgres`.

**Backend báo "lỗi kết nối pool" lúc khởi động**
→ PostgreSQL chưa sẵn sàng. `docker compose` đã cấu hình `depends_on` + healthcheck, thử `docker compose restart backend`.
