package accrual

import (
	"context"
	"github.com/dlc/go-market/internal/config"
	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/model"
	"github.com/dlc/go-market/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func TestListenOtherServ(t *testing.T) {
	tests := []struct {
		name         string
		status       int
		statusAcrual int
		orderAccrual model.Order
		OrderStor    []model.Order
	}{
		{
			name:         "first",
			status:       http.StatusOK,
			statusAcrual: http.StatusOK,
			orderAccrual: model.Order{ID: "666", Accrual: 555, Status: model.PROCESSED},
			OrderStor:    []model.Order{{ID: "666", Accrual: 0, Status: model.NEW}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.ServerConfig{AccrualAddress: "http://localhost:8081", Poll: 10}
			gin.SetMode(gin.TestMode)
			route := gin.Default()
			route.GET("/api/orders/:id", func(c *gin.Context) {
				if tt.statusAcrual != http.StatusOK {
					c.AbortWithStatus(tt.statusAcrual)
				} else {
					c.AbortWithStatusJSON(tt.statusAcrual, tt.orderAccrual)
				}

			})
			route.Run(cfg.AccrualAddress)

			mockUserService := new(storage.TestStore)
			mockUserService.On("CollectOrders", mock.AnythingOfType("*context.emptyCtx")).Return(tt.OrderStor, nil)
			mockUserService.On("UpdateOrders", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("model.Order")).Return(tt.orderAccrual, nil)
			logger.InitLogger()
			storage.InitTestStorage(mockUserService)
			ListenOtherServ(context.TODO(), &cfg)
		})
	}
}
