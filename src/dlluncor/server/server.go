package server

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func Serve() {
	http.HandleFunc("/hi", helloHandler)
	http.ListenAndServe("localhost:8080", nil)
}
