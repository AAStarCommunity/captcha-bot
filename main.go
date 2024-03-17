package main

import (
	"fmt"
	"github.com/assimon/captcha-bot/bootstrap"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

var counter = 0

func main() {
	go func() {
		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			counter++
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("counter", strconv.Itoa(counter))
			fmt.Fprintf(w, "OK")
		})
		_ = http.ListenAndServe(":80", nil)
	}()

	go func() {
		t := time.Tick(time.Second * 10)
		for {
			if resp, err := http.Get("https://captcha-bot-zg9p.onrender.com/healthz"); err != nil {
				log.Default().Printf("error: " + err.Error())
			} else {
				if b, err := io.ReadAll(resp.Body); err == nil {
					log.Default().Printf("health: " + string(b))
				}
			}
			<-t
		}
	}()

	bootstrap.Start()
}
