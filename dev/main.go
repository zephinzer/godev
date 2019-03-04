package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("hello world")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("request received")
		w.Write([]byte("hello world"))
	})
	http.ListenAndServe("0.0.0.0:1337", nil)
}
