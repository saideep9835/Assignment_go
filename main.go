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
	CustomerName  string  `json:"customername"`
	CustomerEmail string  `json:"customeremail"` 
	
  }
  type Payment struct{
	gorm.Model
	Payments  int64 `json:"payments"`
	CustomerEmail string `json:"customeremail"`
  }
//   *gorm.DB
  func Init()  {
	// dbURL := "postgres://postgres:reddy123@localhost:5432/customer"
  	dbURL := "host=localhost user=postgres password=reddy123 dbname=customers port=5432 sslmode=disable"
	Database, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
  
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to DB")
	}
	Database.AutoMigrate(&User{})
	Database.AutoMigrate(&Payment{})
	DB=Database
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

// func SingleUserPayment(c *gin.Context) {
// 	var payments Payment

// 	if err := DB.Select("payments").Where("customer_email=?", c.Param("customer_email")).Find(&payments).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "record not found"})
// 	}
// 	c.JSON(http.StatusOK, gin.H{"data": payments})
// }
// func SumSingleUserPayment(c *gin.Context) {
// 	var userpayment Payment

// 	if err := DB.Select("sum(payments)").Where("customer_email=?", c.Param("customer_email")).Find(&userpayment).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "record not found"})
// 	}
// 	c.JSON(http.StatusOK, gin.H{"data": userpayment})
// }

func main() {
  

  Init()
  
  r := gin.Default()

  r.GET("/users", GetUsers)
  r.GET("/users/:customer_name", SingleUser)
  r.POST("/users", CreateUser)
//   r.GET("/payments/:customer_email", SingleUserPayment)
//   r.GET("/userpayment/:customer_email", SumSingleUserPayment)
  r.Run(":9000")
}