package main

import (
	"net/http"

	"./manager"
)

func main() {
	http.HandleFunc("/", manager.ContainerHandler)
	http.ListenAndServe(":8081", nil)
}
