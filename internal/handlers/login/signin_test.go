package login

import (
	"bytes"
	"encoding/json"
	"github.com/dlc/go-market/internal/hash"
	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/model"
	"github.com/dlc/go-market/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHConfig_Login(t *testing.T) {

	tests := []struct {
		name   string
		status int
		realP  string
		user   *model.User
	}{
		{
			name: "signInTrue",
			user: &model.User{Info: model.UserInfo{
				Password: "supersecret1234",
				Login:    "bob"},
			},
			realP:  "supersecret1234",
			status: http.StatusOK,
		},
		{
			name: "signInFalse",
			user: &model.User{Info: model.UserInfo{
				Password: "supersecret1234",
				Login:    "bob"},
			},
			realP:  "notSupersecret1234",
			status: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gin.SetMode(gin.TestMode)
			if err := logger.InitLogger(); err != nil {
				log.Fatal(err)
			}
			pas, _ := hash.HashPassword(tt.realP)
			stor := &model.UserInfo{Login: tt.user.Info.Login, Password: pas}
			mockUserService := new(storage.TestStore)
			mockUserService.On("findByLogin", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*string")).Return(stor, nil)
			rr := httptest.NewRecorder()
			storage.InitTestStorage(mockUserService)
			router := gin.Default()

			router.POST("/signin", Login)

			reqBody, err := json.Marshal(gin.H{
				"password": tt.user.Info.Password,
				"login":    tt.user.Info.Login,
			})
			assert.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(rr, request)

			assert.Equal(t, tt.status, rr.Code)
			mockUserService.AssertNotCalled(t, "FindByEmail")

		})
	}
}
