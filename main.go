package main

import (
	"fmt"
	"net/http"
)

func main() {

	server := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(basicHandler),
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func basicHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello World\n"))
}
