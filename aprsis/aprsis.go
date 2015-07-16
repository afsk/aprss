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

package aprsis

import (
	"net"
	"strconv"
	"strings"
	"time"
)

func GeneratePassword(callsign string) string {
	parts := strings.Split(callsign, "-")
	simpleCallsign := strings.ToUpper(parts[0])
	length := len(simpleCallsign)
	if length < 1 {
		panic("Invalid callsign")
	}
	hash := 0x73e2
	for i := 0; i < length; {
		hash = hash ^ int(simpleCallsign[i])<<8
		i++
		if i < length {
			hash = hash ^ int(simpleCallsign[i])
			i++
		}
	}
	hash = hash & 0x7fff
	result := strconv.FormatInt(int64(hash), 10)

	return result
}

func SendPacket(packet string, conn net.Conn) {
	time.Sleep(2 * time.Second)
	conn.Write([]byte(packet + "\r\n"))
}
