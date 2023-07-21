package orders

import (
	"bytes"
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

func TestAddNewStore(t *testing.T) {

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
			user:        &model.User{Info: model.UserInfo{Login: "bob"}, Orders: []model.Order{model.Order{ID: "4029177534"}}},
			userS:       &model.User{Info: model.UserInfo{Login: ""}, Orders: []model.Order{}},
			ContentType: "text/plain",
			secretK:     "supersecret1234",
			status:      http.StatusAccepted,
		},
		{
			name:        "not true jwt token",
			user:        &model.User{Info: model.UserInfo{Login: "bob"}, Orders: []model.Order{{ID: ""}}},
			userS:       &model.User{Info: model.UserInfo{Login: ""}, Orders: []model.Order{}},
			ContentType: "text/plain",
			secretK:     "nosupersecret1234",
			status:      http.StatusUnauthorized,
		},
		{
			name:        "not true content type",
			user:        &model.User{Info: model.UserInfo{Login: "bob"}, Orders: []model.Order{{ID: ""}}},
			userS:       &model.User{Info: model.UserInfo{Login: ""}, Orders: []model.Order{}},
			ContentType: "",
			secretK:     "supersecret1234",
			status:      http.StatusUnsupportedMediaType,
		},
		{
			name:        "not true id order",
			user:        &model.User{Info: model.UserInfo{Login: "bob"}, Orders: []model.Order{{ID: "23"}}},
			userS:       &model.User{Info: model.UserInfo{Login: ""}, Orders: []model.Order{}},
			ContentType: "text/plain",
			secretK:     "supersecret1234",
			status:      http.StatusUnprocessableEntity,
		},
		{
			name:        "id was in store with this log",
			user:        &model.User{Info: model.UserInfo{Login: "bob"}, Orders: []model.Order{model.Order{ID: "4029177534"}}},
			userS:       &model.User{Info: model.UserInfo{Login: "bob"}, Orders: []model.Order{model.Order{ID: "4029177534"}}},
			ContentType: "text/plain",
			secretK:     "supersecret1234",
			status:      http.StatusOK,
		},
		{
			name:        "id was in store with other log",
			user:        &model.User{Info: model.UserInfo{Login: "bob"}, Orders: []model.Order{model.Order{ID: "4029177534"}}},
			userS:       &model.User{Info: model.UserInfo{Login: "nobob"}, Orders: []model.Order{model.Order{ID: "4029177534"}}},
			ContentType: "text/plain",
			secretK:     "supersecret1234",
			status:      http.StatusConflict,
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
			mockUserService.On("AddNewOrder", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.User")).Return(tt.userS, nil)

			rr := httptest.NewRecorder()
			storage.InitTestStorage(mockUserService)
			router := gin.Default()
			auth.SetSecretKey("supersecret1234")
			router.Use(auth.AuthMidlleware())
			router.POST("/post", NewOrder)

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"login": tt.user.Info.Login,
			})
			tokenString, err := token.SignedString([]byte(tt.secretK))
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/post", bytes.NewBuffer([]byte(tt.user.Orders[0].ID)))
			req.Header.Set("Content-Type", tt.ContentType)
			assert.NoError(t, err)
			req.AddCookie(&http.Cookie{Name: "Authorise", Value: tokenString, Expires: time.Now().Add(5 * time.Minute)})
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.status, rr.Code)

		})
	}
}
