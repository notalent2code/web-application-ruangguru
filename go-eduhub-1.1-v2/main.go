package main

import (
	"a21hc3NpZ25tZW50/api"
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/model"
	repo "a21hc3NpZ25tZW50/repository"
	"fmt"
	"time"

	_ "embed"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type APIHandler struct {
	StudentAPIHandler api.StudentAPI
	CourseAPIHandler  api.CourseAPI
}

func main() {
	gin.SetMode(gin.ReleaseMode) //release

	router := gin.New()
	db := db.NewDB()
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] \"%s %s %s\"\n",
			param.TimeStamp.Format(time.RFC822),
			param.Method,
			param.Path,
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())

	dbCredential := model.Credential{
		Host:         "localhost",
		Username:     "postgres",
		Password:     "kambing",
		DatabaseName: "be4_ruangguru",
		Port:         5432,
		Schema:       "public",
	}

	conn, err := db.Connect(&dbCredential)
	if err != nil {
		panic(err)
	}

	conn.AutoMigrate(&model.Student{}, &model.Course{})

	router = RunServer(conn, router)

	fmt.Println("Server is running on port 8080")
	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func RunServer(db *gorm.DB, gin *gin.Engine) *gin.Engine {
	studentRepo := repo.NewStudentRepo(db)
	courseRepo := repo.NewCourseRepo(db)

	studentAPIHandler := api.NewStudentAPI(studentRepo)
	courseAPIHandler := api.NewCourseAPI(courseRepo)

	apiHandler := APIHandler{
		StudentAPIHandler: studentAPIHandler,
		CourseAPIHandler:  courseAPIHandler,
	}

	student := gin.Group("/student")
	{
		student.POST("/add", apiHandler.StudentAPIHandler.AddStudent)
		student.GET("/gets", apiHandler.StudentAPIHandler.GetStudents)
		student.GET("/get/:id", apiHandler.StudentAPIHandler.GetStudentByID)
	}

	course := gin.Group("/course")
	{
		course.POST("/add", apiHandler.CourseAPIHandler.AddCourse)
	}

	return gin
}
