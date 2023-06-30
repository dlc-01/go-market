package withdraw

import (
	"github.com/dlc/go-market/internal/handlers"
	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/model/apperrors"
	"github.com/dlc/go-market/internal/storage"
	"github.com/gin-gonic/gin"
)

func GetAllWithdraws(ginC *gin.Context) {
	log, err := handlers.GetLogin(ginC)
	if err != nil {
		logger.Error(err.Error())
		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	u, err := storage.GetAllWithdrawsByLogin(ginC, &log)
	if err != nil {
		if apperrors.Status(err) < 300 {
			ginC.AbortWithStatus(apperrors.Status(err))
			return
		} else if apperrors.Status(err) == 500 {
			logger.Errorf("error while getting all orders by login: %s", err)
		}
		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ginC.AbortWithStatusJSON(apperrors.Status(apperrors.NewStatusOK()), &u.Withdraws)

}