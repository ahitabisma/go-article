package handler

import (
	"go-article/internal/handler/request"
	"go-article/internal/service"
	"go-article/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest

	// Validasi input JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		response := utils.APIResponse("Register account failed", http.StatusBadRequest, "error", nil, errors)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Panggil service untuk registrasi
	user, err := h.authService.Register(req)
	if err != nil {
		// Handle error spesifik duplikasi email
		if err.Error() == "email already registered" {
			response := utils.APIResponse("Register account failed", http.StatusConflict, "error", nil, gin.H{"email": "Email already registered"})
			c.JSON(http.StatusConflict, response)
			return
		}

		response := utils.APIResponse("Register account failed", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Berhasil registrasi
	response := utils.APIResponse("Account registered successfully", http.StatusCreated, "success", user, nil)
	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest

	// Validasi input JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		response := utils.APIResponse("Login failed", http.StatusBadRequest, "error", nil, errors)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Panggil service untuk login
	user, token, err := h.authService.Login(req)
	if err != nil {
		response := utils.APIResponse("Login failed", http.StatusUnauthorized, "error", nil, err.Error())
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	formatter := gin.H{
		"token": token,
		"user":  user,
	}

	response := utils.APIResponse("Successfuly logged in", http.StatusOK, "success", formatter, nil)
	c.JSON(http.StatusOK, response)
}
