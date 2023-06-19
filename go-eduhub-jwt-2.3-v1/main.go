package main

import (
	"a21hc3NpZ25tZW50/api"
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/middleware"
	"a21hc3NpZ25tZW50/model"
	repo "a21hc3NpZ25tZW50/repository"
	"a21hc3NpZ25tZW50/service"
	"fmt"

	_ "embed"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type APIHandler struct {
	UserAPIHandler   api.UserAPI
	CourseAPIHandler api.CourseAPI
}

func main() {
	gin.SetMode(gin.ReleaseMode) //release

	router := gin.New()
	db := db.NewDB()

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

	conn.AutoMigrate(&model.User{}, &model.Course{})

	router = RunServer(conn, router)

	fmt.Println("Server is running on port 8080")
	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func RunServer(db *gorm.DB, gin *gin.Engine) *gin.Engine {
	userRepo := repo.NewUserRepo(db)
	courseRepo := repo.NewCourseRepo(db)

	userService := service.NewUserService(userRepo)
	courseService := service.NewCourseService(courseRepo)

	userAPIHandler := api.NewUserAPI(userService)
	courseAPIHandler := api.NewCourseAPI(courseService)

	apiHandler := APIHandler{
		UserAPIHandler:   userAPIHandler,
		CourseAPIHandler: courseAPIHandler,
	}

	users := gin.Group("/user")
	{
		users.POST("/login", apiHandler.UserAPIHandler.Login)
		users.POST("/register", apiHandler.UserAPIHandler.Register)
	}

	course := gin.Group("/course")
	{
		course.Use(middleware.Auth())
		course.DELETE("/delete/:course_id", apiHandler.CourseAPIHandler.DeleteCourse)
	}

	return gin
}
