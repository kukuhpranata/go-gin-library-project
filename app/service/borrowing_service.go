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

type BorrowingService interface {
	Create(ctx context.Context, request *web.BorrowingCreateRequest) (*web.BorrowingResponse, *response.CustomError)
	Return(ctx context.Context, borrowingId int, userId int) (*web.BorrowingResponse, *response.CustomError)
	Find(ctx context.Context, borrowingId int) (*web.BorrowingResponse, *response.CustomError)
	FindAll(ctx context.Context) ([]web.BorrowingResponse, *response.CustomError)
}

type BorrowingServiceImpl struct {
	BorrowingRepository repository.BorrowingRepository
	BookRepository      repository.BookRepository
	DB                  *gorm.DB
	Validate            *validator.Validate
}

func NewBorrowingService(borrowingRepository repository.BorrowingRepository, bookRepository repository.BookRepository, DB *gorm.DB, validate *validator.Validate) BorrowingService {
	return &BorrowingServiceImpl{
		BorrowingRepository: borrowingRepository,
		BookRepository:      bookRepository,
		DB:                  DB,
		Validate:            validate,
	}
}

func (s *BorrowingServiceImpl) Create(ctx context.Context, request *web.BorrowingCreateRequest) (*web.BorrowingResponse, *response.CustomError) {
	err := s.Validate.Struct(request)
	if err != nil {
		return nil, response.BadRequestError(err.Error())
	}

	var book model.Book
	err = s.BookRepository.Find(s.DB, &book, request.BookId)
	if err != nil {
		return nil, response.RepositoryError(err.Error())
	}

	if book.Quantity < 1 {
		return nil, response.BadRequestError("Out of stock!")
	}

	var t time.Time
	dueDate, err := time.Parse("2006-01-02", request.DueDate)
	if err != nil {
		return nil, response.GeneralError(err.Error())
	}

	borrowing := model.Borrowing{
		BookId:     request.BookId,
		UserId:     request.UserId,
		BorrowDate: time.Now(),
		DueDate:    dueDate,
		Status:     "borrowed",
		ReturnDate: t,
	}

	err = s.BorrowingRepository.Save(s.DB, &borrowing)
	if err != nil {
		return nil, response.RepositoryError(err.Error())
	}

	book.Quantity--
	err = s.BookRepository.UpdateQuantity(s.DB, &book)
	if err != nil {
		return nil, response.RepositoryError(err.Error())
	}

	BorrowingResponse := web.BorrowingResponse{
		Id:         borrowing.Id,
		BookId:     request.BookId,
		UserId:     request.UserId,
		BorrowDate: borrowing.BorrowDate,
		DueDate:    borrowing.DueDate,
		Status:     borrowing.Status,
	}

	return &BorrowingResponse, nil
}

func (s *BorrowingServiceImpl) Return(ctx context.Context, borrowingId int, userId int) (*web.BorrowingResponse, *response.CustomError) {
	var borrowing model.Borrowing

	err := s.BorrowingRepository.Find(s.DB, &borrowing, borrowingId)
	if err != nil {
		return nil, response.RepositoryError(err.Error())
	}

	if userId != borrowing.UserId {
		return nil, response.BadRequestError("User not match!")
	}

	var book model.Book
	err = s.BookRepository.Find(s.DB, &book, borrowing.BookId)
	if err != nil {
		return nil, response.RepositoryError(err.Error())
	}

	status := "returned"
	if time.Now().After(borrowing.DueDate) {
		status = "late_returned"
	}

	borrowing.Status = status
	borrowing.ReturnDate = time.Now()

	err = s.BorrowingRepository.UpdateStatus(s.DB, &borrowing)
	if err != nil {
		return nil, response.RepositoryError(err.Error())
	}

	book.Quantity++
	err = s.BookRepository.UpdateQuantity(s.DB, &book)
	if err != nil {
		return nil, response.RepositoryError(err.Error())
	}

	BorrowingResponse := web.BorrowingResponse{
		Id:         borrowing.Id,
		BookId:     borrowing.BookId,
		UserId:     borrowing.UserId,
		BorrowDate: borrowing.BorrowDate,
		DueDate:    borrowing.DueDate,
		Status:     borrowing.Status,
		ReturnDate: borrowing.ReturnDate,
	}

	return &BorrowingResponse, nil
}

func (s *BorrowingServiceImpl) Find(ctx context.Context, borrowingId int) (*web.BorrowingResponse, *response.CustomError) {
	var borrowing model.Borrowing

	err := s.BorrowingRepository.Find(s.DB, &borrowing, borrowingId)
	if err != nil {
		return nil, response.NotFoundError(err.Error())
	}

	BorrowingResponse := web.BorrowingResponse{
		Id:         borrowing.Id,
		BookId:     borrowing.BookId,
		UserId:     borrowing.UserId,
		BorrowDate: borrowing.BorrowDate,
		DueDate:    borrowing.DueDate,
		Status:     borrowing.Status,
		ReturnDate: borrowing.ReturnDate,
	}

	return &BorrowingResponse, nil
}

func (s *BorrowingServiceImpl) FindAll(ctx context.Context) ([]web.BorrowingResponse, *response.CustomError) {
	var borrowings []model.Borrowing
	err := s.BorrowingRepository.FindAll(s.DB, &borrowings)
	if err != nil {
		return nil, response.RepositoryError(err.Error())
	}

	var borrowingResponses []web.BorrowingResponse
	for _, borrowing := range borrowings {
		borrowingResponse := web.BorrowingResponse{
			Id:         borrowing.Id,
			BookId:     borrowing.BookId,
			UserId:     borrowing.UserId,
			BorrowDate: borrowing.BorrowDate,
			DueDate:    borrowing.DueDate,
			Status:     borrowing.Status,
			ReturnDate: borrowing.ReturnDate,
		}
		borrowingResponses = append(borrowingResponses, borrowingResponse)
	}

	return borrowingResponses, nil
}
