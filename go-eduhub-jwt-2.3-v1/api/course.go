package api

import (
	"a21hc3NpZ25tZW50/service"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourseAPI interface {
	DeleteCourse(c *gin.Context)
}

type courseAPI struct {
	courseService service.CourseService
}

func NewCourseAPI(courseService service.CourseService) *courseAPI {
	return &courseAPI{courseService}
}

func (cr *courseAPI) DeleteCourse(c *gin.Context) {
	courseId, err := strconv.Atoi(c.Param("course_id"))
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(400, gin.H{
			"message": "invalid course id",
		})
		return
	}

	err = cr.courseService.Delete(courseId)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "error internal server",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "course delete success",
	})
}
