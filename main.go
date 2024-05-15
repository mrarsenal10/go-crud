package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TodoItem struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type TodoCreationItem struct {
	Id          string `json:"-" gorm:"column:id"`
	Title       string `json:"title" binding:"required" gorm:"column:title"`
	Description string `json:"description" gorm:"column:description"`
	Status      string `json:"status" gorm:"column:status"`
}

func main() {
	// dsn := "user:pwd@tcp(localhost:3306)/database?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// if err != nil {
	// 	log.Fatal(err)
	// }

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		items := v1.Group("/items")
		{
			items.POST("", createItem())
			items.GET("")
			items.GET("/:id")
			items.PUT("/:id")
			items.DELETE("/:id")
		}
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func createItem() func(c *gin.Context) {
	return func(c *gin.Context) {
		var data TodoCreationItem

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": data})
	}
}
