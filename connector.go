package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// GetWeather
// @Summary Get weather
// @Description Get weather description
// @Tags todos
// @Produce json
// @Success 201 {object}
// @Router /api/weather/:id [get]
func GetWeather(c string) string {
	body := GetResponse(c)

	return string(body)
}

// GetWeatherCache
// @Summary Get weather
// @Description Get weather description
// @Tags todos
// @Produce json
// @Success 201 {object}
// @Router /api/weatherCache/:id [get]
func GetWeatherCache(c string) string {
	body := GetCachedResponse(c)

	return body
}

// GetCity
// @Summary Get weather
// @Description Get weather description
// @Tags todos
// @Produce json
// @Success 201 {object}
// @Router /api/city/:id [get]
func GetCity(c string) CityResult {
	body := GetResponse(c)
	var city City
	err := json.Unmarshal(body, &city)

	if err != nil {
		log.Println("error:", err)
	}
	return city.getResult()
}

// GetCities
// @Summary Get weather
// @Description Get weather description
// @Tags todos
// @Produce json
// @Success 201 {object}
// @Router /api/cities/:id [get]
func GetCities(c string) CitiesResult {
	start := time.Now()
	var r CitiesResult
	cities := strings.Split(c, ",")
	for _, city := range cities {
		cityResult := GetCity(city)
		r.addCity(cityResult)
	}
	r.Duration = time.Since(start).Milliseconds()

	return r
}

// GetCitiesAsync
// @Summary Get weather
// @Description Get weather description
// @Tags todos
// @Produce json
// @Success 201 {object}
// @Router /api/citiesAsync/:id [get]
func GetCitiesAsync(c string) CitiesResult {
	start := time.Now()
	var r CitiesResult
	cities := strings.Split(c, ",")

	var wg sync.WaitGroup

	for _, city := range cities {
		wg.Add(1)
		go worker(city, &r, &wg)
	}

	wg.Wait()
	r.Duration = time.Since(start).Milliseconds()

	return r
}

func worker(c string, r *CitiesResult, wg *sync.WaitGroup) {
	defer wg.Done()
	r.addCity(GetCity(c))
}

func GetCitiesWs(ws *websocket.Conn, c string) {
	cities := strings.Split(c, ",")

	for _, city := range cities {
		go workerWs(city, ws)
	}
}

func workerWs(c string, ws *websocket.Conn) {
	msg, err := json.Marshal(GetCity(c))
	if err != nil {
		log.Println("error:", err)
	}

	err = websocket.Message.Send(ws, msg)
	if err != nil {
		log.Error(err)
	}
}

func GetCachedResponse(c string) string {
	if "" == c {
		c = "Dublin"
	}

	k := []byte(c)
	cached := getCache(k)
	if cached != nil {
		log.Print("hit cache")
		return string(cached)
	}

	log.Print("miss cache")

	body := GetResponse(c)

	setCache(k, body)

	return string(body)
}

func GetResponse(c string) []byte {
	req, err := http.NewRequest("GET", GetUrl(c), nil)
	if err != nil {
		log.Print(err)
	}

	res, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print("unclosed connection")
		}
	}(res.Body)

	body, _ := ioutil.ReadAll(res.Body)

	return body
}

func GetUrl(c string) string {
	apiUrl, err := url.Parse("https://api.openweathermap.org/data/2.5/weather")
	if err != nil {
		log.Print(err)
	}

	params := url.Values{}
	params.Add("APPID", "2dab67244cb2e115701cde03e5f9a7f7")
	params.Add("units", "metric")
	params.Add("q", c)

	apiUrl.RawQuery = params.Encode()

	return apiUrl.String()
}
