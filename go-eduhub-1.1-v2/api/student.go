package api

import (
	"a21hc3NpZ25tZW50/model"
	repo "a21hc3NpZ25tZW50/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StudentAPI interface {
	AddStudent(c *gin.Context)
	GetStudents(c *gin.Context)
	GetStudentByID(c *gin.Context)
}

type studentAPI struct {
	studentRepo repo.StudentRepository
}

func NewStudentAPI(studentRepo repo.StudentRepository) *studentAPI {
	return &studentAPI{studentRepo}
}

func (s *studentAPI) AddStudent(c *gin.Context) {
	var student model.Student
	err := c.ShouldBindJSON(&student)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	err = s.studentRepo.Store(&student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "add student success"})
}

func (s *studentAPI) GetStudents(c *gin.Context) {
	students, err := s.studentRepo.FetchAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, students)
}

func (s *studentAPI) GetStudentByID(c *gin.Context) {
	studentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid student ID"})
		return
	}

	var students []model.Student
	students, err = s.studentRepo.FetchAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	var student model.Student
	var found bool = false

	for _, s := range students {
		if s.ID == studentId {
			student = s
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "student not found"})
		return
	}

	c.JSON(http.StatusOK, student)
}
