package accrual

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/dlc/go-market/internal/config"
	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/model"
)

const ordersChanSize = 1000

var (
	zero          float64 = 0
	client                = resty.New()
	ordersRequest         = make(chan []model.Order, ordersChanSize)
	ordersSave            = make(chan model.Order, ordersChanSize)
)

func ListenOtherServ(ctx context.Context, cfg *config.ServerConfig) {
	logger.Info("start routine")
	poolTicker := time.NewTicker(time.Second * time.Duration(cfg.Poll))
	go collectOrders(ctx, ordersRequest, poolTicker)
	go workAccrual(ordersRequest, ordersSave, cfg, client)
	go saveOrder(ctx, ordersSave)
}

func workAccrual(ordersRequest chan []model.Order, ordersSave chan model.Order, cfg *config.ServerConfig, client *resty.Client) {
	for orders := range ordersRequest {
		for _, order := range orders {
			client.SetTimeout(time.Second)
			resp, err := client.R().Get(cfg.AccrualAddress + "/api/orders/" + order.ID)
			if err != nil {
				logger.Errorf("error while sending resp to accrual: %s", err)
				return
			}
			if resp.StatusCode() == http.StatusNoContent {
				logger.Infof("order %s not found in accrual ", order.ID)
				continue
			}
			if resp.StatusCode() == http.StatusTooManyRequests {
				logger.Info("too many requests")
				time.Sleep(60 * time.Second)
			}
			if resp.StatusCode() == http.StatusInternalServerError {
				logger.Info("Internal accrual server error")
				time.Sleep(2 * time.Second)
				continue
			}
			var buf bytes.Buffer
			buf.Write(resp.Body())

			var externalData model.Order
			err = json.Unmarshal(buf.Bytes(), &externalData)

			if err != nil {
				logger.Errorf("error while unmarshal data %s", err)
			}
			externalData.ID = order.ID
			ordersSave <- externalData

		}
	}
}
