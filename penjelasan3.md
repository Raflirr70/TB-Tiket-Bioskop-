# Penjelasan 3
# Redirect Berdasarkan Role Setelah Login (Admin vs User)

Dokumentasi ini menjelaskan langkah-langkah agar:

1. Saat login dengan role **admin** → masuk ke halaman **dashboard**.
2. Saat login dengan role **user** → kembali ke halaman **landing page**.
3. Role **user** **tidak bisa** mengakses halaman dashboard.

---

## 1. Kondisi Saat Ini (Masalah)

Cek `internal/delivery/http/handler/auth_handler.go:35-47` (Login):

```go
token, err := h.authUsecase.Login(req)
if err != nil {
	response.Error(c, http.StatusUnauthorized, err.Error())
	return
}

// Set cookie "token" untuk browser
c.SetCookie("token", token, 3600*24, "/", "", false, true)

response.Success(c, http.StatusOK, gin.H{
	"token":    token,
	"redirect": "/dashboard",   // <-- SELALU ke dashboard
})
```

Dan `internal/delivery/http/handler/page_handler.go:49-58`:

```go
func (h *PageHandler) DashboardPage(c *gin.Context) {
	...
	c.HTML(http.StatusOK, "base", gin.H{
		"page": "dashboard",
		...
	})
}
```

Masalah:
- Response login **selalu** `redirect: "/dashboard"` tanpa lihat role.
  Baik admin maupun user sama-sama diarahkan ke dashboard.
- Route `/dashboard` (route.go:32) cuma pakai `middleware.Auth` yang
  hanya mengecek **apakah login**, tidak mengecek **role apa**.

---

## 2. Konsep Role di Project Ini

Berdasarkan konvensi database (register pakai `RoleID: 1`, lihat
`internal/usecase/auth_usecase.go:60`):

| Role             | RoleID |
|------------------|--------|
| Admin            | 0      |
| User (default)   | 1      |

Role disimpan dalam token JWT (pkg/helper/jwt.go:11):
```go
RoleID uint `json:"role"`
```

Lalu dibaca oleh middleware dan disimpan ke context (auth_middleware.go:39):
```go
c.Set("role", claims.RoleID)
```

Jadi handler bisa ambil role dengan:
```go
role, _ := c.Get("role")  // uint (0 = admin, 1 = user)
```

---

## 3. Alur yang Diinginkan

```
User login (POST /api/v1/auth/login)
        │
        ▼
  AUTH USECASE LOGIN
  (cek email + password, generate token berisi RoleID)
        │
        ▼
  AUTH HANDLER LOGIN
        │
   cek role dari token
        │
   ┌────┴───────────┐
   │ role == 0      │ role == 1
   │ (admin)        │ (user)
   └────┬───────────┘
        │               │
        ▼               ▼
   /dashboard      / (landing page)

   USER akses /dashboard langsung
        │
        ▼
   MIDDLEWARE RequireAdmin
        │
   role == 1 (user)?
        │ (ya)
        ▼
   redirect ke landing page (TOLAK akses dashboard)
```

---

## 4. Implementasi

### LANGKAH 1 — Middleware `RequireAdmin` (memblokir user)

Buat middleware baru di `internal/delivery/http/middleware/auth_middleware.go`.
Fungsinya: hanya role **admin (0)** yang boleh lanjut; selain itu redirect
ke landing page.

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
			// Bukan admin (role user 1, dst) -> tolak akses dashboard
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

> Kalau mau pakai middleware `Auth` yang sudah ada di atasnya, satukan
> logika (cek login dulu, baru cek role). Untuk mudah baca, contoh di atas
> sudah menangani cek cookie + validasi + cek role sekaligus.

### LANGKAH 2 — Update route /dashboard

Ubah `internal/delivery/http/router/route.go:32`:

```go
// sebelum: cuma cek login
r.GET("/dashboard", middleware.Auth(cfg.JWT), pageHandler.DashboardPage)

// sesudah: cek login + cek role admin
r.GET("/dashboard", middleware.RequireAdmin(cfg.JWT), pageHandler.DashboardPage)
```

Sekarang user biasa (role 1) yang membuka `/dashboard` akan diarahkan
kembali ke `/` (landing page).

### LANGKAH 3 — Login redirect berdasar role (AuthHandler.Login)

Ubah `internal/delivery/http/handler/auth_handler.go:25-48` supaya
mementukan `redirect` dari role di dalam token.

Cara termudah: **decode token yang baru dibuat** untuk dapat RoleID,
lalu pilih redirect. Tambahkan helper `DecodeRole` atau decode langsung
pakai `helper.ValidateToken`.

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

	token, err := h.authUsecase.Login(req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Set cookie "token"
	c.SetCookie("token", token, 3600*24, "/", "", false, true)

	// Decode token untuk tahu role user yang login
	claims, _ := helper.ValidateToken(token, h.cfg.JWT) // butuh akses ke config

	redirect := "/" // default: user -> landing page
	if claims != nil && claims.RoleID == 0 {
		redirect = "/dashboard" // admin -> dashboard
	}

	response.Success(c, http.StatusOK, gin.H{
		"token":    token,
		"redirect": redirect,
	})
}
```

> Catatan: `AuthHandler` saat ini (`auth_handler.go:12-17`) belum menyimpan
> `config`. Supaya handler bisa decode token, tambahkan field `config
> *config.Config` pada struct `AuthHandler` dan isi saat `NewAuthHandler`
> dipanggil di `cmd/server/main.go`. Alternatif lebih bersih: buat usecase
> login mengembalikan role juga, atau buat fungs `GetRedirect(role uint)`.

Contoh alternatif (paling bersih): pindahkan logika penentuan redirect ke
**usecase** agar handler tidak perlu config. Ubah return `Login` menjadi
memberikan `redirect`:

```go
// di AuthUsecaseImpl.Login, ganti return (string, error) menjadi (string, string, error):
func (u *AuthUsecaseImpl) Login(req du.LoginRequest) (string, string, error) {
	...
	if user.RoleID == 0 { // admin
		return token, "/dashboard", nil
	}
	return token, "/", nil // user -> landing page
}
```

---

## 5. Ringkasan Perubahan

| File                        | Perubahan                                                  |
|-----------------------------|------------------------------------------------------------|
| auth_middleware.go          | Tambah middleware `RequireAdmin` (tolak non-admin)         |
| route.go:32                 | `/dashboard` pakai `middleware.RequireAdmin`               |
| auth_handler.go (Login)     | Redirect `/dashboard` utk admin, `/` utk user              |

Hasil akhir:

| Role | Login redirect | Akses /dashboard |
|------|----------------|------------------|
| Admin (0) | /dashboard | Boleh |
| User (1)  | / (landing) | Ditolak → redirect / |

---

> Catatan: nilai role (admin=0, user=1) di hardcode dalam contoh. Kalau
> nanti role bertambah atau berubah, pindahkan ke konstanta/config supaya
> mudah dirawat.
