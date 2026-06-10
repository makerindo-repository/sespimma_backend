package middleware

import (
	"github.com/gin-gonic/gin"
	"sespima_api/pkg/response"
)

// RequireRoles returns a middleware that aborts with 403 if the authenticated
// user's role is not in the allowed list.
func RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			response.Forbidden(c, "akses ditolak")
			c.Abort()
			return
		}
		if _, ok := allowed[role.(string)]; !ok {
			response.Forbidden(c, "akses ditolak: role tidak diizinkan")
			c.Abort()
			return
		}
		c.Next()
	}
}
