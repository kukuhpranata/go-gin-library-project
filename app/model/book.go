package model

import "time"

type Book struct {
	Id              int
	Title           string
	Author          string
	Isbn            string
	PublicationYear int
	Quantity        int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
