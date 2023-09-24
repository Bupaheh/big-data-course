package main

import (
	"container/list"
	"log"
	"strings"
	"sync"
	"time"
)

var queries = make(chan string)
var getResults = make(chan string)

var db struct {
	mu   sync.Mutex
	data string
}
var snapshot string
var dbLog = list.New()

func startDb() {
	go transactionManager()
	go snapshotManager()
}

func transactionManager() {
	for {
		str := <-queries

		log.Printf("Database. Query: %s\n", str)

		db.mu.Lock()

		switch {
		case str == "get":
			getResults <- db.data
		case strings.HasPrefix(str, "replace "):
			body, _ := strings.CutPrefix(str, "replace ")
			db.data = body
		default:
			log.Println("Incorrect query")
		}

		db.mu.Unlock()

		dbLog.PushBack(str)
	}
}

func snapshotManager() {
	for {
		db.mu.Lock()
		snapshot = db.data
		db.mu.Unlock()

		log.Printf("Database. Current snapshot: %s\n", snapshot)
		time.Sleep(time.Minute)
	}
}
