// The MIT License (MIT)

// Copyright (c) 2015 Alex Mirea <yo3igc@gmail.com>

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package weather

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Coordinates struct {
	Lon float32
	Lat float32
}

type Weather struct {
	Id          int32
	Main        string
	Description string
	Icon        string
}

type Main struct {
	Temp     float32
	Pressure int32
	Humidity int32
	Temp_min float32
	Temp_max float32
}

type Wind struct {
	Speed float32
	Deg   int32
}

type Clouds struct {
	All int32
}

type Sys struct {
	Type    int32
	Id      int32
	Message float32
	Country string
	Sunrise int32
	Sunset  int32
}

type WeatherResponse struct {
	Coord      Coordinates
	Weather    []Weather
	Base       string
	Main       Main
	Visibility int32
	Wind       Wind
	Clouds     Clouds
	Dt         int32
	Sys        Sys
	Id         int32
	Name       string
	Cod        int16
	Message    string
}

func GetByCity(city string, apiKey string) (WeatherResponse, error) {
	return get("q="+city, apiKey)
}

func GetByLocation(lat string, lon string, apiKey string) (WeatherResponse, error) {
	return get("lat="+lat+"&lon="+lon, apiKey)
}

func get(query string, apiKey string) (WeatherResponse, error) {
	var weather WeatherResponse

	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiKey + "&units=metric&" + query)
	if err != nil {
		return weather, errors.New("Error getting weather data")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return weather, errors.New("Error parsing weather data")
	}

	json.Unmarshal([]byte(string(body)), &weather)

	if weather.Cod != 200 {
		return weather, errors.New(weather.Message)
	}

	return weather, nil
}
