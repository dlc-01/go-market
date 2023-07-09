package auth

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/model"
	"github.com/dlc/go-market/internal/model/apperrors"
)

func AuthMidlleware() gin.HandlerFunc {
	return func(ginC *gin.Context) {
		cookie, _ := ginC.Cookie("Authorise")
		if cookie == "" {
			logger.Info("cannot find cookie \"Authorise\"")
			err := apperrors.NewUnauthorized("user not logged in")
			ginC.AbortWithStatusJSON(apperrors.Status(err), err)
		} else {
			u, err := GetData(secretK, cookie)
			if err != nil {
				logger.Errorf("cannot get claims from token: %s ", err)
				ginC.AbortWithStatusJSON(apperrors.Status(apperrors.NewUnauthorized("")), gin.H{
					"error": err,
				})
			}
			ginC.Set("login", u.Login)
		}
		ginC.Next()

	}

}

func GetData(secretK string, cookie string) (*model.UserInfo, error) {
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			res := fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])
			return nil, errors.New(res)
		}
		return []byte(secretK), nil
	})

	if err != nil {
		return nil, err
	}
	return getClaims(token)

}

func getClaims(token *jwt.Token) (*model.UserInfo, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		u := &model.UserInfo{Login: claims["login"].(string)}

		return u, nil

	} else {
		return &model.UserInfo{}, fmt.Errorf("invalid token")
	}

}
