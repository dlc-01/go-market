package gzip

import (
	"compress/gzip"
	"github.com/dlc/go-market/internal/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression
)

func Gzip(level int) gin.HandlerFunc {
	return func(ginC *gin.Context) {
		headerContentGzip := strings.Contains(ginC.Request.Header.Get("Content-Encoding"), "gzip")
		headerAcceptGzip := strings.Contains(ginC.Request.Header.Get("Accept-Encoding"), "gzip")

		if headerContentGzip {
			newCompressReader(ginC)
			newCompressWriter(ginC, level)
		} else if headerAcceptGzip && !headerContentGzip {
			newCompressWriter(ginC, level)
		}
	}
}

func newCompressReader(ginC *gin.Context) {
	r, err := gzip.NewReader(ginC.Request.Body)
	if err != nil {
		ginC.AbortWithError(http.StatusBadRequest, err)
		logger.Errorf("cannot uncompressed request body: %s", err)
		return
	}

	ginC.Request.Body = r
	defer r.Close()

	ginC.Next()
}

func newCompressWriter(ginC *gin.Context, level int) {
	gz, err := gzip.NewWriterLevel(ginC.Writer, level)
	if err != nil {
		logger.Errorf("cannot compress request body: %s", err)
		return
	}

	ginC.Writer = &gzipWriter{ginC.Writer, gz}
	defer gz.Close()
	ginC.Header("Content-Encoding", "gzip")

	ginC.Next()
}

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}
