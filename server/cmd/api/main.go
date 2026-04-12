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
	"github.com/aboloredev/armory/internal/logger"
	"github.com/aboloredev/armory/internal/server"
	"github.com/gin-gonic/gin"
)

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
	gin.SetMode(cfg.Server.GinMode)

	srv := server.New(cfg, db, log)

	router := srv.SetupRoutes()

	httpServer := &http.Server{
		Addr: fmt.Sprintf("%s", cfg.Server.Port),
		Handler: router,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func(){
		log.Info().Str("port", cfg.Server.Port).Msg("http server started")
		err := httpServer.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed){
			log.Fatal().Err(err).Msg("Failed to start http server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
	defer cancel()

	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to shutdown server")
	}

	log.Info().Msg("Shutting down server")
}