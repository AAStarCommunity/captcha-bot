package main

import (
	"fmt"
	"github.com/assimon/captcha-bot/bootstrap"
	"io"
	"net/http"
	"time"
)

func main() {
	go func() {
		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			fmt.Fprintf(w, "OK")
		})
		http.ListenAndServe(":80", nil)
	}()

	go func() {
		t := time.NewTimer(time.Second * 10)
		for {
			if resp, err := http.Get("https://captcha-bot-zg9p.onrender.com/healthz"); err != nil {
				fmt.Println(fmt.Sprintf("err:%v", err))
			} else {
				if b, err := io.ReadAll(resp.Body); err == nil {
					fmt.Println("health: " + string(b))
				}
			}
			<-t.C
		}
	}()

	bootstrap.Start()
}
