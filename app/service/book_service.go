package service

import (
	"context"
	"kukuh/go-gin-library-project/app/model"
	"kukuh/go-gin-library-project/app/repository"
	"kukuh/go-gin-library-project/app/web"
	"kukuh/go-gin-library-project/response"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type BookService interface {
	Create(ctx context.Context, request *web.BookCreate) (*web.BookResponse, *response.CustomError)
	Update(ctx context.Context, request *web.BookUpdate) (*web.BookResponse, *response.CustomError)
	Find(ctx context.Context, bookId int) (*web.BookResponse, *response.CustomError)
	Delete(ctx context.Context, bookId int) *response.CustomError
	FindAll(ctx context.Context) ([]web.BookResponse, *response.CustomError)
}

type BookServiceImpl struct {
	BookRepository repository.BookRepository
	DB             *gorm.DB
	Validate       *validator.Validate
}

func NewBookService(bookRepository repository.BookRepository, DB *gorm.DB, validate *validator.Validate) BookService {
	return &BookServiceImpl{
		BookRepository: bookRepository,
		DB:             DB,
		Validate:       validate,
	}
}

func (s *BookServiceImpl) Create(ctx context.Context, request *web.BookCreate) (*web.BookResponse, *response.CustomError) {
	err := s.Validate.Struct(request)
	if err != nil {
		return nil, response.BadRequestError(err.Error())
	}

	book := model.Book{
		Title:           request.Title,
		Author:          request.Author,
		Isbn:            request.Isbn,
		PublicationYear: request.PublicationYear,
		Quantity:        request.Quantity,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err = s.BookRepository.Save(s.DB, &book)
	if err != nil {
		return nil, response.RepositoryError(err.Error())
	}

	BookResponse := web.BookResponse{
		Id:              book.Id,
		Title:           book.Title,
		Author:          book.Author,
		Isbn:            book.Isbn,
		PublicationYear: book.PublicationYear,
		Quantity:        book.Quantity,
	}

	return &BookResponse, nil
}

func (s *BookServiceImpl) Update(ctx context.Context, request *web.BookUpdate) (*web.BookResponse, *response.CustomError) {
	err := s.Validate.Struct(request)
	if err != nil {
		return nil, response.BadRequestError(err.Error())
	}

	var book model.Book

	err = s.BookRepository.Find(s.DB, &book, request.Id)
	if err != nil {
		return nil, response.NotFoundError(err.Error())
	}

	book.Title = request.Title
	book.Author = request.Author
	book.Isbn = request.Isbn
	book.PublicationYear = request.PublicationYear
	book.Quantity = request.Quantity
	book.UpdatedAt = time.Now()

	err = s.BookRepository.Update(s.DB, &book)
	if err != nil {
		return nil, response.RepositoryError(err.Error())
	}

	BookResponse := web.BookResponse{
		Title:           book.Title,
		Author:          book.Author,
		Isbn:            book.Isbn,
		PublicationYear: book.PublicationYear,
		Quantity:        book.Quantity,
	}

	return &BookResponse, nil
}

func (s *BookServiceImpl) Find(ctx context.Context, bookId int) (*web.BookResponse, *response.CustomError) {
	var book model.Book

	err := s.BookRepository.Find(s.DB, &book, bookId)
	if err != nil {
		return nil, response.NotFoundError(err.Error())
	}

	BookResponse := web.BookResponse{
		Id:              bookId,
		Title:           book.Title,
		Author:          book.Author,
		Isbn:            book.Isbn,
		PublicationYear: book.PublicationYear,
		Quantity:        book.Quantity,
	}

	return &BookResponse, nil
}

func (s *BookServiceImpl) Delete(ctx context.Context, bookId int) *response.CustomError {
	var book model.Book

	err := s.BookRepository.Find(s.DB, &book, bookId)
	if err != nil {
		return response.NotFoundError(err.Error())
	}

	s.BookRepository.Delete(s.DB, book.Id)

	return nil
}

func (s *BookServiceImpl) FindAll(ctx context.Context) ([]web.BookResponse, *response.CustomError) {
	var books []model.Book
	err := s.BookRepository.FindAll(s.DB, &books)
	if err != nil {
		return nil, response.RepositoryError(err.Error())
	}

	var bookResponses []web.BookResponse
	for _, book := range books {
		bookResponse := web.BookResponse{
			Id:              book.Id,
			Title:           book.Title,
			Author:          book.Author,
			Isbn:            book.Isbn,
			PublicationYear: book.PublicationYear,
			Quantity:        book.Quantity,
		}
		bookResponses = append(bookResponses, bookResponse)
	}

	return bookResponses, nil
}
