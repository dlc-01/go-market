package register

import (
	"github.com/dlc/go-market/internal/auth"
	"github.com/dlc/go-market/internal/handlers"
	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/model"
	"github.com/dlc/go-market/internal/model/apperrors"
	"github.com/dlc/go-market/internal/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(ginC *gin.Context) {
	var req model.AuthReq
	var err error
	if invalidArgs, err := handlers.BindData(ginC, &req); err != nil {
		logger.Infof("cannot bind data %s", err)
		if apperrors.Status(err) == http.StatusBadRequest {
			ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
				"error":       err,
				"InvalidArgs": invalidArgs,
			})
		} else {
			ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
				"error": err,
			})
		}

		return
	}

	u := model.UserInfo{Password: req.Password, Login: req.Username}

	err = storage.CreateUser(ginC, &u)
	if err != nil {
		logger.Errorf("cannot create user %s", err)
		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	auth.SetToken(ginC, u)
	ginC.AbortWithStatus(apperrors.Status(apperrors.NewStatusOK()))

}
