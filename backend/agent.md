# AntriSehat Backend — Agent Rules & Conventions

> **Baca file ini sebelum menulis atau memodifikasi kode apapun.**
> File ini adalah sumber kebenaran tunggal (single source of truth) untuk standar penulisan kode backend AntriSehat.

---

## 1. Stack & Versi

| Komponen | Teknologi | Versi |
|---|---|---|
| Language | Go | 1.22+ |
| HTTP Framework | Gin | v1.12+ |
| Database | PostgreSQL | 16 |
| Driver DB | `go.mongodb.org/...` ❌ — gunakan `pgx/v5` | v5 |
| Migration | `golang-migrate` | v4 |
| Auth | JWT (`golang-jwt/jwt/v5`) | v5 |
| Password Hash | `golang.org/x/crypto/bcrypt` | latest |
| Validation | `go-playground/validator/v10` | v10 |
| Real-time | SSE atau `gorilla/websocket` | latest |
| Config | `goccy/go-yaml` atau env | — |
| Logger | `log/slog` (stdlib, structured) | Go 1.21+ |

---

## 2. Struktur Direktori (WAJIB diikuti)

```
backend/
├── cmd/
│   └── api/
│       └── main.go              # Entrypoint: init config, DB, router, server
├── internal/
│   ├── config/
│   │   └── config.go            # Struct config + loader dari env/yaml
│   ├── domain/
│   │   ├── entity/              # Pure Go structs (no DB tag, no JSON tag bercampur)
│   │   │   ├── user.go
│   │   │   ├── patient.go
│   │   │   ├── queue.go
│   │   │   └── ...
│   │   └── repository/          # Interface repository (contract)
│   │       ├── user_repository.go
│   │       ├── patient_repository.go
│   │       └── ...
│   ├── handler/                 # HTTP handler (thin layer, hanya parse & respond)
│   │   ├── auth_handler.go
│   │   ├── patient_handler.go
│   │   ├── queue_handler.go
│   │   └── ...
│   ├── middleware/              # Gin middleware
│   │   ├── auth.go              # JWT validation
│   │   ├── rbac.go              # Role-based access
│   │   ├── rate_limit.go
│   │   ├── audit.go             # Audit log writer
│   │   └── error_handler.go     # Global error handler
│   ├── repository/              # Implementasi interface domain/repository
│   │   ├── postgres/
│   │   │   ├── user_repo.go
│   │   │   ├── patient_repo.go
│   │   │   └── ...
│   │   └── db.go                # Init koneksi pgx pool
│   ├── service/                 # Business logic (orchestrate repo calls)
│   │   ├── auth_service.go
│   │   ├── patient_service.go
│   │   ├── queue_service.go
│   │   └── ...
│   ├── dto/                     # Request & Response structs (dengan json/validate tag)
│   │   ├── auth_dto.go
│   │   ├── patient_dto.go
│   │   └── ...
│   ├── apperror/                # Custom error types & sentinel errors
│   │   └── errors.go
│   ├── worker/                  # Background jobs (notifikasi, cleanup)
│   │   └── notification_worker.go
│   └── router/
│       └── router.go            # Route registration
├── migrations/                  # File SQL up/down
│   ├── 000001_init_schema.up.sql
│   └── 000001_init_schema.down.sql
├── pkg/                         # Shared utility yang bisa dipakai ulang
│   ├── hash/
│   │   └── bcrypt.go            # Hash & verify password
│   ├── jwt/
│   │   └── jwt.go               # Generate & parse JWT
│   ├── qrcode/
│   │   └── qr.go
│   └── response/
│       └── response.go          # Standard HTTP response wrapper
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── agent.md                     # ← file ini
└── go.mod
```

---

## 3. Konvensi Penamaan

### File
- Semua huruf kecil dengan underscore: `patient_handler.go`, `auth_service.go`
- Satu domain per file, jangan gabungkan dua entity dalam satu file

### Struct & Interface
```go
// Entity → nama benda, tanpa suffix
type Patient struct { ... }
type Queue struct { ... }

// Interface repository → suffix "Repository"
type PatientRepository interface { ... }

// Interface service → suffix "Service"
type PatientService interface { ... }

// DTO → suffix "Request" atau "Response"
type CreatePatientRequest struct { ... }
type PatientResponse struct { ... }

// Handler → suffix "Handler"
type PatientHandler struct { ... }
```

