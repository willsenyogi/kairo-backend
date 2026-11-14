package handler

import (
	"kairo-backend/internal/dto"
	"kairo-backend/internal/enum"
	"kairo-backend/internal/model"
	"kairo-backend/internal/response"
	"kairo-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func InitAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: s,
	}
}
func (h *AuthHandler) GetAllUsers(c *gin.Context) {

	users, err := h.authService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{
			Status:  http.StatusInternalServerError,
			Method:  "GET",
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Status:         http.StatusOK,
		Method:         "GET",
		Success:        true,
		TotalDataIndex: len(users),
		Message:        "All Users Retrieved Successfully",
		Data:           users,
	})
}

func (h *AuthHandler) RegisterNewUser(c *gin.Context) {
	var req dto.RegisterRequest

	// Body validation
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Method:  "POST",
			Success: false,
			Message: err.Error(),
		})
		return
	}

	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     enum.RoleBuyer,
	}

	if err := h.authService.RegisterNewUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{
			Status:  http.StatusInternalServerError,
			Method:  "POST",
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.Response{
		Status:  http.StatusCreated,
		Method:  "POST",
		Success: true,
		Message: "User Account Successfully Created",
		Data:    user,
	})
}

func (h *AuthHandler) LoginValidateCreds(c *gin.Context) {
	var req dto.LoginRequest

	// Body validation
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Method:  "POST",
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if err := h.authService.LoginValidateCreds(&req); err != nil {
		c.JSON(http.StatusUnauthorized, response.Response{
			Status:  http.StatusUnauthorized,
			Method:  "POST",
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Status:  http.StatusOK,
		Method:  "POST",
		Success: true,
		Message: "Login validated",
	})
}
