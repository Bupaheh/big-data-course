package main

import (
	jsonpatch "github.com/evanphx/json-patch/v5"
	"log"
)

type Transaction struct {
	Source  string
	Id      uint64
	Payload string
}

var snap = "{}"
var transactions = make(chan Transaction, 1337)
var journal = make([]Transaction, 0)
var vclock = make(map[string]uint64)

func startDb() {
	go websocketReplication()
	go transactionManager()
}

func transactionManager() {
	for {
		transaction := <-transactions
		log.Printf("Database. Query: %+v\n", transaction)

		localSourceClock, contains := vclock[transaction.Source]

		if contains && localSourceClock >= transaction.Id {
			log.Println("Database. Already handled transaction")
			return
		}

		patch, err := jsonpatch.DecodePatch([]byte(transaction.Payload))

		if err != nil {
			log.Println("Database. Incorrect patch")
			return
		}

		newSnap, err := patch.Apply([]byte(snap))

		if err != nil {
			log.Println("Database. Failed to apply patch")
			return
		}

		vclock[transaction.Source] = transaction.Id

		snap = string(newSnap)
		journal = append(journal, transaction)

		log.Printf("Database. Patch applied: %s\n", snap)
	}
}

var peers = []string{"localhost:8081"}

func websocketReplication() {
	for _, peer := range peers {
		go websocketClient(peer)
	}
}
