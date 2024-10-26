package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Weather struct {
	Coordinates struct {
		Longitude float64 `json:"lon"`
		Latitude  float64 `json:"lat"`
	} `json:"coord"`
	Main struct {
		Temperature float64 `json:"temp"`
	} `json:"main"`
	Name string `json:"name"`
}

func main() {

	// place := "Boston"

	load := godotenv.Load()
	if load != nil {
		panic(load)
	}

	key := os.Getenv("API_KEY")
	/*
		if len(os.Args) >= 2 {
			place = os.Args[1]
		}
	*/
	res, err := http.Get("http://api.openweathermap.org/data/2.5/weather?lat=42.3611&lon=-71.0570&appid=" + key)
	if err != nil {
		panic(err)

	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("API not available")

	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	degrees, name := weather.Main.Temperature, weather.Name

	fmt.Printf("%s, %.0fK \n",
		name,
		degrees,
	)

}