### Fungsi
```go
// Repository → CRUD eksplisit
func (r *patientRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Patient, error)
func (r *patientRepo) FindByNIKHash(ctx context.Context, nikHash string) (*domain.Patient, error)
func (r *patientRepo) Create(ctx context.Context, p *domain.Patient) error
func (r *patientRepo) Update(ctx context.Context, p *domain.Patient) error
func (r *patientRepo) SoftDelete(ctx context.Context, id uuid.UUID) error

// Service → action-oriented
func (s *patientService) RegisterPatient(ctx context.Context, req dto.CreatePatientRequest) (*dto.PatientResponse, error)
func (s *queueService) AdvanceQueueStatus(ctx context.Context, queueID uuid.UUID, newStatus domain.QueueStatus) error
```

---

## 4. Error Handling (WAJIB)

### 4.1 Buat Custom Error di `internal/apperror/errors.go`

```go
package apperror

import "net/http"

type AppError struct {
    Code       string `json:"code"`
    Message    string `json:"message"`
    StatusCode int    `json:"-"`
    Err        error  `json:"-"` // wrapped original error
}

func (e *AppError) Error() string { return e.Message }
func (e *AppError) Unwrap() error { return e.Err }

// Konstruktor
func New(statusCode int, code, message string, err error) *AppError {
    return &AppError{StatusCode: statusCode, Code: code, Message: message, Err: err}
}

// Sentinel helpers
func NotFound(entity string, err error) *AppError {
    return New(http.StatusNotFound, "NOT_FOUND", entity+" tidak ditemukan", err)
}
func Unauthorized(err error) *AppError {
    return New(http.StatusUnauthorized, "UNAUTHORIZED", "Akses tidak diizinkan", err)
}
func Forbidden(err error) *AppError {
    return New(http.StatusForbidden, "FORBIDDEN", "Anda tidak memiliki hak akses", err)
}
func BadRequest(message string, err error) *AppError {
    return New(http.StatusBadRequest, "BAD_REQUEST", message, err)
}
func Conflict(message string, err error) *AppError {
    return New(http.StatusConflict, "CONFLICT", message, err)
}
func Internal(err error) *AppError {
    return New(http.StatusInternalServerError, "INTERNAL_ERROR", "Terjadi kesalahan server", err)
}
func ValidationError(message string) *AppError {
    return New(http.StatusUnprocessableEntity, "VALIDATION_ERROR", message, nil)
}
```

### 4.2 Global Error Handler Middleware

```go
// internal/middleware/error_handler.go
// Semua handler TIDAK boleh langsung c.JSON error — gunakan c.Error(err) lalu
// middleware ini yang handle di akhir chain.

func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        if len(c.Errors) == 0 {
            return
        }
        err := c.Errors.Last().Err
        var appErr *apperror.AppError
        if errors.As(err, &appErr) {
            // Log internal error tanpa expose ke client
            slog.Error("app error", "code", appErr.Code, "err", appErr.Err)
            c.JSON(appErr.StatusCode, response.Error(appErr.Code, appErr.Message))
            return
        }
        // Unhandled error
        slog.Error("unhandled error", "err", err)
        c.JSON(500, response.Error("INTERNAL_ERROR", "Terjadi kesalahan server"))
    }
}
```

### 4.3 Aturan Try-Catch (defer + recover) di Goroutine

Setiap goroutine yang diluncurkan (worker, WebSocket) **WAJIB** memiliki recover:

```go
func safeGo(fn func()) {
    go func() {
        defer func() {
            if r := recover(); r != nil {
                slog.Error("goroutine panic recovered", "panic", r)
            }
        }()
        fn()
    }()
}
```

### 4.4 Aturan Penulisan Error di Setiap Layer

| Layer | Aturan |
|---|---|
| **Repository** | Wrap error DB dengan `fmt.Errorf("methodName: %w", err)`. Kembalikan `apperror.NotFound` jika `pgx.ErrNoRows`. |
| **Service** | Tambah konteks bisnis. Jangan expose detail DB ke atas. |
| **Handler** | Panggil `c.Error(err)` — **jangan** `c.JSON` langsung untuk error. |
| **Middleware** | Tangkap via `ErrorHandler()` middleware. |

---

## 5. Keamanan (SECURITY — WAJIB)

### 5.1 Password: Bcrypt dengan Cost Factor

