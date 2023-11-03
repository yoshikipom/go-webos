package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/gorilla/websocket"

	webos "github.com/yoshikipom/go-webos"
)

const (
	youtubeAppID = "youtube.leanback.v4"
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

func sleepBeforeLoginClick(currentAppIsYoutube bool) {
	var sleep_time time.Duration
	if currentAppIsYoutube {
		sleep_time = 2 * time.Second
	} else {
		sleep_time = 11 * time.Second
	}
	time.Sleep(sleep_time)
}

func main() {
	if len(os.Args) != 4 {
		log.Fatalln("args are invalid. example) ./auth 192.168.1.2 abcde(client-id) https://www.youtube.com/watch?v=something")
	}
	ipAddress := os.Args[1]
	clientId := os.Args[2]
	videoURL := os.Args[3]
	fmt.Printf("%+v\n", clientId)

	tv := connectToTV(ipAddress, clientId)
	defer tv.Close()

	a, err := tv.CurrentApp()
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
	}
	currentAppIsYoutube := a.AppID == youtubeAppID

	if !currentAppIsYoutube {
		tv.OpenApp("youtube.leanback.v4")
	}

	_, err = tv.Command(webos.SystemLauncherLaunchCommand, webos.Payload{"id": youtubeAppID, "params": map[string]interface{}{"contentTarget": videoURL}})
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
	}

	sleepBeforeLoginClick(currentAppIsYoutube)

	tv.KeyEnter()
}
