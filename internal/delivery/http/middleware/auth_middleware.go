package middleware

import (
	"Project/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware membaca cookie token, decode JWT, lalu
// menyimpan data user ke context.
func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1) Ambil token dari cookie
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			c.Abort()
			return
		}

		// 2) Validate & decode token
		claims, err := helper.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
			})
			c.Abort()
			return
		}

		// 3) Simpan data user ke context
		//    -> SEMUA handler di route yang memakai middleware ini
		//       bisa ambil lewat c.Get("user_id"), c.Get("email"), dll.
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.RoleID)
		c.Set("firstname", claims.Firstname)
		c.Set("lastname", claims.Lastname)

		// 4) Lanjut ke handler berikutnya
		c.Next()
	}
}

func OptionalAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")

		if err != nil {
			c.Next()
			return
		}

		claims, err := helper.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			c.Next()
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
