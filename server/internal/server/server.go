package server

import (
	"net/http"

	_ "github.com/aboloredev/armory/docs"

	"github.com/aboloredev/armory/internal/config"
	"github.com/aboloredev/armory/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config         *config.Config
	database       *gorm.DB
	logger         zerolog.Logger
	authService    *services.AuthService
	productService *services.ProductService
	userService    *services.UserService
	uploadService  *services.UploadService
	cartService    *services.CartServices
	orderService   *services.OrderService
}

func New(
	cfg *config.Config, db *gorm.DB,
	log zerolog.Logger,
	authService *services.AuthService,
	productService *services.ProductService,
	userService *services.UserService,
	uploadService *services.UploadService,
	cartService *services.CartServices,
	orderServies *services.OrderService) *Server {
	return &Server{
		config:         cfg,
		database:       db,
		logger:         log,
		authService:    authService,
		productService: productService,
		userService:    userService,
		uploadService:  uploadService,
		cartService:    cartService,
		orderService:   orderServies,
	}

}

func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.New()

	// Middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.CORS())

	// Health check
	router.GET("/health", s.HealthCheck)

	// Add documentation routes
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.StaticFile("/api-docs", "./docs/rapiddoc.html")

	// stattic uploads
	router.Static("/uploads", "./uploads")

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			// Auth Routes
			auth.POST("/register", s.RegisterHandler)
			auth.POST("/login", s.LoginHandler)
			auth.POST("/logout", s.LogoutHandler)
			auth.POST("/refresh", s.RefreshTokenHandler)
		}
		protected := api.Group("/")
		protected.Use(s.AuthMiddleware())
		{
			users := protected.Group("/users")
			{
				// Users Routes
				users.GET("/profile", s.GetUserProfileHandler)
				users.PATCH("/profile", s.UpdateUserProfileHandler)
				users.GET("/all", s.AdminMiddleware(), s.GetAllUsersHandler)
			}

			categories := protected.Group("/categories")
			{ // Categories Route
				categories.POST("/", s.AdminMiddleware(), s.CreateCategoryHandler)
				categories.PATCH("/:id", s.AdminMiddleware(), s.UpdateCategoryHandler)
				categories.DELETE("/:id", s.AdminMiddleware(), s.DeleteCategoryHandler)
				categories.GET("/", s.AdminMiddleware(), s.GetAllCategoriesHandler)
			}

			products := protected.Group("/products")
			{
				// Products Routes
				products.POST("/", s.AdminMiddleware(), s.CreateProductHandler)
				products.PATCH("/:id", s.AdminMiddleware(), s.UpdateProductHandler)
				products.DELETE("/:id", s.AdminMiddleware(), s.DeleteProduct)
				products.POST("/:id/images", s.AdminMiddleware(), s.UploadProductImageService)
				products.DELETE("/:id/images", s.AdminMiddleware(), s.DeleteProductImageService)
			}
			cart := protected.Group("/cart")
			{
				// Cart Routes
				cart.GET("/", s.GetUserCart)
				cart.POST("/", s.AddToCartHandler)
				cart.PATCH("/:id", s.UpdateCartItemHandler)
				cart.DELETE("/:id", s.RemoveCartItemHandler)
				cart.DELETE("/", s.ClearUserCartHandler)
			}
			order := protected.Group("/orders")
			{
				// Order routes
				order.GET("/", s.GetAllOrdersHandler)
				order.POST("/", s.CreateOrderHandler)
				order.GET("/:id", s.GetOrderHandler)
			}
		}

		// Publuc route
		api.GET("/category/all", s.GetAllCategoriesHandler)
		api.GET("/products/all", s.AllProductHandler)
		api.GET("/products/:id", s.GetProductHandler)
	}

	return router
}


func (s *Server) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (s *Server) CORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}
