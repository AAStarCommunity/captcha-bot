package main

import (
	"fmt"
	"github.com/assimon/captcha-bot/bootstrap"
	"net/http"
)

func main() {
	go func() {
		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			fmt.Fprintf(w, "OK")
		})
		http.ListenAndServe(":80", nil)
	}()

	bootstrap.Start()
}
