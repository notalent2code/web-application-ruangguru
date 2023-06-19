package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Store(Category *model.Category) error
	Update(id int, category model.Category) error
	Delete(id int) error
	GetByID(id int) (*model.Category, error)
	GetList() ([]model.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) *categoryRepository {
	return &categoryRepository{db}
}

func (c *categoryRepository) Store(Category *model.Category) error {
	err := c.db.Create(Category).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *categoryRepository) Update(id int, category model.Category) error {// TODO: replace this
	err := c.db.Where(id).Updates(category).Error
	return err
}

func (c *categoryRepository) Delete(id int) error {	// TODO: replace this
	err := c.db.Where(id).Delete(&model.Category{})
	return err.Error 
}

func (c *categoryRepository) GetByID(id int) (*model.Category, error) {
	var Category model.Category
	err := c.db.Where("id = ?", id).First(&Category).Error
	if err != nil {
		return nil, err
	}

	return &Category, nil
}

func (c *categoryRepository) GetList() ([]model.Category, error) {// TODO: replace this
	var result []model.Category
	rows, err := c.db.Table("categories").Rows()
	if err != nil{
		return []model.Category{}, err
	}
	defer rows.Close()

	for rows.Next() { 
		c.db.ScanRows(rows, &result)
	}
	return result, nil 
}
