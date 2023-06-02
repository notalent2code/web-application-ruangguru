package main

import (
	"a21hc3NpZ25tZW50/api"
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/model"
	repo "a21hc3NpZ25tZW50/repository"
	"a21hc3NpZ25tZW50/service"
	"fmt"
	"time"

	_ "embed"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type APIHandler struct {
	CategoryAPIHandler api.CategoryAPI
	TaskAPIHandler     api.TaskAPI
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

	conn.AutoMigrate(&model.Category{}, &model.Task{})

	router = RunServer(conn, router)

	fmt.Println("Server is running on port 8080")
	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func RunServer(db *gorm.DB, gin *gin.Engine) *gin.Engine {
	categoryRepo := repo.NewCategoryRepo(db)
	taskRepo := repo.NewTaskRepo(db)

	categoryService := service.NewCategoryService(categoryRepo)
	taskService := service.NewTaskService(taskRepo)

	categoryAPIHandler := api.NewCategoryAPI(categoryService)
	taskAPIHandler := api.NewTaskAPI(taskService)

	apiHandler := APIHandler{
		CategoryAPIHandler: categoryAPIHandler,
		TaskAPIHandler:     taskAPIHandler,
	}

	task := gin.Group("/task")
	{
		task.POST("/add", apiHandler.TaskAPIHandler.AddTask)
		task.GET("/get/:id", apiHandler.TaskAPIHandler.GetTaskByID)
		task.GET("/list", apiHandler.TaskAPIHandler.GetTaskList)
		task.PUT("/update/:id", apiHandler.TaskAPIHandler.UpdateTask)
		task.DELETE("/delete/:id", apiHandler.TaskAPIHandler.DeleteTask)
		task.GET("/category/:id", apiHandler.TaskAPIHandler.GetTaskListByCategory)

	}

	category := gin.Group("/category")
	{
		category.POST("/add", apiHandler.CategoryAPIHandler.AddCategory)
		category.GET("/get/:id", apiHandler.CategoryAPIHandler.GetCategoryByID)
		category.GET("/list", apiHandler.CategoryAPIHandler.GetCategoryList)
		category.PUT("/update/:id", apiHandler.CategoryAPIHandler.UpdateCategory)
		category.DELETE("/delete/:id", apiHandler.CategoryAPIHandler.DeleteCategory)
	}

	return gin
}
