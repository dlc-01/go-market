package balance

import (
	"github.com/dlc/go-market/internal/handlers"
	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/model"
	"github.com/dlc/go-market/internal/model/apperrors"
	"github.com/dlc/go-market/internal/storage"
	"github.com/gin-gonic/gin"
)

func ShowBalance(ginC *gin.Context) {
	log, err := handlers.GetLogin(ginC)
	if err != nil {
		logger.Error(err.Error())
		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	var resp *model.BalanceResp

	if resp, err = storage.GetUserBalanceWithdraw(ginC, &log); err != nil {
		logger.Error(err.Error())
		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ginC.AbortWithStatusJSON(apperrors.Status(apperrors.NewStatusOK()), &resp)
}
