package accrual

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/dlc/go-market/internal/config"
	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/model"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

const ordersChanSize = 1000

var (
	zero    float64 = 0
	client          = resty.New()
	ordersR         = make(chan []model.Order, ordersChanSize)
	ordersS         = make(chan model.Order, ordersChanSize)
)

func ListenOtherServ(ctx context.Context, cfg *config.ServerConfig) {
	logger.Info("start routine")
	poolTicker := time.NewTicker(time.Second * time.Duration(cfg.Poll))
	go collectOrders(ctx, ordersR, poolTicker)
	go workAccrual(ordersR, ordersS, cfg, client)
	go saveOrder(ctx, ordersS)
}

func workAccrual(ordersR chan []model.Order, ordersS chan model.Order, cfg *config.ServerConfig, c *resty.Client) {
	for orders := range ordersR {
		for _, order := range orders {
			client.SetTimeout(time.Second)
			resp, err := client.R().Get(cfg.AccrualAddress + "/api/orders/" + order.ID)
			if err != nil {
				logger.Errorf("error while sending resp to accrual: %s", err)
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

			switch externalData.Status {
			case model.NEW:
				if order.Status == model.NEW {
					continue
				}
			case model.PROCESSING:
				if order.Status == model.NEW {
					externalData.Accrual = 0
					ordersS <- externalData

				} else {
					continue
				}

			case model.PROCESSED:
				externalData.Accrual = 0
				ordersS <- externalData

			case model.INVALID:
				externalData.Accrual = 0
				ordersS <- externalData
			}
			time.Sleep(1 * time.Second)
		}
	}
}
