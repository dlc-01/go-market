package register

import (
	"github.com/dlc/go-market/internal/auth"
	"github.com/dlc/go-market/internal/handlers"
	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/model"
	"github.com/dlc/go-market/internal/model/apperrors"
	"github.com/dlc/go-market/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(ginC *gin.Context) {
	var req model.AuthReq
	var err error
	handlers.BindData(ginC, &req)

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
	ginC.AbortWithStatus(http.StatusOK)

}
