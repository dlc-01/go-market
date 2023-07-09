package balance

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dlc/go-market/internal/auth"
	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/model"
	"github.com/dlc/go-market/internal/storage"
)

func TestGetBalance(t *testing.T) {

	tests := []struct {
		name        string
		status      int
		secretK     string
		user        *model.UserInfo
		balanceStor *model.BalanceResp
		balanceRes  *model.BalanceResp
	}{
		{
			name: "all good",
			user: &model.UserInfo{
				Login: "bob",
			},

			secretK:     "supersecret1234",
			status:      http.StatusOK,
			balanceStor: &model.BalanceResp{Balance: 0, Sum: 0},
			balanceRes:  &model.BalanceResp{Balance: 0, Sum: 0},
		},
		{
			name: "not true jwt token",
			user: &model.UserInfo{
				Login: "bob",
			},

			secretK:     "nosupersecret1234",
			status:      http.StatusUnauthorized,
			balanceStor: &model.BalanceResp{Balance: 0, Sum: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.InitLogger()
			gin.SetMode(gin.TestMode)
			if err := logger.InitLogger(); err != nil {
				log.Fatal(err)
			}

			mockUserService := new(storage.TestStore)
			mockUserService.On("GetBalanceWithdraw", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*string")).Return(tt.balanceStor, nil)

			rr := httptest.NewRecorder()
			storage.InitTestStorage(mockUserService)
			router := gin.Default()
			auth.SetSecretKey("supersecret1234")
			router.Use(auth.AuthMidlleware())
			router.GET("/balance", ShowBalance)

			reqBody, err := json.Marshal(tt.balanceRes)
			assert.NoError(t, err)

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"login": tt.user.Login,
			})
			tokenString, err := token.SignedString([]byte(tt.secretK))
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodGet, "/balance", bytes.NewBuffer([]byte("")))
			assert.NoError(t, err)
			req.AddCookie(&http.Cookie{Name: "Authorise", Value: tokenString, Expires: time.Now().Add(5 * time.Minute)})
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.status, rr.Code)
			if tt.status != http.StatusUnauthorized {
				body, err := io.ReadAll(rr.Body)
				assert.NoError(t, err)
				assert.Equal(t, body, reqBody)
			}

		})
	}
}
