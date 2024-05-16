package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	Id          int    `json:"-" gorm:"column:id;"`
	Title       string `json:"title" gorm:"column:title;"`
	Description string `json:"description" gorm:"column:description;"`
	Status      string `json:"status" gorm:"column:status;"`
}

type Paging struct {
	Page  int   `json:"page" form:"page"`
	Limit int   `json:"limit" form:"limit"`
	Total int64 `json:"total" form:"total"`
}

func (TodoItem) TableName() string { return "todo_items" }

func (TodoCreationItem) TableName() string { return TodoItem{}.TableName() }

func (p *Paging) Process() {
	if p.Page < 0 {
		p.Page = 1
	}

	if p.Limit < 0 || p.Limit > 100 {
		p.Limit = 10
	}
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/dev_database?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	v1 := r.Group("/v1/api")
	{
		items := v1.Group("items")
		{
			items.POST("", createItem(db))
			items.GET("", getItems(db))
			items.GET("/:id", getItem(db))
			items.PUT("/:id")
			items.DELETE("/:id", deleteItem(db))
		}
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func createItem(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var data TodoCreationItem

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&data).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": data.Id})
	}
}

func getItems(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var result []TodoItem
		var paging Paging

		if err := db.Table(TodoItem{}.TableName()).Count(&paging.Total).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"data": err.Error()})
		}

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		paging.Process()

		if err := db.
			Offset((paging.Page - 1) * paging.Limit).
			Limit(paging.Limit).
			Find(&result).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"data": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"data": result, "paging": paging})
	}
}

func getItem(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var data TodoItem

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"data": err.Error()})
			return
		}

		if err := db.Where("id = ?", id).First(&data).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"data": "Item does not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": data})
	}
}

func deleteItem(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"data": err.Error()})
			return
		}

		if err := db.Table(TodoItem{}.TableName()).Where("id = ?", id).Delete(nil).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"data": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": true})
	}
}
