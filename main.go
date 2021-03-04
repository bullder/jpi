package main

import (
	"github.com/coocood/freecache"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"io/ioutil"
	_ "jpi/docs"
	"log"
	"net/http"
	"os"
)

var cache = freecache.NewCache(1 * 1024 * 1024)

// @title JPI app
// @description This is a jpi management application
// @version 1.0
// @host localhost:1323
// @BasePath /
func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Smartbox!")
	})
	e.GET("/api", func(c echo.Context) error {
		return c.String(http.StatusOK, GetWeather())
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(":" + port))
}

// Get weather
// @Summary Get weather
// @Description Get weather description
// @Tags todos
// @Produce json
// @Success 201 {object}
// @Router /api [get]
func GetWeather() string {
	url := "https://community-open-weather-map.p.rapidapi.com/weather?q=Dublin%2Cie&units=metric"
	key := []byte("Dublin")
	cached := getCache(key)
	if cached != "" {
		log.Print("hit cache")
		return cached
	}

	log.Print("real request")

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-rapidapi-key", "f76e1ae731msh4629dae9182758fp1a95a4jsnb031d2d37365")
	req.Header.Add("x-rapidapi-host", "community-open-weather-map.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	setCache(key, body)

	return string(body)
}

func getCache(k []byte) string {
	got, err := cache.Get(k)
	if err != nil {
		return ""
	}
	return string(got)
}

func setCache(k []byte, v []byte)  {
	cache.Set(k, v, 60)
}
