// Package user provides user management functionality for the cloud disk system.
// It supports registration, login, JWT authentication, and profile updates.
package user

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DB() *gorm.DB {
	return db
}

// AuthMiddleware is a middleware that checks for JWT in the Authorization header.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("username", claims["username"].(string))
		c.Set("user_id", uint64(claims["user_id"].(float64)))
		c.Next()
	}
}

// User 模型
type User struct {
	ID       uint64 `gorm:"primarykey"`
	Username string `gorm:"unique;size:255"`
	Password string `gorm:"size:255"`
	Avatar   string `gorm:"size:255"` //头像
	Bio      string `gorm:"size:512"` //个人简介
}

var db *gorm.DB
var jwtKey = []byte("secret")

func InitDB() {
	var err error
	db, err = gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/pan?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&User{})
}

// Register a new user with username and password
func Register(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.TrimSpace(user.Username) == "" || strings.TrimSpace(user.Password) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password cannot be empty"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	newUser := User{
		Username: user.Username,
		Password: string(hashedPassword),
	}

	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "User registered successfully",
		"user_id":  newUser.ID,
		"username": newUser.Username,
	})
}

// Login a user and return a JWT token
func Login(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var dbUser User
	if err := db.Where("username = ?", user.Username).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  dbUser.ID,
		"username": dbUser.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
	})
}

// GetUserInfo information by user ID
func GetUserInfo(c *gin.Context) {
	uid := c.MustGet("user_id").(uint64)
	var user User
	if err := db.First(&user, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// UpdateUserInfo updates a user's avatar and bio information.
func UpdateUserInfo(c *gin.Context) {
	uid := c.MustGet("user_id").(uint64)
	var userUpdates struct {
		Avator string `json:"avator"`
		Bio    string `json:"bio"`
	}
	if err := c.ShouldBindJSON(&userUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := db.First(&user, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if userUpdates.Avator != "" {
		db.Model(&user).Update("avator", userUpdates.Avator)
	}
	if userUpdates.Bio != "" {
		db.Model(&user).Update("bio", userUpdates.Bio)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User info updated successfully",
	})
}

// 启动gin服务器
func main() {
	InitDB()
	r := gin.Default()

	r.POST("/register", Register)
	r.POST("/login", Login)

	auth := r.Group("/user")
	auth.Use(AuthMiddleware())
	{
		auth.GET("/info", GetUserInfo)
		auth.PUT("/update", UpdateUserInfo)
	}

	r.Run(":8080")
}
