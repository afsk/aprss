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

package message

import (
	"errors"
	"regexp"
	"strings"
)

type AprsMessage struct {
	From    string
	To      string
	Path    string
	Message string
	Code    string
}

func Parse(data string) (*AprsMessage, error) {
	regex := "^([A-Z0-9\\-]+)\\>(.+)::([A-Z0-9\\-\\s]+):([^\\|^\\{^\\~]+)\\{?(.{0,5})"
	re := regexp.MustCompile(regex)
	all := re.FindAllStringSubmatch(data, -1)
	if len(all) > 0 {
		matched := all[0]
		if len(matched) == 6 {
			result := new(AprsMessage)
			result.From = strings.TrimSpace(matched[1])
			result.Path = strings.TrimSpace(matched[2])
			result.To = strings.TrimSpace(matched[3])
			result.Message = strings.TrimSpace(matched[4])
			result.Code = strings.TrimSpace(matched[5])

			return result, nil
		}
	}
	return nil, errors.New("Data not properly parsed")
}

func (m *AprsMessage) GetAck() (string, error) {
	if len(m.Code) <= 0 {
		return "", errors.New("Message does not have a code")
	}
	if len(m.Code) > 5 {
		return "", errors.New("Code is longer than 5 chars")
	}
	if len(m.From) <= 0 {
		return "", errors.New("Empty From field")
	}
	if len(m.To) <= 0 {
		return "", errors.New("Empty To field")
	}

	ack := m.To + ">TCPIP*::" + padRight(m.From, " ", 9) + ":ack" + m.Code
	return ack, nil
}

func (m *AprsMessage) GetData() (string, error) {
	if len(m.Code) > 5 {
		return "", errors.New("Code is longer than 5 chars")
	}
	if len(m.From) <= 0 {
		return "", errors.New("Empty From field")
	}
	if len(m.To) <= 0 {
		return "", errors.New("Empty To field")
	}
	if len(m.Message) <= 0 {
		return "", errors.New("Empty Message field")
	}

	data := m.From + ">TCPIP*::" + padRight(m.To, " ", 9) + ":" + m.Message
	if len(m.Code) > 0 {
		data += "{" + m.Code
	}

	return data, nil
}

func padRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}
