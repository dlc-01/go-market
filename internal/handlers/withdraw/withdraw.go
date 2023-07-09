package withdraw

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/dlc/go-market/internal/handlers"
	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/luhn"
	"github.com/dlc/go-market/internal/model"
	"github.com/dlc/go-market/internal/model/apperrors"
	"github.com/dlc/go-market/internal/storage"
)

func WithdrawalOfFunds(ginC *gin.Context) {
	req := model.Withdraw{TimeCreated: time.Now()}

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

	if err := luhn.ValidIDErr(req.Order); err != nil {
		logger.Errorf("error while checking by algorithm luhn %s", err)
		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
	if req.Sum <= 0 {
		err := apperrors.NewUnprocessableContent("sum <= 0 ")
		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	log, err := handlers.GetLogin(ginC)
	if err != nil {
		logger.Error(err.Error())
		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	mapOrd := make([]model.Order, 1)
	mapOrd[0] = model.Order{
		ID:          req.Order,
		Status:      model.PROCESSED,
		Accrual:     0,
		TimeCreated: time.Now(),
	}

	withdraws := make([]model.Withdraw, 1)
	withdraws[0] = req

	u := model.User{
		Info: model.UserInfo{
			Login: log},
		Orders:    mapOrd,
		Withdraws: withdraws,
		Balance:   req.Sum}
	if err = storage.AddNewOderWithdraw(ginC, &u); err != nil {
		if apperrors.Status(err) < 300 {
			ginC.AbortWithStatus(apperrors.Status(err))
			return
		}
		logger.Errorf("error while adding new order: %s", err)
		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return

	}

}
