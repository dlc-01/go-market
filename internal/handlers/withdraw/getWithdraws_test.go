package withdraw

import (
	"bytes"
	"encoding/json"
	"github.com/dlc/go-market/internal/auth"
	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/model"
	"github.com/dlc/go-market/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetAllWithdraws(t *testing.T) {

	tests := []struct {
		name    string
		status  int
		secretK string

		user  *model.User
		userS *model.User
	}{
		{
			name:  "all good",
			user:  &model.User{Info: model.UserInfo{Login: "bob"}, Withdraws: []model.Withdraw{model.Withdraw{Order: "4029177534", Sum: 666}, {Order: "466455", Sum: 111}}},
			userS: &model.User{Info: model.UserInfo{Login: ""}, Withdraws: []model.Withdraw{model.Withdraw{Order: "4029177534", Sum: 666}, {Order: "466455", Sum: 111}}},

			secretK: "supersecret1234",
			status:  http.StatusOK,
		},
		{
			name:    "not true jwt token",
			user:    &model.User{Info: model.UserInfo{Login: "bob"}, Orders: []model.Order{{ID: ""}}},
			userS:   &model.User{Info: model.UserInfo{Login: ""}, Orders: []model.Order{}},
			secretK: "nosupersecret1234",
			status:  http.StatusUnauthorized,
		},
		{
			name:    "no content",
			user:    &model.User{Info: model.UserInfo{Login: "bob"}, Withdraws: []model.Withdraw{}},
			userS:   &model.User{Info: model.UserInfo{Login: ""}, Withdraws: []model.Withdraw{}},
			secretK: "supersecret1234",
			status:  http.StatusNoContent,
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
			mockUserService.On("GetAllWithdrawsByLogin", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*string")).Return(tt.userS, nil)

			rr := httptest.NewRecorder()
			storage.InitTestStorage(mockUserService)
			router := gin.Default()
			auth.SetSecretKey("supersecret1234")
			router.Use(auth.AuthMidlleware())
			router.GET("/withdraws", GetAllWithdraws)

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"login": tt.user.Info.Login,
			})
			tokenString, err := token.SignedString([]byte(tt.secretK))
			assert.NoError(t, err)
			resp, err := json.Marshal(tt.user.Withdraws)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodGet, "/withdraws", bytes.NewBuffer([]byte("")))
			assert.NoError(t, err)
			req.AddCookie(&http.Cookie{Name: "Authorise", Value: tokenString, Expires: time.Now().Add(5 * time.Minute)})
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.status, rr.Code)
			if rr.Code == http.StatusOK {
				body, err := io.ReadAll(rr.Body)
				assert.NoError(t, err)
				assert.Equal(t, resp, body)
			}

		})
	}
}
