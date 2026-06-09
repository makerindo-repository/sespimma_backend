package helpers

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id uint) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": id,
			"exp": time.Now().Add(
				time.Hour * 24,
			).Unix(),
		})

	return token.SignedString(
		[]byte(os.Getenv("JWT_SECRET")),
	)
}