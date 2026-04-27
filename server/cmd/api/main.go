package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aboloredev/armory/internal/config"
	"github.com/aboloredev/armory/internal/database"
	"github.com/aboloredev/armory/internal/events"
	"github.com/aboloredev/armory/internal/interfaces"
	"github.com/aboloredev/armory/internal/logger"
	"github.com/aboloredev/armory/internal/providers"
	"github.com/aboloredev/armory/internal/server"
	"github.com/aboloredev/armory/internal/services"
	"github.com/gin-gonic/gin"
)

// @title Armory E-Commerce API
// @version 1.0
// @description A modern e-commerce API built with Go, Gin, and GORM
// @termsOfService http://swagger.io/terms/

// @contact.name   Alabi Fathiu
// @contact.url    http://linkedin.com/in/fathiu-alabi
// @contact.email  no-email@no-email

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemas http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	log := logger.New()

	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("could not load config")
	}

	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to database")
	}

	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get database connection")
	}
	defer mainDB.Close()

	// define context so the event can use it
	ctx := context.Background()

	publisherEvent, err := events.NewEventPublisher(
		ctx,
		&config.AWSConfig{
			EventQueueName: cfg.AWSServices.EventQueueName,
			Key:            cfg.AWSServices.Key,
			KeyId:          cfg.AWSServices.KeyId,
			S3Endpoint:     cfg.AWSServices.S3Endpoint,
			Region:         cfg.AWSServices.S3Endpoint,
		})
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialise events")
		return
	}
	gin.SetMode(cfg.Server.GinMode)

	authService := services.NewAuthService(db, cfg)
	productService := services.NewProductService(db)
	userService := services.NewUserService(db)
	cartService := services.NewCartServices(db)
	orderServices := services.NewOrderService(db, publisherEvent)

	// switch between upload providers (s3 and local)
	var uploadProvider interfaces.UploadProvider
	if cfg.Uploads.UploadProvider == "s3" {
		uploadProvider = providers.NewS3Provider(cfg)
	} else {
		uploadProvider = providers.NewLocalUploadProvider(cfg.Uploads.Path)
	}

	uploadService := services.NewUploadService(uploadProvider)

	srv := server.New(
		cfg,
		db,
		log,
		authService,
		productService,
		userService,
		uploadService,
		cartService,
		orderServices,
	)

	router := srv.SetupRoutes()

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Info().Str("port", cfg.Server.Port).Msg("http server started")
		err = httpServer.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("Failed to start http server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to shutdown server")
		return
	}

	log.Info().Msg("Shutting down server")
}
