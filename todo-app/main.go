package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type Todo struct {
    ID     uint   `json:"id" gorm:"primaryKey"`
    Task   string `json:"task"`
    Status string `json:"status"`
}

var db *gorm.DB

func initDB() {
    var err error
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    db.AutoMigrate(&Todo{})
}

func main() {
    initDB()

    r := gin.Default()

    r.POST("/todos", createTodo)
    r.GET("/todos", getTodos)
    r.DELETE("/todos/:id", deleteTodo)

    r.Run(":8080")
}

func createTodo(c *gin.Context) {
    var todo Todo
    if err := c.ShouldBindJSON(&todo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    db.Create(&todo)
    c.JSON(http.StatusOK, todo)
}

func getTodos(c *gin.Context) {
    var todos []Todo
    db.Find(&todos)
    c.JSON(http.StatusOK, todos)
}

func deleteTodo(c *gin.Context) {
    id := c.Param("id")
    db.Delete(&Todo{}, id)
    c.Status(http.StatusNoContent)
}
