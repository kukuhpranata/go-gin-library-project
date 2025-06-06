package controller

import (
	"kukuh/go-gin-library-project/app/service"
	"kukuh/go-gin-library-project/app/web"
	"kukuh/go-gin-library-project/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookController interface {
	Create(ctx *gin.Context)
	Find(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindAll(ctx *gin.Context)
}

type BookControllerImpl struct {
	BookService service.BookService
}

func NewBookController(bookService service.BookService) BookController {
	return &BookControllerImpl{
		BookService: bookService,
	}
}

func (c *BookControllerImpl) Create(ctx *gin.Context) {
	bookCreateRequest := new(web.BookCreate)
	if err := ctx.ShouldBindJSON(bookCreateRequest); err != nil {
		customErr := response.BadRequestError("Invalid request body: " + err.Error())
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	bookResponse, customErr := c.BookService.Create(ctx.Request.Context(), bookCreateRequest)
	if customErr != nil {
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	webResponse := response.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   bookResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (c *BookControllerImpl) Find(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)

	bookResponse, customErr := c.BookService.Find(ctx.Request.Context(), idInt)
	if customErr != nil {
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	webResponse := response.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   bookResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (c *BookControllerImpl) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	bookUpdateRequest := new(web.BookUpdate)
	if err := ctx.ShouldBindJSON(bookUpdateRequest); err != nil {
		customErr := response.BadRequestError("Invalid request body: " + err.Error())
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	idInt, _ := strconv.Atoi(id)
	bookUpdateRequest.Id = idInt

	bookResponse, customErr := c.BookService.Update(ctx.Request.Context(), bookUpdateRequest)
	if customErr != nil {
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	webResponse := response.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   bookResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (c *BookControllerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)

	customErr := c.BookService.Delete(ctx.Request.Context(), idInt)
	if customErr != nil {
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	webResponse := response.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   gin.H{"message": "Book deleted successfully"},
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (c *BookControllerImpl) FindAll(ctx *gin.Context) {
	bookResponses, customErr := c.BookService.FindAll(ctx.Request.Context())
	if customErr != nil {
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	webResponse := response.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   bookResponses,
	}

	ctx.JSON(http.StatusOK, webResponse)
}
