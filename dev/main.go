package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "1337"
	}
	fmt.Printf("hello world from port %s\n", port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("request received")
		w.Write([]byte("hello world"))
	})
	http.ListenAndServe(
		fmt.Sprintf("0.0.0.0:%s", port),
		nil,
	)
}
