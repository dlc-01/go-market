package app

import (
	"github.com/dlc/go-market/internal/auth"
	"github.com/dlc/go-market/internal/config"
	"github.com/dlc/go-market/internal/gzip"
	"github.com/dlc/go-market/internal/handlers/balance"
	"github.com/dlc/go-market/internal/handlers/login"
	"github.com/dlc/go-market/internal/handlers/orders"
	"github.com/dlc/go-market/internal/handlers/register"
	"github.com/dlc/go-market/internal/handlers/withdraw"
	"github.com/dlc/go-market/internal/logger"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.ServerConfig) {
	auth.SetSecretKey(cfg.SecretKey)
	router := gin.Default()

	router.Use(logger.GetMiddlewareLogger())
	router.Use(gzip.Gzip(gzip.BestSpeed))

	user := router.Group("/")
	user.POST("/api/user/register", register.Register)
	user.POST("/api/user/login", login.Login)

	router.Use(auth.AuthMidlleware())

	router.POST("/api/user/orders", orders.NewOrder)
	router.GET("/api/user/orders", orders.GetAllOrders)

	router.GET("/api/user/balance", balance.ShowBalance)
	router.POST("/api/user/balance/withdraw", withdraw.WithdrawalOfFunds)
	router.GET("/api/user/withdraws", withdraw.GetAllWithdraws)
	router.Run(cfg.ServerAddress)
}
