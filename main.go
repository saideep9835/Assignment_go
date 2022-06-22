package main

import (
	// "encoding/json"
	// "log"
	"fmt"
	"net/http"
	//   "github.com/gorilla/mux"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	gorm.Model
	CustomerName  string `json:"customername"`
	CustomerEmail string `json:"customeremail"`
}

type UpdateUserInput struct {
	CustomerName  string `json:"customername"`
	CustomerEmail string `json:"customeremail"`
}

func DatabaseCreation() {
	// dbURL := "postgres://postgres:reddy123@localhost:5432/customer"
	dbURL := "host=localhost user=postgres password=reddy123 dbname=customers port=5432 sslmode=disable"
	Database, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to DB")
	}
	Database.AutoMigrate(&User{})

	DB = Database
	// return DB
}

// func GetUsers(c *gin.Context) {
// 	c.JSON
// 	var users []User
// 	DB.Find(&users)
// 	json.NewEncoder(w).Encode(users)
// }
func GetUsers(c *gin.Context) {
	var users []User

	DB.Limit(10).Find(&users)
	c.JSON(http.StatusOK, gin.H{"data": users})
}
func SingleUser(c *gin.Context) {
	var users User

	if err := DB.Where("customer_name=?", c.Param("customer_name")).First(&users).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record not found"})
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}
func CreateUser(c *gin.Context) {
	// Validate input
	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create book
	user_details := User{CustomerName: input.CustomerName, CustomerEmail: input.CustomerEmail}
	DB.Create(&user_details)

	c.JSON(http.StatusOK, gin.H{"data": user_details})
}

func DeleteUser(c *gin.Context) {
	// Get model if exist
	var deluser User
	if err := DB.Where("customer_name = ?", c.Param("customer_name")).First(&deluser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	DB.Delete(&deluser)

	c.JSON(http.StatusOK, gin.H{"data": true})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var body UpdateUserInput
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var updateuser User
	if res := DB.First(&updateuser, id); res.Error != nil {
		c.AbortWithError(http.StatusNotFound, res.Error)
		return
	}
	updateuser.CustomerName = body.CustomerName

	updateuser.CustomerEmail = body.CustomerEmail

	DB.Save(&updateuser)
	c.JSON(http.StatusOK, &updateuser)
}
func main() {

	DatabaseCreation()

	r := gin.Default()

	r.GET("/users", GetUsers)
	r.GET("/users/:customer_name", SingleUser)
	r.POST("/users", CreateUser)
	r.DELETE("/delusers/:customer_name", DeleteUser)
	r.POST("/updater/:id", UpdateUser)
	r.Run(":9000")
}
