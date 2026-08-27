# Alur Perbaikan Auth (Admin → Dashboard, User → Landing)

Dokumen ini berisi **alur perbaikan + kode** agar:

1. Login role **admin** → ke `/dashboard`.
2. Login role **user** → ke `/` (landing page).
3. Role **user** tidak bisa buka `/dashboard`.
4. Login/register pages tidak bisa dibuka user yang sudah login.

Konvensi role di DB: `admin = 0`, `user = 1` (register pakai RoleID 1).

---

## Ringkasan Kesalahan yang Ada

| File | Kesalahan |
|------|-----------|
| `internal/usecase/auth_usecase.go:31` | `Login` return `(string, error)` — **tidak cocok** interface `(string, string, error)` → tidak compile |
| `internal/delivery/http/handler/auth_handler.go` | Struct `AuthHandler` tak punya field `config`, padahal `main.go:32` mem-pass `cg` → tidak compile |
| `auth_handler.go:44-47` | Login response `redirect` selalu `"/dashboard"` tanpa lihat role |
| `page_handler.go:28,39` | `LoginPages`/`RegisterPages` redirect selalu `/dashboard` saat sudah login (bukan berdasar role) |
| `route.go:32` | `/dashboard` cuma `middleware.Auth`, tak cek role admin |

---

## Alur Setelah Perbaikan

```
Login (POST /api/v1/auth/login)
   │
   ▼
AuthUsecase.Login
   │ cek email+password (FindByEmail + bcrypt)
   │ generate token berisi RoleID
   │ pilih redirect berdasar RoleID
   │   admin(0) -> "/dashboard"
   │   user(1)  -> "/"
   ▼
return (token, redirect, nil)
   │
   ▼
AuthHandler.Login
   │ set cookie "token"
   │ response.Success({ token, redirect })
   ▼
JSON ke browser
   │
   ▼
app.js:52  window.location.href = data.data.redirect || "/dashboard"
   │
   ├─ admin -> /dashboard   (middleware RequireAdmin → lolos)
   └─ user  -> /            (landing, OptionalAuth)

Saat user buka /dashboard langsung:
   └─ RequireAdmin cek RoleID == 1 → redirect "/" (TOLAK)

Saat sudah login buka /login atau /register:
   └─ handler cek cookie valid
         admin -> redirect "/dashboard"
         user  -> redirect "/"
```

---

## Kode Perbaikan

### 1. `internal/usecase/auth_usecase.go` — role-aware redirect

```go
func (u *AuthUsecaseImpl) Login(req du.LoginRequest) (string, string, error) {
	user, err := u.userRepository.FindByEmail(req.Email)
	if err != nil {
		return "", "", errors.New("Invalid email or password")
	}

	if !helper.CheckPassword(req.Password, user.Password) {
		return "", "", errors.New("Invalid email or password")
	}

	token, err := helper.GenerateToken(user.ID, user.Email, user.Firstname, user.Lastname, user.RoleID, u.config.JWT)
	if err != nil {
		return "", "", err
	}

	// pilih redirect berdasar role
	redirect := "/" // user -> landing page
	if user.RoleID == 0 {
		redirect = "/dashboard" // admin -> dashboard
	}

	return token, redirect, nil
}
```

### 2. `internal/delivery/http/handler/auth_handler.go` — tambah config + pakai redirect

Ubah struct & constructor:

```go
type AuthHandler struct {
	authUsecase du.AuthUsecase
	cfg         *config.Config
}

func NewAuthHandler(authUsecase du.AuthUsecase, cfg *config.Config) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase, cfg: cfg}
}
```

Ubah `Login` (buang hardcode `/dashboard`):

```go
func (h *AuthHandler) Login(c *gin.Context) {
	var req du.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := validator.Validate(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	token, redirect, err := h.authUsecase.Login(req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.SetCookie("token", token, 3600*24, "/", "", false, true)

	response.Success(c, http.StatusOK, gin.H{
		"token":    token,
		"redirect": redirect,
	})
}
```

