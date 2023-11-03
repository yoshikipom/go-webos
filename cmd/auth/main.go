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

func fetchClientKey(ip string) {
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
		// the TV uses a self-signed certificate
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		NetDial:         (&net.Dialer{Timeout: time.Second * 5}).Dial,
	}

	tv, err := webos.NewTV(&dialer, ip)
	if err != nil {
		log.Fatalf("could not dial TV: %v", err)
	}
	defer tv.Close()

	// the MessageHandler must be started to read responses from the TV
	go tv.MessageHandler()

	// AuthorisePrompt shows the authorisation prompt on the TV screen
	key, err := tv.AuthorisePrompt()
	if err != nil {
		log.Fatalf("could not authorise using prompt: %v", err)
	}

	// the key returned can be used for future request to the TV using the
	// AuthoriseClientKey(<key>) method, instead of AuthorisePrompt()
	fmt.Print(key)

	// see commands.go for available methods
	tv.Notification("📺👌")
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("args are invalid. example) ./auth 192.168.1.2")
	}
	ip := os.Args[1]

	fetchClientKey(ip)
}
