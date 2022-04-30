package utils

import (
	"os"
	"time"

	"github.com/Reynadi531/creatica2022-be/internal/entities"
	"github.com/golang-jwt/jwt/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

var (
	JWT_EXPIRE_TIME    = time.Now().Add(time.Hour * 48).Unix()
	JWT_SIGNATURE_KEY  = []byte(os.Getenv("JWT_SECRET"))
	JWT_SIGNING_METHOD = jwt.SigningMethodHS256
)

type MyClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	UserId   string `json:"userid"`
}

func GenerateJWTToken(user entities.User) (string, int64, error) {
	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "Pendekin",
			ExpiresAt: JWT_EXPIRE_TIME,
		},
		Username: user.Username,
		UserId:   user.ID.String(),
	}

	token := jwt.NewWithClaims(JWT_SIGNING_METHOD, claims)
	signedToken, err := token.SignedString(JWT_SIGNATURE_KEY)
	if err != nil {
		return "", 0, err
	}

	return signedToken, JWT_EXPIRE_TIME, nil
}

func GenerateRefreshToken() (string, error) {
	id, err := gonanoid.New()
	return id, err
}
