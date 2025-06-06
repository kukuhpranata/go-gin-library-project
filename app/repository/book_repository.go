package repository

import (
	"errors"
	"kukuh/go-gin-library-project/app/model"

	"gorm.io/gorm"
)

type BookRepository interface {
	Save(db *gorm.DB, book *model.Book) error
	Find(db *gorm.DB, book *model.Book, bookId int) error
	Update(db *gorm.DB, book *model.Book) error
	Delete(db *gorm.DB, bookId int) error
	FindAll(db *gorm.DB, books *[]model.Book) error
	UpdateQuantity(db *gorm.DB, book *model.Book) error
}

type BookRepositoryImpl struct {
}

func NewBookRepository() BookRepository {
	return &BookRepositoryImpl{}
}

func (r BookRepositoryImpl) Save(db *gorm.DB, book *model.Book) error {
	query := `INSERT INTO books (title, author, isbn, publication_year, quantity, created_at, updated_at) 
	VALUES (?,?,?,?,?,?,?)`
	result := db.Exec(query, book.Title, book.Author, book.Isbn, book.PublicationYear, book.Quantity, book.CreatedAt, book.UpdatedAt)

	if result.RowsAffected == 0 {
		return errors.New("failed to insert")
	}
	return nil
}

func (r BookRepositoryImpl) Find(db *gorm.DB, book *model.Book, bookId int) error {
	err := db.Raw("SELECT * FROM books WHERE id = ?", bookId).Scan(&book).Error
	if err != nil {
		return err
	}
	return nil
}

func (r BookRepositoryImpl) Update(db *gorm.DB, book *model.Book) error {
	result := db.Exec("UPDATE books set title = ?, author = ?, isbn = ?, publication_year = ?, quantity = ?, updated_at = ? WHERE id = ?", book.Title, book.Author, book.Isbn, book.PublicationYear, book.Quantity, book.UpdatedAt, book.Id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return errors.New("data not found")
	}
	return nil
}

func (r BookRepositoryImpl) Delete(db *gorm.DB, bookId int) error {
	result := db.Exec("DELETE FROM books where id = ?", bookId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return errors.New("data not found")
	}
	return nil
}

func (r BookRepositoryImpl) FindAll(db *gorm.DB, books *[]model.Book) error {
	result := db.Raw("SELECT * from books").Scan(&books)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return errors.New("data not found")
	}
	return nil
}

func (r BookRepositoryImpl) UpdateQuantity(db *gorm.DB, book *model.Book) error {
	result := db.Exec("UPDATE books set quantity = ? WHERE id = ?", book.Quantity, book.Id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return errors.New("data not found")
	}
	return nil
}
