package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func replace(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler. Replace request")

	var buffer [1000]byte
	n, err := r.Body.Read(buffer[:])

	if err != nil && err != io.EOF {
		fmt.Println(n)
		fmt.Printf("err: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	queries <- "replace " + string(buffer[:n])
	w.WriteHeader(http.StatusOK)
}

func get(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler. Get request")

	queries <- "get"
	str := <-getResults

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(str))
}

func main() {
	logfile, _ := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer logfile.Close()
	log.SetOutput(logfile)

	startDb()

	http.HandleFunc("/replace", replace)
	http.HandleFunc("/get", get)

	port := ":3333"
	fmt.Printf("Listenning on localhost%s\n", port)

	http.ListenAndServe(port, nil)
}
