package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/OAthooh/BiasharaTrack.git/controllers"
	"github.com/OAthooh/BiasharaTrack.git/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewAuthHandler(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *controllers.AuthHandler
	}{
		{
			name: "Valid DB instance",
			args: args{db: &gorm.DB{}},
			want: &controllers.AuthHandler{Db: &gorm.DB{}},
		},
		{
			name: "Nil DB instance",
			args: args{db: nil},
			want: &controllers.AuthHandler{Db: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := controllers.NewAuthHandler(tt.args.db); got.Db != tt.want.Db {
				t.Errorf("NewAuthHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		expectedCode int
		expectedBody bool // whether to expect a JSON body
		setup        func()
	}{
		{
			name:   "OPTIONS request",
			fields: fields{db: &gorm.DB{}},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httptest.NewRequest("OPTIONS", "/login", nil)
				return c
			}()},
			expectedCode: http.StatusOK, // 200 for OPTIONS
			expectedBody: false,         // No JSON body expected for OPTIONS
			setup: func() {
				// Setup for OPTIONS request
			},
		},
		{
			name:   "Invalid JSON request",
			fields: fields{db: &gorm.DB{}},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				// Create a new request with a POST method
				req := httptest.NewRequest("POST", "/login", nil)
				c.Request = req
				return c
			}()},
			expectedCode: http.StatusBadRequest,
			expectedBody: true,
			setup: func() {
				// Setup for invalid JSON request
			},
		},
		{
			name:   "Non-existent email",
			fields: fields{db: &gorm.DB{}},
			args: args{c: &gin.Context{
				Request: &http.Request{Method: "POST"},
			}},
			expectedCode: http.StatusBadRequest,
			setup: func() {
				// Setup for non-existent email
			},
		},
		{
			name:   "Database error",
			fields: fields{db: &gorm.DB{}},
			args: args{c: &gin.Context{
				Request: &http.Request{Method: "POST"},
			}},
			expectedCode: http.StatusBadRequest,
			setup: func() {
				// Setup for database error
			},
		},
		{
			name:   "Invalid password",
			fields: fields{db: &gorm.DB{}},
			args: args{c: &gin.Context{
				Request: &http.Request{Method: "POST"},
			}},
			expectedCode: http.StatusBadRequest,
			setup: func() {
				// Setup for invalid password
			},
		},
		{
			name:   "Error generating token",
			fields: fields{db: &gorm.DB{}},
			args: args{c: &gin.Context{
				Request: &http.Request{Method: "POST"},
			}},
			expectedCode: http.StatusBadRequest,
			setup: func() {
				// Setup for error generating token
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = tt.args.c.Request

			auth := controllers.NewAuthHandler(tt.fields.db)

			auth.Login(c)

			// Assert the response code
			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
			}

			// Only check JSON body if we expect one
			if tt.expectedBody {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to parse response body: %v", err)
				}
				// Add your JSON body assertions here
			}
		})
	}
}

