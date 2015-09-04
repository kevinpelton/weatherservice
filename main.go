package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
)

//------------------------------------
// TYPES
//------------------------------------

// WeatherResponse
type WeatherResponse struct {
	Location    string
	Temperature float64
}

//Weather is
type Weather struct {
	Coord struct {
		Longitude float64 `json:"lon"`
		Latitude  float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
	} `json:"main"`
	Wind struct {
		Speed     float64 `json:"speed"`
		Direction int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		Percent int `json:"all"`
	} `json:"clouds"`
	Rain struct {
		LastThreeHours int `json:"3h"`
	} `json:"rain"`
	Snow struct {
		LastThreeHours int `json:"3h"`
	} `json:"snow"`
	Dt  int `json:"dt"`
	Sys struct {
		Type        int     `json:"type"`
		ID          int     `json:"id"`
		Message     float64 `json:"message"`
		CountryCode string  `json:"country"`
		Sunrise     int     `json:"sunrise"`
		Sunset      int     `json:"sunset"`
	} `json:"sys"`
	CityID   int    `json:"id"`
	CityName string `json:"name"`
	Cod      int    `json:"cod"`
}

//------------------------------------
// MAIN
//------------------------------------

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/weather/#location", func(w rest.ResponseWriter, req *rest.Request) {
			weather, err := GetWeather(req.PathParam("location"))
			if err != nil {
				rest.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//w.WriteJson(&weather)
			//
			w.WriteJson(
				&WeatherResponse{
					Location:    fmt.Sprintf("%s", weather.CityName),
					Temperature: KelvinToFarenheit(weather.Main.Temp),
				},
			)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

//------------------------------------
// FUNCS
//------------------------------------

// getContent is a generic URL get function
func getContent(url string) ([]byte, error) {
	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Send the request via a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// Defer the closing of the body
	defer resp.Body.Close()
	// Read the content into a byte array
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// GetWeather calls the openweathermap API
func GetWeather(location string) (*Weather, error) {

	content, err := getContent(
		fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s", location))
	if err != nil {
		return nil, err
	}
	// Fill the record with the data from the JSON
	var weather Weather
	err = json.Unmarshal(content, &weather)
	if err != nil {
		return nil, err
	}
	return &weather, err
}

// KelvinToFarenheit calculates fahrenheit values
func KelvinToFarenheit(temp float64) float64 {
	return (temp-273.15)*1.8000 + 32.00
}
