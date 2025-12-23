package routes

import (
	"dormitory-management/config"
	"dormitory-management/handlers"
	"dormitory-management/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, config *config.Config) {
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, config)
	roomHandler := handlers.NewRoomHandler(db)
	paymentHandler := handlers.NewPaymentHandler(db)
	repairHandler := handlers.NewRepairHandler(db)
	dashboardHandler := handlers.NewDashboardHandler(db)
	residentHandler := handlers.NewResidentHandler(db)
	userHandler := handlers.NewUserHandler(db)

	bookingHandler := handlers.NewBookingHandler(db)

	// Auth routes (no authentication required)
	router.POST("/api/auth/login", authHandler.Login)
	router.POST("/api/auth/register", authHandler.Register)

	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(config))

	// Profile
	api.GET("/profile", authHandler.GetProfile)
	api.PUT("/profile", authHandler.UpdateProfile)
	api.POST("/profile/password", authHandler.ChangePassword)

	// Dashboard
	api.GET("/dashboard/stats", dashboardHandler.GetDashboardStats)

	// Rooms
	rooms := api.Group("/rooms")
	{
		rooms.GET("", roomHandler.GetRooms)
		rooms.GET("/available", roomHandler.GetAvailableRooms)
		rooms.GET("/:id", roomHandler.GetRoom)
		rooms.POST("", middleware.AdminMiddleware(), roomHandler.CreateRoom)
		rooms.PUT("/:id", middleware.AdminMiddleware(), roomHandler.UpdateRoom)
		rooms.DELETE("/:id", middleware.AdminMiddleware(), roomHandler.DeleteRoom)
		rooms.POST("/:id/assign", middleware.AdminMiddleware(), roomHandler.AssignResident)
	}

	// Bookings
	bookings := api.Group("/bookings")
	{
		bookings.GET("", bookingHandler.GetBookings)
		bookings.POST("", bookingHandler.CreateBooking)
		bookings.PUT("/:id/status", middleware.AdminMiddleware(), bookingHandler.UpdateBookingStatus)
	}

	// Payments
	payments := api.Group("/payments")
	{
		payments.GET("", paymentHandler.GetPayments)
		payments.POST("", middleware.AdminMiddleware(), paymentHandler.CreatePayment)
		payments.PUT("/:id/status", middleware.AdminMiddleware(), paymentHandler.UpdatePaymentStatus)
		payments.GET("/report", paymentHandler.GetMonthlyReport)
	}

	// Repair requests
	repairs := api.Group("/repairs")
	{
		repairs.GET("", repairHandler.GetRepairRequests)
		repairs.POST("", repairHandler.CreateRepairRequest)
		repairs.PUT("/:id/status", middleware.AdminMiddleware(), repairHandler.UpdateRepairStatus)
		repairs.GET("/stats", repairHandler.GetRepairStats)
	}

	// Residents
	residents := api.Group("/residents")
	residents.Use(middleware.AdminMiddleware())
	{
		residents.GET("", residentHandler.GetResidents)
		residents.GET("/:id", residentHandler.GetResident)
		residents.POST("", residentHandler.CreateResident)
		residents.PUT("/:id", residentHandler.UpdateResident)
		residents.DELETE("/:id", residentHandler.DeleteResident)
	}

	// Users
	users := api.Group("/users")
	users.Use(middleware.AdminMiddleware())
	{
		users.GET("", userHandler.GetUsers)
		users.POST("/:id/reset-password", userHandler.ResetPassword)
		users.DELETE("/:id", userHandler.DeleteUser)
	}

	// Serve static files
	router.Static("/frontend", "../frontend")
	router.StaticFile("/", "../frontend/html/index.html")
	router.StaticFile("/login", "../frontend/html/login.html")
	router.StaticFile("/dashboard", "../frontend/html/dashboard.html")
	router.StaticFile("/rooms", "../frontend/html/rooms.html")
	router.StaticFile("/payments", "../frontend/html/payments.html")
	router.StaticFile("/repairs", "../frontend/html/repairs.html")
	router.StaticFile("/residents", "../frontend/html/residents.html")
	router.StaticFile("/users", "../frontend/html/users.html")
	router.StaticFile("/profile", "../frontend/html/profile.html")
	router.StaticFile("/bookings", "../frontend/html/bookings.html")
}
