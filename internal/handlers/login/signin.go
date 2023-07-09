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
	"net/http"
)

func Login(ginC *gin.Context) {
	var req model.AuthReq

	if invalidArgs, err := handlers.BindData(ginC, &req); err != nil {
		logger.Errorf("cannot bind data %s", err)
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

	u := &model.UserInfo{Login: req.Username}

	u, err := storage.FindByLogin(ginC.Request.Context(), &u.Login)
	if err != nil {
		logger.Errorf("cannot find user %s", err)
		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return

	}

	if err := hash.CheckPassword(req.Password, u.Password); err != nil {
		logger.Errorf("wrong password %s", err)
		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
	auth.SetToken(ginC, *u)

}
