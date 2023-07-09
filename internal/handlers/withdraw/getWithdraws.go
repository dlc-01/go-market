package withdraw

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dlc/go-market/internal/handlers"
	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/model/apperrors"
	"github.com/dlc/go-market/internal/storage"
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
		}
		logger.Errorf("error while getting all orders by login: %s", err)
		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ginC.AbortWithStatusJSON(http.StatusOK, &u.Withdraws)

}
