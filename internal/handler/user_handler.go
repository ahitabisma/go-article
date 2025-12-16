package handler

import (
	"go-article/internal/service"
	"go-article/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Profile(c *gin.Context) {
	currentUser := c.MustGet("user_id").(float64)
	userID := uint64(currentUser)

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		response := utils.APIResponse("Failed to get user profile", http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("User profile", http.StatusOK, "success", user, nil)
	c.JSON(http.StatusOK, response)
}
