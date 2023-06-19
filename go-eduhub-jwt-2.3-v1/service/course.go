package service

import (
	repo "a21hc3NpZ25tZW50/repository"
)

type CourseService interface {
	Delete(id int) error
}

type courseService struct {
	courseRepository repo.CourseRepository
}

func NewCourseService(courseRepository repo.CourseRepository) CourseService {
	return &courseService{courseRepository}
}

func (c *courseService) Delete(id int) error {
	err := c.courseRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
