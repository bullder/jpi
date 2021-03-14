package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "jpi/docs"
	"net/http"
	"os"
)

// @title JPI app
// @description This is a jpi management application
// @version 1.0
// @host localhost:1323
// @BasePath /
func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Smartbox!")
	})
	e.GET("/api/weather/:id", func(c echo.Context) error {
		return c.String(http.StatusOK, GetWeather(c.Param("id")))
	})

	e.GET("/api/city/:id", func(c echo.Context) error {
		return c.JSON(http.StatusOK, GetCity(c.Param("id")))
	})

	e.GET("/api/cities/:id", func(c echo.Context) error {
		return c.JSON(http.StatusOK, GetCities(c.Param("id")))
	})

	e.GET("/api/citiesAsync/:id", func(c echo.Context) error {
		return c.JSON(http.StatusOK, GetCitiesAsync(c.Param("id")))
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(":" + port))
}
