package repository

import (
	"errors"
	"kukuh/go-gin-library-project/app/model"

	"gorm.io/gorm"
)

type BorrowingRepository interface {
	Save(db *gorm.DB, borrowing *model.Borrowing) error
	Find(db *gorm.DB, borrowing *model.Borrowing, borrowingId int) error
	FindAll(db *gorm.DB, borrowings *[]model.Borrowing) error
	UpdateStatus(db *gorm.DB, borrowing *model.Borrowing) error
}

type BorrowingRepositoryImpl struct {
}

func NewBorrowingRepository() BorrowingRepository {
	return &BorrowingRepositoryImpl{}
}

func (r BorrowingRepositoryImpl) Save(db *gorm.DB, borrowing *model.Borrowing) error {
	query := `INSERT INTO borrowings (book_id, user_id, borrow_date, due_date, status) 
	VALUES (?,?,?,?,?)`
	result := db.Exec(query, borrowing.BookId, borrowing.UserId, borrowing.BorrowDate, borrowing.DueDate, borrowing.Status)

	if result.RowsAffected == 0 {
		return errors.New("failed to insert")
	}
	return nil
}

func (r BorrowingRepositoryImpl) Find(db *gorm.DB, borrowing *model.Borrowing, borrowingId int) error {
	err := db.Raw("SELECT * FROM borrowings WHERE id = ?", borrowingId).Scan(&borrowing).Error
	if err != nil {
		return err
	}
	return nil
}

func (r BorrowingRepositoryImpl) FindAll(db *gorm.DB, borrowings *[]model.Borrowing) error {
	result := db.Raw("SELECT * from borrowings").Scan(&borrowings)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return errors.New("data not found")
	}
	return nil
}

func (r BorrowingRepositoryImpl) UpdateStatus(db *gorm.DB, borrowing *model.Borrowing) error {
	result := db.Exec("UPDATE borrowings set status = ?, return_date = ? WHERE id = ?", borrowing.Status, borrowing.ReturnDate, borrowing.Id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return errors.New("data not found")
	}
	return nil
}
