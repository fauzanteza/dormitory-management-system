package routes

import (
	"dormitory-management/config"
	"dormitory-management/handlers"
	"dormitory-management/middleware"
	"net/http"
	"time"

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

	// Apply rate limiting globally
	router.Use(middleware.RateLimitMiddleware())

	// Public routes (no authentication required)
	router.POST("/api/auth/login", authHandler.Login)
	router.POST("/api/auth/register", authHandler.Register)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"version":   "1.0.0",
		})
	})

	// ================= ADMIN API ROUTES =================
	adminAPI := router.Group("/api/admin")
	adminAPI.Use(middleware.AuthMiddleware(config), middleware.AdminMiddleware())
	{
		// Dashboard
		adminAPI.GET("/dashboard/stats", dashboardHandler.GetDashboardStats)

		// Rooms Management
		adminAPI.GET("/rooms", roomHandler.GetRooms)
		adminAPI.GET("/rooms/available", roomHandler.GetAvailableRooms)
		adminAPI.GET("/rooms/:id", roomHandler.GetRoom)
		adminAPI.POST("/rooms", roomHandler.CreateRoom)
		adminAPI.PUT("/rooms/:id", roomHandler.UpdateRoom)
		adminAPI.DELETE("/rooms/:id", roomHandler.DeleteRoom)
		adminAPI.POST("/rooms/:id/assign", roomHandler.AssignResident)

		// Bookings Management
		adminAPI.GET("/bookings", bookingHandler.GetBookings)
		adminAPI.PUT("/bookings/:id/status", bookingHandler.UpdateBookingStatus)

		// Residents Management
		adminAPI.GET("/residents", residentHandler.GetResidents)
		adminAPI.GET("/residents/:id", residentHandler.GetResident)
		adminAPI.POST("/residents", residentHandler.CreateResident)
		adminAPI.PUT("/residents/:id", residentHandler.UpdateResident)
		adminAPI.DELETE("/residents/:id", residentHandler.DeleteResident)

		// Payments Management
		adminAPI.GET("/payments", paymentHandler.GetPayments)
		adminAPI.POST("/payments", paymentHandler.CreatePayment)
		adminAPI.PUT("/payments/:id/status", paymentHandler.UpdatePaymentStatus)
		adminAPI.GET("/payments/report", paymentHandler.GetMonthlyReport)

		// Repairs Management
		adminAPI.GET("/repairs", repairHandler.GetRepairRequests)
		adminAPI.PUT("/repairs/:id/status", repairHandler.UpdateRepairStatus)
		adminAPI.GET("/repairs/stats", repairHandler.GetRepairStats)

		// Users Management
		adminAPI.GET("/users", userHandler.GetUsers)
		adminAPI.POST("/users/:id/reset-password", userHandler.ResetPassword)
		adminAPI.DELETE("/users/:id", userHandler.DeleteUser)

		// Profile
		adminAPI.GET("/profile", authHandler.GetProfile)
		adminAPI.PUT("/profile", authHandler.UpdateProfile)
		adminAPI.POST("/change-password", authHandler.ChangePassword)
	}

	// ================= STUDENT API ROUTES =================
	studentAPI := router.Group("/api/student")
	studentAPI.Use(middleware.AuthMiddleware(config))
	{
		// Student Dashboard
		studentAPI.GET("/dashboard/stats", dashboardHandler.GetStudentDashboardStats)

		// Room - View available rooms and own room
		studentAPI.GET("/rooms/available", roomHandler.GetAvailableRooms)
		studentAPI.GET("/my-room", roomHandler.GetMyRoom)

		// Bookings - Create and view own bookings
		studentAPI.GET("/bookings", bookingHandler.GetBookings)
		studentAPI.POST("/bookings", bookingHandler.CreateBooking)

		// Payments - View own payments
		studentAPI.GET("/my-payments", paymentHandler.GetMyPayments)

		// Repairs - Create and view own repairs
		studentAPI.GET("/my-repairs", repairHandler.GetRepairRequests)
		studentAPI.POST("/repairs", repairHandler.CreateRepairRequest)

		// Profile
		studentAPI.GET("/profile", authHandler.GetProfile)
		studentAPI.PUT("/profile", authHandler.UpdateProfile)
		studentAPI.POST("/change-password", authHandler.ChangePassword)
	}

	// ================= COMMON ROUTES =================
	commonAPI := router.Group("/api")
	commonAPI.Use(middleware.AuthMiddleware(config))
	{
		// Logout
		commonAPI.POST("/auth/logout", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
		})
	}

	// ================= STATIC FILES SERVING =================
	// Serve admin frontend
	router.Static("/admin", "../frontend/admin")

	// Serve student frontend
	router.Static("/student", "../frontend/student")

	// Serve auth pages
	router.Static("/auth", "../frontend/auth")

	// Serve common resources
	router.Static("/common", "../frontend/common")

	// Root redirect based on role
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/auth/login.html")
	})

	// Login and register shortcuts
	router.GET("/login", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/auth/login.html")
	})

	router.GET("/register", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/auth/register.html")
	})
}
