package main

import (
	"crypto/tls"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	webos "github.com/yoshikipom/go-webos"
)

var dialer = websocket.Dialer{
	HandshakeTimeout: 10 * time.Second,
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
	NetDial: (&net.Dialer{
		Timeout: time.Second * 5,
	}).Dial,
}

func connectToTV(ipAddress string, clientId string) *webos.TV {
	tv, err := webos.NewTV(&dialer, ipAddress)
	if err != nil {
		log.Fatalf("could not dial: %v", err)
	}

	go tv.MessageHandler()

	if err = tv.AuthoriseClientKey(clientId); err != nil {
		log.Fatalf("could not authoise using client key: %v", err)
	}

	return tv
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("args are invalid. example) ./auth 192.168.1.2 abcde(client-id) message")
	}
	ipAddress := os.Args[1]
	clientId := os.Args[2]
	target := os.Args[3]

	tv := connectToTV(ipAddress, clientId)
	defer tv.Close()

	num, _ := strconv.Atoi(target)
	if num > 20 {
		log.Fatalf("Volume: %d might be too loud. Skip to change volume", num)
	}
	tv.SetVolume(num)
}