Tambah import `"Project/internal/config"`.

### 3. `internal/delivery/http/handler/page_handler.go` — redirect role-aware saat sudah login

Ubah `LoginPages` & `RegisterPages` (ganti cek cookie saja jadi cek role):

```go
func (h *PageHandler) homeByRole(c *gin.Context) {
	role, _ := c.Get("role")
	if r, ok := role.(uint); ok && r == 0 {
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}
	c.Redirect(http.StatusSeeOther, "/")
}

func (h *PageHandler) LoginPages(c *gin.Context) {
	if _, err := c.Cookie("token"); err == nil {
		h.homeByRole(c)
		return
	}
	c.HTML(http.StatusOK, "login", gin.H{
		"title": "Login BooCinS",
		"nav":   false,
	})
}

func (h *PageHandler) RegisterPages(c *gin.Context) {
	if _, err := c.Cookie("token"); err == nil {
		h.homeByRole(c)
		return
	}
	c.HTML(http.StatusOK, "register", gin.H{
		"title": "Register BooCinS",
		"nav":   false,
	})
}
```

> Karena `/login` & `/register` tidak lewat middleware auth, `c.Get("role")` kosong.
> Solusi: route ini perlu middleware ringan yang decode role (lihat langkah 4).

**Versi lebih bersih:** beri middleware `OptionalAuth` pada `/login` dan
`/register` supaya `c.Get("role")` terisi, lalu `homeByRole` di atas berfungsi.
Ubah route:

```go
r.GET("/login", middleware.OptionalAuth(cfg.JWT), pageHandler.LoginPages)
r.GET("/register", middleware.OptionalAuth(cfg.JWT), pageHandler.RegisterPages)
```

### 4. `internal/delivery/http/middleware/auth_middleware.go` — RequireAdmin

```go
// RequireAdmin : hanya role admin (RoleID == 0) yang boleh akses.
// Kalau user (RoleID != 0) -> redirect ke landing page.
func RequireAdmin(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}

		if claims.RoleID != 0 {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("firstname", claims.Firstname)
		c.Set("lastname", claims.Lastname)
		c.Set("role", claims.RoleID)

		c.Next()
	}
}
```

### 5. `internal/delivery/http/router/route.go` — pasang middleware

```go
r.GET("/", middleware.OptionalAuth(cfg.JWT), pageHandler.LandingPages)
r.GET("/login", middleware.OptionalAuth(cfg.JWT), pageHandler.LoginPages)
r.GET("/register", middleware.OptionalAuth(cfg.JWT), pageHandler.RegisterPages)

r.POST("/api/v1/auth/login", authHandler.Login)
r.POST("/api/v1/auth/logout", authHandler.Logout)
r.POST("/api/v1/auth/register", authHandler.Register)

r.GET("/dashboard", middleware.RequireAdmin(cfg.JWT), pageHandler.DashboardPage)
```

---

## Pengecekan (Setelah Perbaikan)

1. `go build ./...` — harus compile tanpa error.
2. Login admin → `redirect: "/dashboard"` → buka dashboard (lolos).
3. Login user → `redirect: "/"` → landing.
4. User akses `/dashboard` langsung → di-302 ke `/`.
5. User login buka `/login` → di-302 ke `/`.
6. Admin login buka `/login` → di-302 ke `/dashboard`.

---

## Catatan

- `main.go:32` sudah mem-pass `cg` ke `NewAuthHandler` → setelah struct diperbaiki (langkah 2) konsisten, tak perlu ubah main.go.
- Controller `AuthUsecaseImpl.Login` harus return 3 nilai agar cocok interface `du.AuthUsecase` (auth_usecase.go:4).
- Nilai role hardcode (`0` admin, `1` user). Kalau role bertambah, pindah ke konstanta/config.
