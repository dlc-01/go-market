package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

var sLog zap.SugaredLogger

func InitLogger() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("cannot init logger: %w", err)
	}
	sLog = *logger.Sugar()
	return nil
}

func Fatalf(format string, opts ...any) {
	sLog.Fatalf(format, opts)
}

func Errorf(format string, opts ...any) {
	sLog.Errorf(format, opts)
}
func Error(format string) {
	sLog.Error(format)
}
func Infof(format string, opts ...any) {
	sLog.Infof(format, opts)
}

func Warnf(format string, opts ...any) {
	sLog.Warnf(format, opts)
}
func Panicf(format string, opts ...any) {
	sLog.Panicf(format, opts)
}

func Info(msg string) {
	sLog.Info(msg)
}

func GetMiddlewareLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		sLog.Infoln(
			"type", "request",
			"uri", c.Request.RequestURI,
			"method", c.Request.Method,
			"duration", latency,
		)
		sLog.Infoln(
			"type", "response",
			"status", c.Writer.Status(),
			"size", c.Writer.Size(),
		)
	}
}
