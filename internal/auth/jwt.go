package auth

import (
	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

var secretK string

func SetSecretKey(s string) {
	secretK = s
}
func SetToken(ginC *gin.Context, u model.UserInfo) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": u.Login,
	})
	tokenString, err := t.SignedString([]byte(secretK))
	if err != nil {
		logger.Errorf("failed to create token: %w", err)
		return
	} else {
		ginC.SetSameSite(http.SameSiteLaxMode)
		ginC.SetCookie("Authorise", tokenString, 3600, "", "", false, true)
		return
	}
}
