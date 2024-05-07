package main

import (
	"encoding/json"
	"flag"
	"io"
	"math"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

type WeatherValues struct {
	Temperature float64 `json:"temperature"`
}

type WeatherData struct {
	Values WeatherValues `json:"values"`
}

type WeatherRes struct {
	Data WeatherData `json:"data"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file")
	}
}

// adding this comment literally just to test the release pushing
func main() {
	day := flag.String("day", "today", "What day's weather do you want to see?")
	flag.Parse()

	if flag.NArg() == 0 {
		log.Printf("Here is the weather for %v!\n", *day)
		getWeather()
	} else if flag.Arg(0) == "day" {
		log.Printf("Here is the weather for %v!\n", *day)
		getWeather()
	}
}

func getWeather() {
	apiKey := os.Getenv("TOMORROW_API_KEY")
	url := "https://api.tomorrow.io/v4/weather/realtime?location=64109&apikey=" + apiKey

	req, _ := http.NewRequest("GET", url, nil)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error in api call - %v", err.Error())
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var w WeatherRes
	err = json.Unmarshal(body, &w)
	if err != nil {
		log.Fatalf("error parsing json response - %v", err.Error())
	}

	c := w.Data.Values.Temperature
	f := math.Round((c * 1.8) + 32)

	log.Printf("Temp: %v°F", f)
}
