package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kaa-dan/clean-architecture-go/internal/domain/entities"
	"github.com/kaa-dan/clean-architecture-go/internal/domain/services"
	"github.com/kaa-dan/clean-architecture-go/pkg/response"
	"github.com/kaa-dan/clean-architecture-go/pkg/validator"
)

type UserHandler struct {
	userService services.UserService
	validator   *validator.Validator
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator.New(),
	}
}

func (h *UserHandler) SignUp(c *gin.Context) {
	var req entities.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	result, err := h.userService.SignUp(c.Request.Context(), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, result)
}

func (h *UserHandler) SignIn(c *gin.Context) {
	var req entities.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	result, err := h.userService.SignIn(c.Request.Context(), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.Error(c, http.StatusBadRequest, "User ID is required")
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, user)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	users, err := h.userService.GetAllUsers(c.Request.Context(), limit, offset)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"users":  users,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.Error(c, http.StatusBadRequest, "User ID is required")
		return
	}

	// Check if user can update this profile (self or admin)
	currentUserID, _ := c.Get("user_id")
	currentUserRole, _ := c.Get("user_role")

	if currentUserID != userID && currentUserRole != string(entities.RoleAdmin) {
		response.Error(c, http.StatusForbidden, "You can only update your own profile")
		return
	}

	var req entities.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), userID, &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.Error(c, http.StatusBadRequest, "User ID is required")
		return
	}

	// Check if user can delete this profile (self or admin)
	currentUserID, _ := c.Get("user_id")
	currentUserRole, _ := c.Get("user_role")

	if currentUserID != userID && currentUserRole != string(entities.RoleAdmin) {
		response.Error(c, http.StatusForbidden, "You can only delete your own profile")
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), userID); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	user, err := h.userService.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, user)
}
