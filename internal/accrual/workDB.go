package accrual

import (
	"context"
	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/model"
	"github.com/dlc/go-market/internal/storage"
	"time"
)

func collectOrders(ctx context.Context, ordersC chan<- []model.Order, poolTicker *time.Ticker) {
	for range poolTicker.C {
		orders, err := storage.CollectOrders(ctx)
		if err != nil {
			logger.Errorf("cannot collect orders: %s", err)
		}
		logger.Info("get orders")
		ordersC <- orders
	}
}

func saveOrder(ctx context.Context, ordersS <-chan model.Order) {
	for order := range ordersS {
		if err := storage.UpdateOrders(ctx, order); err != nil {
			logger.Errorf("cannot update order :%s", err)
		}
	}
}
