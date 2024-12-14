package controllers

import (
	"bytes"
	"mime/multipart"
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

func TestNewInventoryManagementHandler(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *controllers.InventoryManagementHandler
	}{
		{
			name: "Valid DB instance",
			args: args{db: &gorm.DB{}},
			want: &controllers.InventoryManagementHandler{Db: &gorm.DB{}},
		},
		{
			name: "Nil DB instance",
			args: args{db: nil},
			want: &controllers.InventoryManagementHandler{Db: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := controllers.NewInventoryManagementHandler(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInventoryManagementHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInventoryManagementHandler_CreateProduct(t *testing.T) {
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
			name:   "User not authenticated",
			fields: fields{db: &gorm.DB{}},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Set("userID", -1)
				return c
			}()},
			expectedCode: http.StatusUnauthorized,
			expectedBody: true,
			setup:        func() {},
		},

		{
			name:   "Invalid file type",
			fields: fields{db: &gorm.DB{}},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("image", "test.txt")
				part.Write([]byte("This is a test file"))
				writer.Close()
				c.Request = httptest.NewRequest("POST", "/product", body)
				c.Request.Header.Set("Content-Type", writer.FormDataContentType())
				c.Set("userID", uint(1))
				return c
			}()},
			expectedCode: http.StatusBadRequest,
			expectedBody: true,
			setup:        func() {},
		},
		{
			name: "Invalid product creation",
			fields: fields{db: func() *gorm.DB {
				db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				db.AutoMigrate(&models.Product{})
				return db
			}()},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				writer.WriteField("name", "Test Product")
				writer.WriteField("price", "100")
				writer.WriteField("quantity", "10")
				writer.Close()
				c.Request = httptest.NewRequest("POST", "/product", body)
				c.Request.Header.Set("Content-Type", writer.FormDataContentType())
				c.Set("userID", uint(1))
				return c
			}()},
			expectedCode: http.StatusBadRequest,
			expectedBody: true,
			setup:        func() {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			im := &controllers.InventoryManagementHandler{
				Db: tt.fields.db,
			}

			// Create a recorder first
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = tt.args.c.Request
			c.Set("userID", tt.args.c.GetUint("userID"))

			im.CreateProduct(c)

			// Now check the actual response
			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, but got %d", tt.expectedCode, w.Code)
			}
		})
	}
}