```go
// pkg/hash/bcrypt.go
package hash

import (
    "golang.org/x/crypto/bcrypt"
)

const BcryptCost = 12 // Minimum 12, jangan kurang

// HashPassword menghasilkan hash bcrypt yang sudah include salt otomatis.
// bcrypt secara internal men-generate random salt — JANGAN tambah salt manual.
func HashPassword(plainPassword string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), BcryptCost)
    if err != nil {
        return "", fmt.Errorf("hash.HashPassword: %w", err)
    }
    return string(bytes), nil
}

// VerifyPassword membandingkan plain password dengan hash.
// Mengembalikan nil jika cocok, error jika tidak.
func VerifyPassword(plainPassword, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plainPassword))
}
```

> **Catatan**: `bcrypt` secara otomatis men-generate dan menyimpan **salt** di dalam hash string (`$2a$12$<22-char-salt><31-char-hash>`). Kolom `password_hash` di tabel `users` menyimpan string lengkap ini. **Tidak perlu kolom salt terpisah.**

### 5.2 NIK: Enkripsi at-rest + Hash untuk pencarian

```go
// pkg/crypto/nik.go
// Gunakan AES-256-GCM untuk enkripsi NIK (symmetric, key dari env)
// Simpan hasil encrypt ke kolom nik_encrypted (BYTEA)
// Simpan SHA-256(NIK) ke kolom nik_hash (VARCHAR 64) untuk query unik
```

### 5.3 JWT

```go
// pkg/jwt/jwt.go
// Access token: expire 15 menit
// Refresh token: expire 7 hari, simpan hash-nya di DB atau Redis
// Gunakan HS256 dengan secret key minimal 32 byte dari env
// Claims wajib: sub (user_id), facility_id, roles, exp, iat, jti
```

### 5.4 Rate Limiting

- Login endpoint: maksimal **5 request/menit** per IP
- OTP endpoint: maksimal **3 request/5 menit** per nomor HP
- Gunakan middleware rate limiter (in-memory atau Redis)

### 5.5 Validasi Input

- **SEMUA** field request DTO **wajib** punya tag `validate:"..."`
- Jangan percaya input dari client — validasi di service layer juga jika perlu
- Sanitasi string untuk mencegah injection sebelum masuk ke query (pgx parameterized query sudah handle SQL injection)

---

## 6. Standard HTTP Response

```go
// pkg/response/response.go
type SuccessResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data"`
    Meta    *Meta       `json:"meta,omitempty"` // untuk pagination
}

type ErrorResponse struct {
    Success bool   `json:"success"`
    Code    string `json:"code"`
    Message string `json:"message"`
}

type Meta struct {
    Page    int `json:"page"`
    Limit   int `json:"limit"`
    Total   int `json:"total"`
}

