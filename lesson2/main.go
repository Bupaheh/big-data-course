package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func replace(w http.ResponseWriter, r *http.Request) {
	var buffer [1000]byte
	n, err := r.Body.Read(buffer[:])

	if err != nil && err != io.EOF {
		fmt.Println(n)
		fmt.Printf("err: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	f, _ := os.Create("cache")
	defer f.Close()
	f.Write(buffer[:n])
}

func get(w http.ResponseWriter, r *http.Request) {
	buffer, err := os.ReadFile("cache")

	if err != nil {
		fmt.Printf("err: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buffer)
}

func main() {
	http.HandleFunc("/replace", replace)
	http.HandleFunc("/get", get)

	http.ListenAndServe(":3333", nil)
}
