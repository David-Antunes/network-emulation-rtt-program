package main

import (
	"bytes"
	"encoding/json"
	"net"
	"strings"
	"time"

	api "github.com/David-Antunes/network-emulation-proxy/api"
)

func main() {
	listenAddr, err := net.ResolveUDPAddr("udp4", ":8000")
	if err != nil {
		panic(err)
	}
	port, err := net.ListenUDP("udp4", listenAddr)
	if err != nil {
		panic(err)
	}

	ief, err := net.InterfaceByName("eth0")
	if err != nil {
		panic(err)
	}
	addrs, err := ief.Addrs()
	if err != nil {
		panic(err)
	}
	splitAddr := strings.Split(addrs[0].String(), ".")

	if len(splitAddr) != 4 {
		panic("something went wrong with Ip address")
	}

	broadcastIp := splitAddr[0] + "." + splitAddr[1] + "." + splitAddr[2] + ".255:8000"

	conn, err := net.Dial("udp4", broadcastIp)

	if err != nil {
		panic(err)
	}

	for {
		buf := make([]byte, 1024)
		size, err := port.Read(buf)

		if err != nil {
			panic(err)
		}

		buf = buf[:size]

		resp := &api.RTTRequest{}
		d := json.NewDecoder(bytes.NewReader(buf))
		err = d.Decode(resp)
		if err != nil {
			panic(err)
		}

		resp.ReceiveTime = time.Now()
		resp.TransmitTime = time.Now()

		req, err := json.Marshal(resp)

		if err != nil {
			panic(err)
		}

		conn.Write(req)
	}

}
