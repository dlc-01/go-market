package orders

import (
	"bytes"
	"github.com/dlc/go-market/internal/handlers"
	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/luhn"
	"github.com/dlc/go-market/internal/model"
	"github.com/dlc/go-market/internal/model/apperrors"
	"github.com/dlc/go-market/internal/storage"
	"github.com/gin-gonic/gin"
	"time"
)

func NewOrder(ginC *gin.Context) {
	if err := handlers.CheckContentType(ginC, "text/plain"); err != nil {
		logger.Info("content type not support")
		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(ginC.Request.Body)
	if err != nil {
		logger.Errorf("error while reading req body: %s", err)
		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
	orderS := buf.String()

	if err = luhn.ValidIDErr(orderS); err != nil {
		if apperrors.Status(err) == 500 {
			logger.Error(err.Error())
		} else {
			logger.Info(err.Error())
		}
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

	var acc float64 = 0
	mapOrd := make([]model.Order, 1)
	mapOrd[0] = model.Order{
		Id:          orderS,
		Status:      model.NEW,
		Accrual:     &acc,
		TimeCreated: time.Now(),
	}
	u := model.User{Info: model.UserInfo{Login: log}, Orders: mapOrd}

	if err = storage.AddNewOrder(ginC, &u); err != nil {
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
