package main

import (
	"github.com/valyala/fasthttp"
	"log"
)

func main() {
	status, body, err := fasthttp.Get(nil, "https://www.baidu.com")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("status: %v, body: %s", status, string(body))
}
