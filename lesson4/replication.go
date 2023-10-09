package main

import (
	"context"
	"fmt"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

func websocketSession(peer string) {
	var ctx = context.Background()

	c, _, err := websocket.Dial(ctx, fmt.Sprintf("ws://%s/ws", peer), nil)

	if err != nil {
		log.Printf("Dial. Failed to connect to %s\n", peer)
		return
	}

	defer c.Close(websocket.StatusGoingAway, "")

	for {
		var transaction Transaction

		err = wsjson.Read(ctx, c, &transaction)

		if err != nil {
			log.Printf("Dial. Failed to receive transaction from %s\n", peer)
			return
		}

		transactions <- transaction
	}
}

func websocketClient(peer string) {
	for {
		websocketSession(peer)
		time.Sleep(5 * time.Second)
	}
}
