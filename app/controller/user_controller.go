package controller

import (
	"kukuh/go-gin-library-project/app/service"
	"kukuh/go-gin-library-project/app/web"
	"kukuh/go-gin-library-project/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	UpdateUserOwn(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
}

type UserControllerImpl struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerImpl{
		UserService: userService,
	}
}

func (c *UserControllerImpl) Register(ctx *gin.Context) {
	userCreateRequest := new(web.Register)
	if err := ctx.ShouldBindJSON(userCreateRequest); err != nil {
		customErr := response.BadRequestError("Invalid request body")
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	userResponse, customErr := c.UserService.Register(ctx.Request.Context(), userCreateRequest)
	if customErr != nil {
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	webResponse := response.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   userResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (c *UserControllerImpl) Login(ctx *gin.Context) {
	loginRequest := new(web.LoginUserRequest)
	if err := ctx.ShouldBindJSON(loginRequest); err != nil {
		customErr := response.BadRequestError("Invalid request body")
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	loginResponse, customErr := c.UserService.Login(ctx.Request.Context(), loginRequest)
	if customErr != nil {
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}
	webResponse := response.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   loginResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (c *UserControllerImpl) UpdateUserOwn(ctx *gin.Context) {
	authId, exists := ctx.Get("authId")
	if !exists {
		customErr := response.UnauthorizedError("Authentication ID not found")
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}
	userId, ok := authId.(string)
	if !ok {
		customErr := response.UnauthorizedError("Invalid authentication ID format")
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	userUpdateRequest := new(web.UpdateUserRequest)
	if err := ctx.ShouldBindJSON(userUpdateRequest); err != nil {
		customErr := response.BadRequestError("Invalid request body")
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}
	userIdint, _ := strconv.Atoi(userId)
	userUpdateRequest.Id = userIdint

	userResponse, customErr := c.UserService.UpdateUserOwn(ctx.Request.Context(), userUpdateRequest)
	if customErr != nil {
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}
	webResponse := response.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   userResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (c *UserControllerImpl) DeleteUser(ctx *gin.Context) {
	authId, exists := ctx.Get("authId")
	if !exists {
		customErr := response.UnauthorizedError("Authentication ID not found")
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}
	userId, ok := authId.(string)
	if !ok {
		customErr := response.UnauthorizedError("Invalid authentication ID format")
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	userIdint, _ := strconv.Atoi(userId)

	customErr := c.UserService.Delete(ctx.Request.Context(), userIdint)
	if customErr != nil {
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	webResponse := response.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   gin.H{"message": "User deleted successfully"},
	}

	ctx.JSON(http.StatusOK, webResponse)
}
