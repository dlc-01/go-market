package accrual

import (
	"context"
	"time"

	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/model"
	"github.com/dlc/go-market/internal/storage"
)

func collectOrders(ctx context.Context, ordersRequest chan<- []model.Order, poolTicker *time.Ticker) {
	for range poolTicker.C {
		orders, err := storage.CollectOrders(ctx)
		if err != nil {
			logger.Errorf("cannot collect orders: %s", err)
		}
		logger.Info("get orders")
		ordersRequest <- orders
	}
}

func saveOrder(ctx context.Context, ordersSave <-chan model.Order) {
	for order := range ordersSave {
		if err := storage.UpdateOrders(ctx, order); err != nil {
			logger.Errorf("cannot update order :%s", err)
		}
	}
}
