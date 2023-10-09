package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	logfile, _ := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer logfile.Close()
	log.SetOutput(logfile)

	startDb()

	http.HandleFunc("/replace", replaceHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/vclock", vclockHandler)
	http.HandleFunc("/ws", wsHandler)

	port := ":" + os.Args[1]
	fmt.Printf("Listenning on localhost%s\n", port)

	http.ListenAndServe(port, nil)
}