func OK(data interface{}) SuccessResponse { ... }
func OKWithMeta(data interface{}, meta Meta) SuccessResponse { ... }
func Error(code, message string) ErrorResponse { ... }
```

**Format response sukses:**
```json
{
  "success": true,
  "data": { ... }
}
```

**Format response error:**
```json
{
  "success": false,
  "code": "NOT_FOUND",
  "message": "Pasien tidak ditemukan"
}
```

---

## 7. Database & Repository

### 7.1 Koneksi

- Gunakan `pgxpool` (connection pool), **bukan** `database/sql`
- Set `MaxConns`, `MinConns`, `MaxConnLifetime` di config
- Selalu pass `context.Context` ke semua query

### 7.2 Transaksi

```go
// Gunakan helper untuk transaksi agar tidak lupa commit/rollback
func WithTx(ctx context.Context, pool *pgxpool.Pool, fn func(pgx.Tx) error) error {
    tx, err := pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("begin tx: %w", err)
    }
    defer func() {
        if p := recover(); p != nil {
            _ = tx.Rollback(ctx)
            panic(p) // re-panic setelah rollback
        }
    }()
    if err := fn(tx); err != nil {
        _ = tx.Rollback(ctx)
        return err
    }
    return tx.Commit(ctx)
}
```

### 7.3 Query

- Gunakan **parameterized queries** (`$1, $2, ...`) — tidak boleh string concatenation
- Untuk pencarian teks (nama pasien), gunakan index `pg_trgm` + `ILIKE` atau `%`
- Nomor antrean di-generate dalam transaksi dengan `SELECT ... FOR UPDATE`

### 7.4 Migration

- Satu migration = satu perubahan logis
- Selalu buat file `.up.sql` dan `.down.sql`
- Jalankan migration otomatis saat startup (dev) atau via CI (prod)
- Nama file: `000001_init_schema.up.sql`, `000002_add_audit_log_index.up.sql`

---

## 8. Context & Timeout

```go
// Setiap handler wajib propagate context dari request
// Tambahkan timeout untuk operasi DB yang berat
ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
defer cancel()
result, err := s.repo.FindByID(ctx, id)
```

---

## 9. Logging

- Gunakan `log/slog` (structured logging, sudah built-in Go 1.21+)
- **Level:**
  - `slog.Debug` → detail flow (hanya dev)
  - `slog.Info` → event penting (request masuk, antrean berubah status)
  - `slog.Warn` → kondisi tidak normal tapi masih handled
  - `slog.Error` → error yang perlu perhatian
- **Jangan log** data sensitif: password, NIK asli, token JWT
- Format JSON di production, text di development

```go
slog.Info("patient registered", 
    "patient_id", patient.ID, 
    "facility_id", patient.FacilityID,
    "created_by", userID,
)
```

---

## 10. RBAC — Role & Permission

| Role | Kode | Akses |
|---|---|---|
| Admin Sistem | `admin` | Full akses semua endpoint |
| Petugas Pendaftaran | `registrar` | Kelola pasien, buat kunjungan, kelola antrean |
| Tenaga Medis | `medic` | Lihat antrean poli, catat kunjungan, diagnosis |
| Kepala Faskes | `head` | Dashboard analitik, audit log (read-only) |

```go
// Middleware RBAC: tambahkan ke route group
router.PATCH("/queues/:id/status", middleware.RequireRoles("medic", "registrar"), handler.UpdateQueueStatus)
```

---

## 11. Audit Log

Setiap operasi yang mengubah data pasien **wajib** ditulis ke tabel `audit_logs`:

```go
// Audit log ditulis ASYNCHRONOUSLY via channel agar tidak memperlambat response
type AuditEntry struct {
    UserID   uuid.UUID
    Action   string // "create" | "read" | "update" | "delete" | "login"
    Entity   string // "patient" | "queue" | "visit"
    EntityID uuid.UUID
    IP       string
    Meta     map[string]interface{}
}
```

---

## 12. WebSocket / SSE untuk Antrean Real-time

- Endpoint: `GET /ws/queues/:polyclinic_id`
- Broadcast ke semua client yang subscribe ke poli tersebut saat status antrean berubah
- Gunakan hub pattern (channel-based) agar goroutine-safe
- Setiap goroutine WebSocket **wajib** punya defer recover (lihat §4.3)

---

## 13. Aturan Pull Request / Commit

- **1 PR = 1 fitur atau 1 bug fix**
- Commit message: `feat(auth): implement JWT refresh token`, `fix(queue): handle concurrent number generation`
- Prefix: `feat`, `fix`, `refactor`, `test`, `chore`, `docs`
- Setiap fungsi public **wajib** punya Go doc comment
- Tambah unit test untuk semua service function (minimal happy path + 1 error case)

---

## 14. Environment Variables (Wajib ada di `.env`)

```env
APP_ENV=development          # development | production
APP_PORT=8080
APP_SECRET_KEY=              # JWT secret, min 32 char

DB_HOST=localhost
DB_PORT=5432
DB_NAME=antri_sehat
DB_USER=postgres
DB_PASSWORD=
DB_MAX_CONNS=20
DB_MIN_CONNS=2

NIK_ENCRYPTION_KEY=          # 32-byte hex untuk AES-256

BCRYPT_COST=12

JWT_ACCESS_EXPIRE=15m
JWT_REFRESH_EXPIRE=168h      # 7 hari

# Notifikasi (opsional MVP)
NOTIF_GATEWAY_URL=
NOTIF_GATEWAY_TOKEN=
```

---

## 15. Checklist Sebelum Merge

- [ ] Semua error di handler menggunakan `c.Error(err)`, bukan `c.JSON` langsung
- [ ] Semua goroutine punya defer recover
- [ ] Password tidak pernah di-log
- [ ] NIK tidak pernah dikirim ke client dalam bentuk plain
- [ ] Semua endpoint yang butuh auth punya middleware JWT
- [ ] Semua endpoint yang butuh role punya middleware RBAC
- [ ] Context di-propagate ke seluruh call chain (handler → service → repo)
- [ ] Ada migration file untuk setiap perubahan schema
- [ ] Unit test ditambahkan / diupdate
