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

package aprsfi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Entries struct {
	Class    string
	Name     string
	Type     string
	Time     string
	Lasttime string
	Lat      string
	Lng      string
	Course   int16
	Speed    float32
	Symbol   string
	Srccall  string
	Dstcall  string
	Mice_msg string
	Comment  string
	Path     string
}

type Callsign struct {
	Command string
	Result  string
	What    string
	Found   int16
	Entries []Entries
}

func GetCallsign(call string, apiKey string) (Callsign, error) {
	var data Callsign

	resp, err := http.Get("http://api.aprs.fi/api/get?what=loc&apikey=" + apiKey + "&format=json&name=" + call)
	if err != nil {
		return data, errors.New("Error getting callsign data")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, errors.New("Error parsing callsign data")
	}

	json.Unmarshal([]byte(string(body)), &data)

	if data.Found < 1 {
		return data, errors.New("Callsign not found")
	}

	return data, nil
}
