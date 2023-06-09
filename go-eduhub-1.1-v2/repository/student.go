package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type StudentRepository interface {
	FetchAll() ([]model.Student, error)
	Store(student *model.Student) error
}

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepo(db *gorm.DB) *studentRepository {
	return &studentRepository{db}
}

func (s *studentRepository) FetchAll() ([]model.Student, error) {
	var students []model.Student
	err := s.db.Find(&students).Error
	if err != nil {
		return nil, err
	}

	return students, nil
}

func (s *studentRepository) Store(student *model.Student) error {
	err := s.db.Create(student).Error
	if err != nil {
		return err
	}

	return nil
}
