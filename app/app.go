package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func onHealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "SUCCESS")
}

func Start() {
	err := godotenv.Load()
	if err != nil {
		log.Print("cannot read .env; please refer .env.sample for how to make it")
	}
	e := echo.New()
	address := fmt.Sprintf(":%s", os.Getenv("PORT"))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: "method=${method}, uri=${uri}, status=${status}\n"}))

	e.POST("/publish", onPublish)
	e.GET("/health", onHealthCheck)

	e.Start(address)
}
