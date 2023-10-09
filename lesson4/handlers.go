package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

//go:embed index.html
var content []byte

var localSource string = "Balay"
var localClock uint64 = 0

func replaceHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler. Replace request")

	var buffer [1000]byte
	n, err := r.Body.Read(buffer[:])

	if err != nil && err != io.EOF {
		fmt.Println(n)
		fmt.Printf("err: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Replace body: %s\n", string(buffer[:n]))

	localClock++

	transactions <- Transaction{
		Source:  localSource,
		Id:      localClock,
		Payload: string(buffer[:n]),
	}

	w.WriteHeader(http.StatusOK)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler. Get request")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(snap))
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler. Test request")

	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

func vclockHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler. Vclock request")

	bytes, _ := json.Marshal(vclock)

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler. Ws request")

	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		OriginPatterns:     []string{"*"},
	})

	if err != nil {
		log.Println("Handler. Failed to accept wx")
		return
	}

	defer c.Close(websocket.StatusGoingAway, "")

	transactionsSent := 0

	for {
		time.Sleep(200 * time.Millisecond)

		for transactionsSent < len(journal) {
			transaction := journal[transactionsSent]

			err = wsjson.Write(r.Context(), c, transaction)

			if err != nil {
				log.Println("Handler. Failed to write into ws")
				return
			}

			transactionsSent++
		}
	}
}
