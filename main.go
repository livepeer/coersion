package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Configuration contains run-time configuration, usually read from the environment
type Configuration struct {
	Port          uint16
	SignalTimeout time.Duration
}

func main() {

	log.Println("coersion started")
	defer log.Println("coersion stopped")

	cfg := Configuration{
		Port:          8088,
		SignalTimeout: 5 * time.Second,
	}

	ctx := context.Background()
	// set up the service first
	// XXX

	// Echo instance
	e := echo.New()

	//e.Use(middleware.LoggerWithConfig(middleware.DefaultLoggerConfig))
	// 	Format: "remote_ip=${remote_ip}, user_agent=${user_agent}, method=${method}, uri=${uri}, status=${status}\n"}))

	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET}, //{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	// Routes
	/**/
	routes(ctx, e)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: e,
	}

	err := s.ListenAndServe()
	if err != nil {
		log.Printf("echo server shutdown: %v\n", err)
	}
}
