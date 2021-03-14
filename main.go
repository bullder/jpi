package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/newrelic/go-agent/v3/integrations/logcontext/nrlogrusplugin"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	"github.com/newrelic/go-agent/v3/integrations/nrlogrus"
	"github.com/newrelic/go-agent/v3/newrelic"
	log "github.com/sirupsen/logrus"
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
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(nrlogrusplugin.ContextFormatter{})

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("JPI"),
		newrelic.ConfigLicense("eu01xx97deebf972bfbfced797e1a1757807NRAL"),
		newrelic.ConfigDistributedTracerEnabled(true),
		func(config *newrelic.Config) {
			log.SetLevel(log.DebugLevel)
			config.Logger = nrlogrus.StandardLogger()
		},
	)
	if nil != err {
		log.Println(err)
		os.Exit(1)
	}
	e := echo.New()
	e.Use(nrecho.Middleware(app))
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
