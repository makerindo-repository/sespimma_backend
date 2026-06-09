package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {

	return func(c *gin.Context) {

		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized,
				gin.H{
					"message": "unauthorized",
				})
			c.Abort()
			return
		}

		c.Next()
	}
}