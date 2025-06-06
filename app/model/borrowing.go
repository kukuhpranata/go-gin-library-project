package model

import "time"

type Borrowing struct {
	Id         int
	UserId     int
	BookId     int
	BorrowDate time.Time
	DueDate    time.Time
	Status     string
	ReturnDate time.Time
}

type BorrowingJoin struct {
	Id         int
	UserId     string
	BookId     string
	BorrowDate time.Time
	DueDate    time.Time
	Status     string
	ReturnDate time.Time
}
