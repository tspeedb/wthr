package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/guptarohit/asciigraph"
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

type WeatherResponse struct {
	Cod     string        `json:"cod"`
	Message int           `json:"message"`
	Cnt     int           `json:"cnt"`
	List    []WeatherData `json:"list"`
}

type WeatherData struct {
	Dt      int64              `json:"dt"`
	Main    MainInfo           `json:"main"`
	Weather []WeatherCondition `json:"weather"`
}

type MainInfo struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
}

type WeatherCondition struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
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
	defaultPlace := "Boston"

	// Use Geocoder API to get lat and long of place in args if longer than 2

	if len(os.Args) >= 2 {
		place := os.Args[1:]

		res, err := http.Get("http://api.openweathermap.org/geo/1.0/direct?q=" + strings.Join(place, "_") + "&limit=5&appid=" + key)
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

	res, err := http.Get("http://api.openweathermap.org/data/2.5/forecast?lat=" + lat + "&lon=" + lon + "&appid=" + key)
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

	var response WeatherResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	var temps []float64
	for _, weather := range response.List {
		temps = append(temps, (weather.Main.Temp-273.15)*1.8+32)
	}

	graph := asciigraph.Plot(temps)

	if len(os.Args) >= 2 {
		place := os.Args[1:]
		fmt.Println(strings.Join(place, "_"), "\n", graph)
	} else {
		fmt.Println(defaultPlace, "\n", graph)
	}

}
