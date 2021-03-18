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
	"golang.org/x/net/websocket"
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
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Static("/", "./public")

	e.GET("/api/weather/:id", func(c echo.Context) error {
		return c.JSON(http.StatusOK, GetWeather(c.Param("id")))
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

	e.GET("/ws", ws)

	e.GET("/api/weatherCache/:id", func(c echo.Context) error {
		return c.String(http.StatusOK, GetWeatherCache(c.Param("id")))
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(":" + port))
}

func ws(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			err := websocket.Message.Send(ws, "Hello!")
			if err != nil {
				c.Logger().Error(err)
			}

			msg := ""
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
			}
			if msg != "" {
				GetCitiesWs(ws, msg)
			}
			log.Printf("%s\n", msg)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