func TestInventoryManagementHandler_DeleteProduct(t *testing.T) {
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
		expectedBody string
		setup        func()
	}{
		{
			name: "Successful product deletion",
			fields: fields{db: func() *gorm.DB {
				db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				db.AutoMigrate(&models.Product{})
				db.Create(&models.Product{ID: 1, UserID: 1, Name: "Test Product", Active: true})
				return db
			}()},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httptest.NewRequest("DELETE", "/product/1", nil)
				c.Set("userID", uint(1))
				c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
				return c
			}()},
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"Product marked as inactive successfully"}`,
			setup:        func() {},
		},

		{
			name: "Unauthorized user",
			fields: fields{db: func() *gorm.DB {
				db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				db.AutoMigrate(&models.Product{})
				return db
			}()},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httptest.NewRequest("DELETE", "/product/1", nil)
				c.Set("userID", uint(0))
				c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
				return c
			}()},
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"Unauthorized"}`,
			setup:        func() {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			im := &controllers.InventoryManagementHandler{
				Db: tt.fields.db,
			}

			// Create a recorder first
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = tt.args.c.Request
			c.Set("userID", tt.args.c.GetUint("userID"))
			c.Params = tt.args.c.Params

			im.DeleteProduct(c)

			// Now check the actual response
			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, but got %d", tt.expectedCode, w.Code)
			}
			if w.Body.String() != tt.expectedBody {
				t.Errorf("Expected body %s, but got %s", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestInventoryManagementHandler_GetAllProducts(t *testing.T) {
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
		expectedBody string
		setup        func()
	}{
		{
			name: "User not authenticated",
			fields: fields{db: func() *gorm.DB {
				db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				return db
			}()},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httptest.NewRequest("GET", "/products", nil)
				c.Set("userID", uint(0))
				return c
			}()},
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"Unauthorized"}`,
			setup:        func() {},
		},
		{
			name: "No products found",
			fields: fields{db: func() *gorm.DB {
				db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				db.AutoMigrate(&models.Product{})
				return db
			}()},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httptest.NewRequest("GET", "/products", nil)
				c.Set("userID", uint(1))
				return c
			}()},
			expectedCode: http.StatusOK,
			expectedBody: `null`,
			setup:        func() {},
		},
		{
			name: "Database error",
			fields: fields{db: func() *gorm.DB {
				db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				return db
			}()},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httptest.NewRequest("GET", "/products", nil)
				c.Set("userID", uint(1))
				return c
			}()},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"Failed to get products"}`,
			setup: func() {
				db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				sqlDB, _ := db.DB()
				sqlDB.Close() // Simulate a database error
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			im := &controllers.																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																												InventoryManagementHandler{
				Db: tt.fields.db,
			}

			// Create a recorder first
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = tt.args.c.Request
			c.Set("userID", tt.args.c.GetUint("userID"))

			im.GetAllProducts(c)

			// Now check the actual response
			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, but got %d", tt.expectedCode, w.Code)
			}
			if w.Body.String() != tt.expectedBody {
				t.Errorf("Expected body %s, but got %s", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestInventoryManagementHandler_GetLowStockAlerts(t *testing.T) {
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
		expectedBody string
		setup        func()
	}{
		{
			name: "User not authenticated",
			fields: fields{db: func() *gorm.DB {
				db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				return db
			}()},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httptest.NewRequest("GET", "/low-stock-alerts", nil)
				return c
			}()},
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"Unauthorized"}`,
			setup:        func() {},
		},
		{
			name: "No low stock alerts",
			fields: fields{db: func() *gorm.DB {
				db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				db.AutoMigrate(&models.LowStockAlert{}, &models.Product{}, &models.Inventory{})
				return db
			}()},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httptest.NewRequest("GET", "/low-stock-alerts", nil)
				c.Set("userID", uint(1))
				return c
			}()},
			expectedCode: http.StatusOK,
			expectedBody: `[]`,
			setup:        func() {},
		},
		{
			name: "Database error",
			fields: fields{db: func() *gorm.DB {
				db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				return db
			}()},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httptest.NewRequest("GET", "/low-stock-alerts", nil)
				c.Set("userID", uint(1))
				return c
			}()},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"Failed to fetch alerts"}`,
			setup: func() {
				db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				sqlDB, _ := db.DB()
				sqlDB.Close() // Simulate a database error
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			im := &controllers.InventoryManagementHandler{
				Db: tt.fields.db,
			}

			// Create a recorder first
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = tt.args.c.Request
			c.Set("userID", tt.args.c.GetUint("userID"))

			im.GetLowStockAlerts(c)

			// Now check the actual response
			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, but got %d", tt.expectedCode, w.Code)
			}
			if w.Body.String() != tt.expectedBody {
				t.Errorf("Expected body %s, but got %s", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestInventoryManagementHandler_SearchProducts(t *testing.T) {
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
		expectedBody string
		setup        func()
	}{
		{
			name: "No search query",
			fields: fields{db: func() *gorm.DB {
				db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				return db
			}()},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httptest.NewRequest("GET", "/search-products", nil)
				c.Set("userID", uint(1))
				return c
			}()},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Search query is required"}`,
			setup:        func() {},
		},
		{
			name: "Database error",
			fields: fields{db: func() *gorm.DB {
				db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				return db
			}()},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httptest.NewRequest("GET", "/search-products?q=Test", nil)
				c.Set("userID", uint(1))
				return c
			}()},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"Failed to search products"}`,
			setup: func() {
				db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				sqlDB, _ := db.DB()
				sqlDB.Close() // Simulate a database error
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			im := &controllers.InventoryManagementHandler{
				Db: tt.fields.db,
			}

			// Create a recorder first
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = tt.args.c.Request
			c.Set("userID", tt.args.c.GetUint("userID"))

			im.SearchProducts(c)

			// Now check the actual response
			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, but got %d", tt.expectedCode, w.Code)
			}
			if w.Body.String() != tt.expectedBody {
				t.Errorf("Expected body %s, but got %s", tt.expectedBody, w.Body.String())
			}
		})
	}
}
