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

package main

import (
	"bufio"
	"fmt"
	"github.com/yo3igc/aprss/aprsis"
	"github.com/yo3igc/aprss/message"
	"net"
	"strings"
	"io"
)

const serviceCallsign string = ""

var servicePassword string

var conn net.Conn

func main() {
	if len(serviceCallsign) < 1 {
		panic("Please set your callsign!")
	}

	servicePassword = aprsis.GeneratePassword(serviceCallsign)

	if !openConnection() {
		panic("Error opening connection")
	}
	defer conn.Close()

	login := fmt.Sprintf("user %s pass %s vers aprss 0.1 filter t/m g/%s u/%s\r\n", serviceCallsign, servicePassword, serviceCallsign, serviceCallsign)
	_, err := conn.Write([]byte(login))
	if err != nil {
		panic(err.Error())
	}

	for {
		line, err := readLine()
		if err != nil {
			if err == io.EOF {
				openConnection()
				continue
			} else {
				panic(err.Error())
			}
		}
		if len(line) < 1 || line[0] == '#' {
			continue
		}
		go handleLine(line)
	}
}

func openConnection() bool {
	var err error
	conn, err = net.Dial("tcp", "euro.aprs2.net:14580")
	if err != nil {
		return false
	}

	return true
}

func readLine() (string, error) {
	line, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}

	return line, nil
}

func handleLine(rawData string) {
	packet, err := message.Parse(rawData)
	if err != nil || packet.To != serviceCallsign {
		return
	}

	fmt.Print(rawData)

	if packet.To == serviceCallsign {
		sendAck(packet)
	}

	if strings.ToLower(packet.Message) == "ping" {
		retPacket := new(message.AprsMessage)
		retPacket.From = serviceCallsign
		retPacket.To = packet.From
		retPacket.Message = "pong"
		retRawData, err := retPacket.GetData()
		if err == nil {
			aprsis.SendPacket(retRawData, conn)
			fmt.Println(retRawData)
		}
	}
}

func sendAck(packet *message.AprsMessage) {
	ack, err := packet.GetAck()
	if err == nil {
		aprsis.SendPacket(ack, conn)
		fmt.Println(ack)
	}
}
