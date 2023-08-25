package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Tasks App")
}

func main() {
	http.HandleFunc("/hello", helloHandler)

	fmt.Println("Server started on :80")
	http.ListenAndServe(":80", nil)
}
