package main

import (
    "fmt"
    "log"
    "os"
    "github.com/OAthooh/BiasharaTrack.git/database"
    "github.com/OAthooh/BiasharaTrack.git/routes"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func main() {
    // Get CALLBACK_URL directly from environment
    callbackURL := os.Getenv("CALLBACK_URL")
    if callbackURL == "" {
        log.Println("Warning: CALLBACK_URL is not set in environment variables")
    }

    fmt.Println("Initializing database connection...")
    // Initialize database connection
    db, err := database.Connect()
    if err != nil || db == nil {
        log.Fatal(err)
    }
    defer db.Close()
    fmt.Println("Database connection initialized successfully")

    // Run database migrations
    fmt.Println("Running database migrations...")
    err = db.Migrate()
    if err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }
    fmt.Println("Database migrations completed successfully")

    // Initialize Gin router with default middleware
    fmt.Println("Initializing Gin router...")
    router := gin.Default()

    // Configure CORS with expanded settings
    config := cors.DefaultConfig()
    config.AllowOrigins = []string{
        "http://localhost:5173",
        "https://biasharatrack-frontend.onrender.com",
        callbackURL,
    }
    config.AllowCredentials = true
    config.AllowHeaders = []string{
        "Origin",
        "Content-Type",
        "Accept",
        "Authorization",
        "Access-Control-Allow-Origin",
        "Access-Control-Allow-Headers",
        "Access-Control-Allow-Methods",
    }
    config.AllowMethods = []string{
        "GET",
        "POST",
        "PUT",
        "PATCH",
        "DELETE",
        "OPTIONS",
    }
    config.ExposeHeaders = []string{
        "Content-Length",
        "Authorization",
        "Access-Control-Allow-Origin",
    }
    router.Use(cors.New(config))

    fmt.Println("Gin router initialized successfully")

    // Register routes
    routes.AuthRoutes(router, db.DB)
    routes.CreditRoutes(router, db.DB)
    routes.InventoryManagementRoutes(router, db.DB)
    routes.SalesManagementRoutes(router, db.DB)
    routes.MpesaRoutes(router, db.DB)
    routes.SetupReceiptRoutes(router, db.DB)

    // Get port from environment variable or use default
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    fmt.Printf("Server is running on port %s\n", port)
    // Start server using the PORT environment variable
    router.Run(":" + port)
}