package withdraw

import (
	"bytes"
	"encoding/json"
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

func TestAddWithdraw(t *testing.T) {

	tests := []struct {
		name        string
		status      int
		secretK     string
		ContentType string
		user        *model.User
		userS       *model.User
	}{
		{
			name:        "all good",
			user:        &model.User{Info: model.UserInfo{Login: "bob"}, Withdraws: []model.Withdraw{model.Withdraw{Order: "4029177534", Sum: 1}}},
			userS:       &model.User{Info: model.UserInfo{Login: ""}, Withdraws: []model.Withdraw{}, Balance: 10},
			ContentType: "application/json",
			secretK:     "supersecret1234",
			status:      http.StatusAccepted,
		},
		{
			name:        "not true jwt token",
			user:        &model.User{Info: model.UserInfo{Login: "bob"}, Withdraws: []model.Withdraw{model.Withdraw{Order: "4029177534", Sum: 1}}},
			userS:       &model.User{Info: model.UserInfo{Login: ""}, Withdraws: []model.Withdraw{}, Balance: 10},
			ContentType: "application/json",
			secretK:     "nosupersecret1234",
			status:      http.StatusUnauthorized,
		},
		{
			name:        "not true content type",
			user:        &model.User{Info: model.UserInfo{Login: "bob"}, Withdraws: []model.Withdraw{model.Withdraw{Order: "4029177534", Sum: 1}}},
			userS:       &model.User{Info: model.UserInfo{Login: ""}, Withdraws: []model.Withdraw{}, Balance: 10},
			ContentType: "",
			secretK:     "supersecret1234",
			status:      http.StatusUnsupportedMediaType,
		},
		{
			name:        "not true id order",
			user:        &model.User{Info: model.UserInfo{Login: "bob"}, Withdraws: []model.Withdraw{model.Withdraw{Order: "12", Sum: 1}}},
			userS:       &model.User{Info: model.UserInfo{Login: ""}, Withdraws: []model.Withdraw{}, Balance: 10},
			ContentType: "application/json",
			secretK:     "supersecret1234",
			status:      http.StatusUnprocessableEntity,
		},
		{
			name:        "sum > balnce",
			user:        &model.User{Info: model.UserInfo{Login: "bob"}, Withdraws: []model.Withdraw{model.Withdraw{Order: "4029177534", Sum: 10}}},
			userS:       &model.User{Info: model.UserInfo{Login: ""}, Withdraws: []model.Withdraw{}, Balance: 0},
			ContentType: "application/json",
			secretK:     "supersecret1234",
			status:      http.StatusPaymentRequired,
		},
		{
			name:        "sum < 0 ",
			user:        &model.User{Info: model.UserInfo{Login: "bob"}, Withdraws: []model.Withdraw{model.Withdraw{Order: "4029177534", Sum: -10}}},
			userS:       &model.User{Info: model.UserInfo{Login: ""}, Withdraws: []model.Withdraw{}, Balance: 0},
			ContentType: "application/json",
			secretK:     "supersecret1234",
			status:      http.StatusUnprocessableEntity,
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
			mockUserService.On("AddNewOderWithdraw", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.User")).Return(tt.userS, nil)

			rr := httptest.NewRecorder()
			storage.InitTestStorage(mockUserService)
			router := gin.Default()
			auth.SetSecretKey("supersecret1234")
			router.Use(auth.AuthMidlleware())
			router.POST("/post", WithdrawalOfFunds)

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"login": tt.user.Info.Login,
			})
			tokenString, err := token.SignedString([]byte(tt.secretK))
			assert.NoError(t, err)
			jsons, err := json.Marshal(tt.user.Withdraws[0])
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/post", bytes.NewBuffer((jsons)))
			req.Header.Set("Content-Type", tt.ContentType)
			assert.NoError(t, err)
			req.AddCookie(&http.Cookie{Name: "Authorise", Value: tokenString, Expires: time.Now().Add(5 * time.Minute)})
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.status, rr.Code)

		})
	}
}
