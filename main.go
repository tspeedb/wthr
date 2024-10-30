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

type Coord struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"lon"`
}

func main() {

	load := godotenv.Load()
	if load != nil {
		panic(load)
	}

	key := os.Getenv("API_KEY")
	lat := "42.3611"
	lon := "-71.0570"

	// Use Geocoder API to get lat and long of place in args if longer than 2

	if len(os.Args) >= 2 {
		place := os.Args[1]

		res, err := http.Get("http://api.openweathermap.org/geo/1.0/direct?q=" + place + "&limit=5&appid=" + key)
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
		var coords []Coord

		err = json.Unmarshal(body, &coords)
		if err != nil {
			panic(err)
		}

		lat, lon = fmt.Sprintf("%f", coords[0].Lat), fmt.Sprintf("%f", coords[0].Long)

	}

	res, err := http.Get("http://api.openweathermap.org/data/2.5/weather?lat=" + lat + "&lon=" + lon + "&appid=" + key)
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

	degrees, name := (weather.Main.Temperature-273.15)*1.8+32, weather.Name

	fmt.Printf("%s, %.0fF \n",
		name,
		degrees,
	)

}
