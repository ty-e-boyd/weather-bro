package main

import (
	"encoding/json"
	"flag"
	"fmt"
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

func main() {
	// on each run, decide if it's the first run by checking for the stored key. No stored key means first run
	// on first run, prompt for tomorrow.io api key
	// get the key, test it, on valid key, store it somehow locally
	// then on this, and subsequent runs, retreive the key and os.Setenv()
	// now the rest of the app can function with an api key that the user is responsible for

	// if first run only
	log.Print("Please enter your tomorrow.io api key..")
	var userKey string
	fmt.Scanln(&userKey)

	// this isn't workig, use creatfile, open it, defer, write to a text file, then figure out opening
	err := os.WriteFile("tmp/secret", []byte(userKey), 0o644)
	if err != nil {
		log.Warnf("error writing api key - %v", err.Error())
	}

	log.Warnf("key entered: %v", userKey)
	os.Setenv("TOMORROW_API_KEY_2", userKey)

	// the rest of the app
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
	apiKey := os.Getenv("TOMORROW_API_KEY_2")
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