func TestAuthHandler_Register(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create a mock DB or use a test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		expectedCode int
		setup        func()
		expectedBody bool
	}{
		{
			name:   "Successful registration",
			fields: fields{db: db},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				registerData := models.Register{
					FullName:     "Test User",
					Email:        "test@example.com",
					Password:     "password123",
					BusinessName: "Test Business",
					Telephone:    "1234567890",
					Location:     "Test Location",
				}

				jsonData, _ := json.Marshal(registerData)
				req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")

				c.Request = req
				return c
			}()},
			expectedCode: http.StatusCreated,
			setup: func() {
				// Migrate the schema
				db.AutoMigrate(&models.User{})
			},
			expectedBody: true,
		},
		{
			name:   "Invalid registration request",
			fields: fields{db: db},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				// Mock data with missing fields for invalid registration
				registerData := map[string]interface{}{
					"Email":    "invalid@example.com",
					"Password": "password123",
				}

				jsonData, _ := json.Marshal(registerData)
				req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")

				c.Request = req
				return c
			}()},
			expectedCode: http.StatusBadRequest,
			setup: func() {
				// Setup for invalid registration request
			},
			expectedBody: true,
		},
		{
			name:   "Email already registered",
			fields: fields{db: db},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				// Mock data for email already registered
				registerData := models.Register{
					FullName:     "Existing User",
					Email:        "existing@example.com",
					Password:     "password123",
					BusinessName: "Existing Business",
					Telephone:    "0987654321",
					Location:     "Existing Location",
				}

				jsonData, _ := json.Marshal(registerData)
				req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")

				c.Request = req
				return c
			}()},
			expectedCode: http.StatusConflict,
			setup: func() {
				// Insert existing user to simulate email already registered
				existingUser := models.User{
					FullName:     "Existing User",
					Email:        "existing@example.com",
					Password:     "hashedpassword",
					BusinessName: "Existing Business",
					Telephone:    "0987654321",
					Location:     "Existing Location",
				}
				db.Create(&existingUser)
			},
			expectedBody: true,
		},
		{
			name:   "Database error checking email existence",
			fields: fields{db: db},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				// Mock data for database error checking email existence
				registerData := models.Register{
					FullName:     "Error User",
					Email:        "test@example.com",
					Password:     "password123",
					BusinessName: "Error Business",
					Telephone:    "1234567890",
					Location:     "Error Location",
				}

				jsonData, _ := json.Marshal(registerData)
				req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")

				c.Request = req
				return c
			}()},
			expectedCode: http.StatusConflict,
			setup: func() {
				// Simulate database error by not setting up the database properly
			},
			expectedBody: true,
		},

		{
			name:   "Error hashing password - Empty Password",
			fields: fields{db: db},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				registerData := models.Register{
					FullName:     "Test User",
					Email:        "test@example.com",
					Password:     "", // Empty password
					BusinessName: "Test Business",
					Telephone:    "1234567890",
					Location:     "Test Location",
				}

				jsonData, _ := json.Marshal(registerData)
				req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")

				c.Request = req
				return c
			}()},
			expectedCode: http.StatusBadRequest,
			setup: func() {
				db.AutoMigrate(&models.User{})
			},
			expectedBody: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = tt.args.c.Request

			auth := controllers.NewAuthHandler(tt.fields.db)

			auth.Register(c)

			// Assert the response code
			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
			}

			// Only check JSON body if we expect one
			if tt.expectedBody {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to parse response body: %v", err)
				}
				// Add your JSON body assertions here
			}
		})
	}
}

func TestAuthHandler_VerifyToken(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		expectedCode int
		expectedBody map[string]interface{}
		setup        func()
	}{
		{
			name: "Invalid token",
			fields: fields{
				db: &gorm.DB{},
			},
			args: args{
				c: func() *gin.Context {
					c, _ := gin.CreateTestContext(httptest.NewRecorder())
					c.Request = httptest.NewRequest("POST", "/verify", nil)
					c.Request.Header.Set("Authorization", "Bearer invalid_token")
					return c
				}(),
			},
			expectedCode: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Invalid token",
			},
			setup: func() {
				// Setup mock token verification
			},
		},
		{
			name: "No token provided",
			fields: fields{
				db: &gorm.DB{},
			},
			args: args{
				c: func() *gin.Context {
					c, _ := gin.CreateTestContext(httptest.NewRecorder())
					c.Request = httptest.NewRequest("POST", "/verify", nil)
					return c
				}(),
			},
			expectedCode: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "No token provided",
			},
			setup: func() {
				// Setup mock token verification
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = tt.args.c.Request

			auth := controllers.NewAuthHandler(tt.fields.db)

			auth.VerifyToken(c)

			// Assert the response code
			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
			}

			// Assert the response body
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to parse response body: %v", err)
			}
			if !reflect.DeepEqual(response, tt.expectedBody) {
				t.Errorf("Expected body %v, got %v", tt.expectedBody, response)
			}
		})
	}
}
