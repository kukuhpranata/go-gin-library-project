package web

import "time"

type BorrowingCreateRequest struct {
	BookId  int    `validate:"required" json:"book_id"`
	UserId  int    `json:"user_id"`
	DueDate string `validate:"required" json:"due_date"`
}

type BorrowingResponse struct {
	Id         int       `json:"id"`
	UserId     int       `json:"user_id"`
	BookId     int       `json:"book_id"`
	BorrowDate time.Time `json:"borrow_date"`
	DueDate    time.Time `json:"due_date"`
	Status     string    `json:"status"`
	ReturnDate time.Time `json:"return_date"`
}

type BorrowingFindResponse struct {
	Id         int       `json:"id"`
	User       string    `json:"user"`
	Book       string    `json:"book"`
	BorrowDate time.Time `json:"borrow_date"`
	DueDate    time.Time `json:"due_date"`
	Status     string    `json:"status"`
	ReturnDate time.Time `json:"return_date"`
}
