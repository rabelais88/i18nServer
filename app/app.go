package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func onHealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "SUCCESS")
}

func Start() {
	godotenv.Load()
	e := echo.New()
	address := fmt.Sprintf(":%s", os.Getenv("PORT"))
	e.Start(address)

	e.POST("/publish", onPublish)
	e.GET("/health", onHealthCheck)
}
