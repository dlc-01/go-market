package withdraw

import (
	"github.com/dlc/go-market/internal/handlers"
	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/luhn"
	"github.com/dlc/go-market/internal/model"
	"github.com/dlc/go-market/internal/model/apperrors"
	"github.com/dlc/go-market/internal/storage"
	"github.com/gin-gonic/gin"
	"time"
)

func WithdrawalOfFunds(ginC *gin.Context) {
	req := model.Withdraw{TimeCreated: time.Now()}

	handlers.BindData(ginC, &req)

	if err := luhn.ValidIDErr(req.Order); err != nil {
		if apperrors.Status(err) == 500 {
			logger.Error(err.Error())
		}
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
	}

	log, err := handlers.GetLogin(ginC)
	if err != nil {
		logger.Error(err.Error())
		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	balance, err := storage.GetUserBalance(ginC, &log)
	if err != nil {
		logger.Error(err.Error())
		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	if balance < req.Sum {
		err := apperrors.NewPaymentRequired("balance < sum of order")
		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	var acc float64 = 0
	mapOrd := make([]model.Order, 1)
	mapOrd[0] = model.Order{
		ID:          req.Order,
		Status:      model.PROCESSED,
		Accrual:     &acc,
		TimeCreated: time.Now(),
	}

	withdraws := make([]model.Withdraw, 1)
	withdraws[0] = req

	u := model.User{Info: model.UserInfo{Login: log}, Orders: mapOrd, Withdraws: withdraws, Balance: balance - req.Sum}
	if err = storage.AddNewOderWithdraw(ginC, &u); err != nil {
		if apperrors.Status(err) < 300 {
			ginC.AbortWithStatus(apperrors.Status(err))
			return
		} else if apperrors.Status(err) == 500 {
			logger.Errorf("error while adding new order: %s", err)
		}
		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return

	}

}
