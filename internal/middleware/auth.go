package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	pkgjwt "sespima_api/pkg/jwt"
	"sespima_api/pkg/response"
)

func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.Unauthorized(c, "authorization header required")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		if tokenStr == header {
			response.Unauthorized(c, "invalid token format")
			c.Abort()
			return
		}

		claims, err := pkgjwt.ParseAccessToken(jwtSecret, tokenStr)
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", float64(claims.UserID))
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("nrp", claims.NRP)
		c.Set("nip", claims.NIP)
		c.Next()
	}
}
