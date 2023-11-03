package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yoshikipom/go-webos"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("args are invalid. example) ./auth 192.168.1.2")
	}
	ip := os.Args[1]

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		NetDial: (&net.Dialer{
			Timeout: time.Second * 5,
		}).Dial,
	}

	tv, err := webos.NewTV(&dialer, ip)
	if err != nil {
		log.Fatalf("could not dial: %v", err)
	}
	defer tv.Close()

	go tv.MessageHandler()

	key, err := tv.AuthorisePrompt()
	if err != nil {
		log.Fatalf("could not authorise using prompt: %v", err)
	}

	fmt.Println(key)
	tv.Notification("📺👌")
}
