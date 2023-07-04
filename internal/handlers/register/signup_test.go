package register

import (
	"bytes"
	"encoding/json"
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

func TestHConfig_SingUP2(t *testing.T) {

	tests := []struct {
		name   string
		status int
		realP  string
		userR  *model.UserInfo
		userS  *model.UserInfo
	}{
		{
			name: "Conflict",
			userR: &model.UserInfo{
				Password: "supersecret1234",
				Login:    "bob",
			},
			userS: &model.UserInfo{
				Password: "supersecret1234",
				Login:    "bob",
			},
			status: http.StatusConflict,
		},
		{
			name:   "invalid username",
			status: http.StatusBadRequest,
			userR: &model.UserInfo{
				Password: "supersecret1234",
			},
			userS: &model.UserInfo{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gin.SetMode(gin.TestMode)
			if err := logger.InitLogger(); err != nil {
				log.Fatal(err)
			}

			mockUserService := new(storage.TestStore)
			mockUserService.On("CreateUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.UserInfo")).Return(tt.userS, nil)
			rr := httptest.NewRecorder()

			storage.InitTestStorage(mockUserService)

			router := gin.Default()

			router.POST("/signup", Register)
			reqBody, err := json.Marshal(gin.H{
				"password": tt.userR.Password,
				"login":    tt.userR.Login,
			})
			assert.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(rr, request)

			assert.Equal(t, tt.status, rr.Code)
			mockUserService.AssertNotCalled(t, "createUser")

		})
	}
}
