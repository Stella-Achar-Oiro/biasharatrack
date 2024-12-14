package controllers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/OAthooh/BiasharaTrack.git/models"
	"github.com/OAthooh/BiasharaTrack.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	Db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{Db: db}
}

func (auth *AuthHandler) Login(c *gin.Context) {
	// Handle preflight OPTIONS request
	if c.Request.Method == "OPTIONS" {
		c.Header("Access-Control-Allow-Methods", "POST")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Status(http.StatusOK)
		return
	}

	// Parse JSON request body
	var loginRequest models.AuthRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		utils.WarningLogger("Invalid login request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get user from database
	var user models.User
	result := auth.Db.Where("email = ?", loginRequest.Email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			utils.WarningLogger("Login attempt with non-existent email: %s", loginRequest.Email)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		utils.ErrorLogger("Database error during login: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check password
	if !auth.checkPassword(loginRequest.Password, user.Password) {
		utils.WarningLogger("Failed login attempt for email: %s", loginRequest.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"email":     user.Email,
		"full_name": user.FullName,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		utils.ErrorLogger("Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	utils.InfoLogger("Successful login for user: %s", loginRequest.Email)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
		"user": gin.H{
			"id":            user.ID,
			"full_name":     user.FullName,
			"email":         user.Email,
			"business_name": user.BusinessName,
			"telephone":     user.Telephone,
			"location":      user.Location,
		},
	})
}

func (auth *AuthHandler) Register(c *gin.Context) {
	// Handle preflight OPTIONS request
	if c.Request.Method == "OPTIONS" {
		c.Header("Access-Control-Allow-Methods", "POST")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Status(http.StatusOK)
		return
	}

	// Parse JSON request body
	var registerRequest models.Register

	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		utils.WarningLogger("Invalid registration request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Ensure all registration data are present
	if registerRequest.FullName == "" || registerRequest.Email == "" || registerRequest.Password == "" ||
		registerRequest.BusinessName == "" || registerRequest.Telephone == "" || registerRequest.Location == "" {
		utils.WarningLogger("Incomplete registration data: %+v", registerRequest)
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	// Check if email already exists
	var existingUser models.User
	result := auth.Db.Where("email = ?", registerRequest.Email).First(&existingUser)
	if result.Error == nil {
		utils.WarningLogger("Registration attempt with existing email: %s", registerRequest.Email)
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		utils.ErrorLogger("Database error checking email existence: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorLogger("Error hashing password during registration: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Create new user
	newUser := models.User{
		FullName:     registerRequest.FullName,
		Email:        registerRequest.Email,
		Password:     string(hashedPassword),
		BusinessName: registerRequest.BusinessName,
		Telephone:    registerRequest.Telephone,
		Location:     registerRequest.Location,
	}

	result = auth.Db.Create(&newUser)
	if result.Error != nil {
		utils.ErrorLogger("Error creating new user: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	// Generate JWT token for the new user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":       newUser.ID,
		"email":         newUser.Email,
		"full_name":     newUser.FullName,
		"business_name": newUser.BusinessName,
		"telephone":     newUser.Telephone,
		"location":      newUser.Location,
		"exp":           time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		utils.ErrorLogger("Error generating token for new user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	utils.InfoLogger("Successfully registered new user: %s", newUser.Email)
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   tokenString,
		"user": gin.H{
			"id":            newUser.ID,
			"full_name":     newUser.FullName,
			"email":         newUser.Email,
			"business_name": newUser.BusinessName,
			"telephone":     newUser.Telephone,
			"location":      newUser.Location,
		},
	})
}

func (auth *AuthHandler) VerifyToken(c *gin.Context) {
	// Handle preflight OPTIONS request
	if c.Request.Method == "OPTIONS" {
		c.Header("Access-Control-Allow-Methods", "POST")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Status(http.StatusOK)
		return
	}

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	// Remove 'Bearer ' prefix if present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("f50559429275498b09d13392269fc0fd02a2f548d8470c3765a8895212080636"), nil // Replace with secure secret key
	})

	if err != nil || !token.Valid {
		utils.WarningLogger("Invalid token verification attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":        claims["user_id"],
			"full_name": claims["full_name"],
			"email":     claims["email"],
		},
	})
}

func (auth *AuthHandler) checkPassword(providedPassword, storedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(providedPassword))
	return err == nil
}
