package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"

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

	rf, err := os.ReadFile("secret.txt")
	if err != nil {
		log.Fatalf("error reading file - %v", err.Error())
	}

	var foundKey string
	lines := strings.Split(string(rf), "\n")
	for i, line := range lines {
		if i == 0 {
			fmt.Println(line)
			foundKey = line
		}
	}

	if len(foundKey) > 0 {
		os.Setenv("TOMORROW_API_KEY", foundKey)
	} else {
		// if first run only
		var userKey string
		log.Print("Please enter your tomorrow.io api key..")
		fmt.Scanln(&userKey)

		// write key to local file
		f, err := os.Create("secret.txt")
		if err != nil {
			log.Fatalf("error creating file - %v", err.Error())
		}
		defer f.Close()

		wf, err := f.WriteString(userKey)
		if err != nil {
			log.Fatalf("error writing to file - %v", err.Error())
		}
		log.Warnf("wrote bytes: %v", wf)

		log.Warnf("key entered: %v", userKey)
		os.Setenv("TOMORROW_API_KEY", userKey)
	}

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
	apiKey := os.Getenv("TOMORROW_API_KEY")
	url := "https://api.tomorrow.io/v4/weather/realtime?location=64109&apikey=" + apiKey

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("error creating new request - %v", err.Error())
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error in api call - %v", err.Error())
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error reading body from api call - %v", err.Error())
	}

	var w WeatherRes
	err = json.Unmarshal(body, &w)
	if err != nil {
		log.Fatalf("error parsing json response - %v", err.Error())
	}

	c := w.Data.Values.Temperature
	f := math.Round((c * 1.8) + 32)

	log.Printf("Temp: %v°F", f)
}
