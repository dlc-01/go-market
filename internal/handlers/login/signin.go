package login

import (
	"github.com/dlc/go-market/internal/auth"
	"github.com/dlc/go-market/internal/handlers"
	"github.com/dlc/go-market/internal/hash"
	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/model"
	"github.com/dlc/go-market/internal/model/apperrors"
	"github.com/dlc/go-market/internal/storage"
	"github.com/gin-gonic/gin"
)

func Login(ginC *gin.Context) {
	var req model.AuthReq

	handlers.BindData(ginC, &req)

	u := &model.UserInfo{Login: req.Username, Password: req.Password}

	u, err := storage.FindByLogin(ginC.Request.Context(), &u.Login)
	if err != nil {
		logger.Errorf("cannot find user %s", err)
		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return

	}

	if err := hash.CheckPassword(req.Password, u.Password); err != nil {
		logger.Infof("wrong password %s", err)
		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
	auth.SetToken(ginC, *u)

}
