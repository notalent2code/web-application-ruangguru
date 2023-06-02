package main

import (
	"strconv"
	"strings"
	"time"

	"encoding/base64"

	"github.com/gin-gonic/gin"
)

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var Posts = []Post{
	{ID: 1, Title: "Judul Postingan Pertama", Content: "Ini adalah postingan pertama di blog ini.", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, Title: "Judul Postingan Kedua", Content: "Ini adalah postingan kedua di blog ini.", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = []User{
	{Username: "user1", Password: "pass1"},
	{Username: "user2", Password: "pass2"},
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
			return
		}

		split := strings.Split(authHeader, " ")

		decoded, err := base64.StdEncoding.DecodeString(split[1])
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
			return
		}

		split = strings.Split(string(decoded), ":")
		username := split[0]
		password := split[1]

		for _, user := range users {
			if user.Username == username && user.Password == password {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
	}
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	//Set up authentication middleware here // TODO: replace this
	r.Use(authMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	r.GET("/posts", func(c *gin.Context) {
		id := c.Query("id")

		if id != "" {
			queryId, err := strconv.Atoi(id)
			if err != nil {
				c.JSON(400, gin.H{"error": "ID harus berupa angka"})
				return
			}

			for _, post := range Posts {
				if queryId == post.ID {
					c.JSON(200, gin.H{"post": post})
					return
				}
			}
			c.JSON(404, gin.H{"error": "Postingan tidak ditemukan"})
			return
		}

		c.JSON(200, gin.H{"posts": Posts})
	})

	r.POST("/posts", func(c *gin.Context) {
		var post Post

		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		post.ID = len(Posts) + 1
		post.CreatedAt = time.Now()
		post.UpdatedAt = time.Now()

		Posts = append(Posts, post)

		c.JSON(201, gin.H{"message": "Postingan berhasil ditambahkan", "post": post})
	})

	return r
}

func main() {
	r := SetupRouter()

	r.Run(":8080")
}
