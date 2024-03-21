package middleware

import (
	"Arkadiy_Servis_authorization/iternal/tools"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Timer() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// продолжаем работать с хэндлером который идет после мидлвейра  (который был вызван изначально))
		c.Next()

		latency := time.Since(t)
		log.WithField("component", "latency").Info(latency)
	}
}

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		strToken := c.Request.Header.Get("Authorization")

		claims, err := tools.ParsTokenClaims(strToken)

		if err != nil {
			tools.CreateError(http.StatusUnauthorized, errors.New("You're Unauthorized"), c)
			return
		}

		if claims.UserStatus != "confirmed" {
			tools.CreateError(http.StatusUnauthorized, errors.New("ваш статус не подтвержден"), c)
			return
		}

		c.Next()
	}
}

func AdminAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		strToken := c.Request.Header.Get("Authorization")

		claims, err := tools.ParsTokenClaims(strToken)

		if err != nil {
			tools.CreateError(http.StatusUnauthorized, errors.New("You're Unauthorized"), c)
			return
		}

		if claims.UserPerm != 3 {
			tools.CreateError(http.StatusUnauthorized, errors.New("You don't have access rights"), c)
			return
		}

		c.Next()
	}
}
