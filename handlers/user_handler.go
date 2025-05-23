package handlers

import (
	"net/http"
	"strconv"

	"goapi/models"
	"goapi/repository/users_sql"
	"goapi/services"

	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.UserInput true "User object"
// @Success 201 {object} models.UserOutput
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.UserInput
	if err := c.ShouldBindJSON(&user); err != nil {
		errorResp := models.ToErrorResponse(models.NewAppError(http.StatusBadRequest, "invalid request body", err))
		c.JSON(errorResp.Code, errorResp)
		return
	}

	createdUser, err := h.userService.CreateUser(c.Request.Context(), &user)
	if err != nil {
		errorResp := models.ToErrorResponse(err)
		c.JSON(errorResp.Code, errorResp)
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

// GetUserByID godoc
// @Summary Get a user by ID
// @Description Get a user by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.UserOutput
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResp := models.ToErrorResponse(models.NewAppError(http.StatusBadRequest, "invalid user ID", err))
		c.JSON(errorResp.Code, errorResp)
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		errorResp := models.ToErrorResponse(err)
		c.JSON(errorResp.Code, errorResp)
		return
	}
	if user == nil {
		errorResp := models.ToErrorResponse(models.ErrNotFound)
		c.JSON(errorResp.Code, errorResp)
		return
	}

	c.JSON(http.StatusOK, user)
}

// ListUsers godoc
// @Summary List users
// @Description Get a paginated list of users
// @Tags users
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param name query string false "Name filter"
// @Param email query string false "Email filter"
// @Param order query string false "Sort order (ASC or DESC)"
// @Success 200 {object} models.UserList
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	params := users_sql.SearchParams{
		Limit:  10, // default limit
		Offset: 0,
	}

	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			params.Limit = l
		}
	}
	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			params.Offset = o
		}
	}

	params.Name = c.Query("name")
	params.Email = c.Query("email")
	params.Order = users_sql.SortOrder(c.Query("order"))

	users, err := h.userService.ListUsers(c.Request.Context(), params)
	if err != nil {
		errorResp := models.ToErrorResponse(err)
		c.JSON(errorResp.Code, errorResp)
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update an existing user's information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.UserInput true "User object"
// @Success 200 {object} models.UserOutput
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResp := models.ToErrorResponse(models.NewAppError(http.StatusBadRequest, "invalid user ID", err))
		c.JSON(errorResp.Code, errorResp)
		return
	}

	var user models.UserOutput
	if err := c.ShouldBindJSON(&user); err != nil {
		errorResp := models.ToErrorResponse(models.NewAppError(http.StatusBadRequest, "invalid request body", err))
		c.JSON(errorResp.Code, errorResp)
		return
	}

	user.ID = &id
	if err := h.userService.UpdateUser(c.Request.Context(), &user); err != nil {
		errorResp := models.ToErrorResponse(err)
		c.JSON(errorResp.Code, errorResp)
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResp := models.ToErrorResponse(models.NewAppError(http.StatusBadRequest, "invalid user ID", err))
		c.JSON(errorResp.Code, errorResp)
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		errorResp := models.ToErrorResponse(err)
		c.JSON(errorResp.Code, errorResp)
		return
	}

	c.Status(http.StatusNoContent)
}
