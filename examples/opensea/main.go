package main

import (
	"encoding/json"
	"github.com/ntop001/phx2"
	"log"
	"net/url"
)

func main() {
	// Create a url.URL to connect to. `ws://` is non-encrypted websocket.
	urlStr := "wss://stream.openseabeta.com/socket?token=ae36b1b4131c421e8c84088ad48abe9b&vsn=2.0.0"
	endPoint, _ := url.Parse(urlStr)
	log.Println("Connecting to", endPoint)

	// Create a new phx.Socket
	socket := phx2.NewSocket(endPoint)
	socket.Logger = phx2.NewSimpleLogger(phx2.LogInfo)

	// Wait for the socket to connect before continuing. If it's not able to, it will keep
	// retrying forever.
	cont := make(chan bool)
	socket.OnOpen(func() {
		cont <- true
	})
	log.Println("connection opened.")

	// Tell the socket to connect (or start retrying until it can connect)
	err := socket.Connect()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connection finished.")


	// Wait for the connection
	<-cont
	log.Println("start create channel:")

	// Create a phx.Channel to connect to the default 'room:lobby' channel with no params
	// (alternatively params could be a `map[string]string`)
	channel := socket.Channel("collection:*", nil)
	log.Println("create channel:", channel.Topic())

	// Join the channel. A phx.Push is returned which can be used to bind to replies and errors
	join, err := channel.Join()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("join channel:", join.Event)

	// Listen for a response and only continue once we know we're joined
	join.Receive("ok", func(response interface{}) {
		log.Println("Joined channel:", channel.Topic(), response)
		cont <- true
	})
	join.Receive("error", func(response interface{}) {
		log.Println("Join error", response)
	})

	// wait to be joined
	<-cont

	// get events
	channel.On("item_listed", func(payload interface{}) {
		data, _ := json.Marshal(payload)
		log.Println("get events:", string(data))
	})

	// Now we will block forever, hit ctrl+c to exit
	select {}
}
