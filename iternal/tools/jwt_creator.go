package tools

import (
	"Arkadiy_Servis_authorization/config"
	"Arkadiy_Servis_authorization/iternal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func JWTCreator(user domain.User, c *gin.Context) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Arkadiy_Service_Authorization",
			Subject:   "user",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:     user.ID,
		UserStatus: user.Status,
		UserPerm:   user.Permissions,
	})

	strToken, err := token.SignedString([]byte(config.Env.SecretKey))
	if err != nil {
		return err
	}

	c.Writer.Header().Set("Authorization", "Bearer "+strToken)

	return nil
}
