package controller

import (
	"kukuh/go-gin-library-project/app/service"
	"kukuh/go-gin-library-project/app/web"
	"kukuh/go-gin-library-project/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BorrowingController interface {
	Create(ctx *gin.Context)
	Return(ctx *gin.Context)
	Find(ctx *gin.Context)
	FindAll(ctx *gin.Context)
}

type BorrowingControllerImpl struct {
	BorrowingService service.BorrowingService
}

func NewBorrowingController(borrowingService service.BorrowingService) BorrowingController {
	return &BorrowingControllerImpl{
		BorrowingService: borrowingService,
	}
}

func (c *BorrowingControllerImpl) Create(ctx *gin.Context) {
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

	borrowingCreateRequest := new(web.BorrowingCreateRequest)
	if err := ctx.ShouldBindJSON(borrowingCreateRequest); err != nil {
		customErr := response.BadRequestError("Invalid request body: " + err.Error())
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	userIdInt, _ := strconv.Atoi(userId)
	borrowingCreateRequest.UserId = userIdInt

	borrowingResponse, customErr := c.BorrowingService.Create(ctx.Request.Context(), borrowingCreateRequest)
	if customErr != nil {
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	webResponse := response.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   borrowingResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (c *BorrowingControllerImpl) Return(ctx *gin.Context) {
	id := ctx.Param("id")

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
	idInt, _ := strconv.Atoi(id)
	userIdInt, _ := strconv.Atoi(userId)

	borrowingReturnResponse, customErr := c.BorrowingService.Return(ctx.Request.Context(), idInt, userIdInt)
	if customErr != nil {
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	webResponse := response.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   borrowingReturnResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (c *BorrowingControllerImpl) Find(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)

	borrowingResponse, customErr := c.BorrowingService.Find(ctx.Request.Context(), idInt)
	if customErr != nil {
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	webResponse := response.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   borrowingResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (c *BorrowingControllerImpl) FindAll(ctx *gin.Context) {
	borrowingResponses, customErr := c.BorrowingService.FindAll(ctx.Request.Context())
	if customErr != nil {
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	webResponse := response.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   borrowingResponses,
	}

	ctx.JSON(http.StatusOK, webResponse)
}
