package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OAthooh/BiasharaTrack.git/controllers"
	"github.com/OAthooh/BiasharaTrack.git/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewCreditManager(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *controllers.CreditManager
	}{
		{
			name: "Valid DB instance",
			args: args{db: &gorm.DB{}},
			want: &controllers.CreditManager{Db: &gorm.DB{}},
		},
		{
			name: "Nil DB instance",
			args: args{db: nil},
			want: &controllers.CreditManager{Db: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := controllers.NewCreditManager(tt.args.db); got.Db != tt.want.Db {
				t.Errorf("NewCreditManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreditManager_GetCreditsHistory(t *testing.T) {
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
			name: "Valid user with transactions",
			fields: fields{db: func() *gorm.DB {
				db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				db.AutoMigrate(&models.CreditTransaction{})
				db.Create(&models.CreditTransaction{UserID: 1, CreditAmount: 100, BalanceDue: 50, Status: "PENDING"})
				return db
			}()},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Set("userID", uint(1))
				return c
			}()},
			expectedCode: http.StatusOK,
			expectedBody: true,
			setup:        func() {},
		},
		{
			name: "Valid user with no transactions",
			fields: fields{db: func() *gorm.DB {
				db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				db.AutoMigrate(&models.CreditTransaction{})
				return db
			}()},
			args: args{c: func() *gin.Context {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Set("userID", uint(1))
				return c
			}()},
			expectedCode: http.StatusOK,
			expectedBody: true,
			setup:        func() {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			cm := controllers.NewCreditManager(tt.fields.db)
			cm.GetCreditsHistory(tt.args.c)
			if tt.expectedBody {
				w := httptest.NewRecorder().Result()
				if w.StatusCode != tt.expectedCode {
					t.Errorf("Expected status code %d, but got %d", tt.expectedCode, w.StatusCode)
				}
			}
		})
	}
}
