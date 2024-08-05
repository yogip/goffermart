package handlers

// import (
// 	"bytes"
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"

// 	"goffermart/internal/core/model"
// 	"goffermart/internal/core/service"
// )

// func TestIAMHandler(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockService := service.NewMockIAM(ctrl)
// 	handler := NewIAMHandler(mockService)

// 	gin.SetMode(gin.TestMode)

// 	t.Run("Test Login with valid credentials", func(t *testing.T) {
// 		userReq := &model.UserRequest{Login: "test", Password: "test"}
// 		mockService.EXPECT().Login(gomock.Any(), userReq).Return("test_token", nil)

// 		router := gin.Default()
// 		router.POST("/login", handler.Login)

// 		body, _ := json.Marshal(userReq)
// 		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Header().Get("Authorization"), "Bearer test_token")
// 	})

// 	t.Run("Test Login with invalid credentials", func(t *testing.T) {
// 		userReq := &model.UserRequest{Login: "test", Password: "test"}
// 		mockService.EXPECT().Login(gomock.Any(), userReq).Return("", errors.New("test_error"))

// 		router := gin.Default()
// 		router.POST("/login", handler.Login)

// 		body, _ := json.Marshal(userReq)
// 		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusInternalServerError, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "Could not Auth user: test_error")
// 	})

// 	t.Run("Test Register with valid credentials", func(t *testing.T) {
// 		userReq := &model.UserRequest{Login: "test", Password: "test"}
// 		mockService.EXPECT().Register(gomock.Any(), userReq).Return("test_token", nil)

// 		router := gin.Default()
// 		router.POST("/register", handler.Register)

// 		body, _ := json.Marshal(userReq)
// 		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Header().Get("Authorization"), "Bearer test_token")
// 	})

// 	t.Run("Test Register with invalid credentials", func(t *testing.T) {
// 		userReq := &model.UserRequest{Login: "test", Password: "test"}
// 		mockService.EXPECT().Register(gomock.Any(), userReq).Return("", errors.New("test_error"))

// 		router := gin.Default()
// 		router.POST("/register", handler.Register)

// 		body, _ := json.Marshal(userReq)
// 		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusInternalServerError, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "User registration error: test_error")
// 	})
// }
