package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":80", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello Jérôme Le Vatant !")
}
