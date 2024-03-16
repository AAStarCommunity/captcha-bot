package main

import (
	"fmt"
	"net/http"

	"github.com/assimon/captcha-bot/bootstrap"
)

func main() {
	go func() {
		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "OK")
		})
		http.ListenAndServe(":8080", nil)
	}()

	bootstrap.Start()
}
