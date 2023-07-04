package handlers

import (
	"fmt"
	"github.com/dlc/go-market/internal/model/apperrors"
	"github.com/gin-gonic/gin"
)

func CheckContentType(ginC *gin.Context, types string) error {
	if ginC.ContentType() != types {
		return apperrors.NewUnsupportedMediaType(fmt.Sprintf("%s only accepts Content-Type %s", ginC.FullPath(), types))
	}
	return nil
}

func GetLogin(ginC *gin.Context) (string, error) {
	log, c := ginC.Get("login")
	if !c {
		return log.(string), apperrors.NewUnauthorized("cannot find login")
	}
	return log.(string), nil
}
