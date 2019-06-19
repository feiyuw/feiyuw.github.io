package main

import (
	"github.com/valyala/fasthttp"
	"log"
	"os"
)

var (
	client = &fasthttp.HostClient{
		Addr: "localhost:19898,localhost:29898",
	}
	body = make([]byte, 1)
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Missing url")
	}
	urls := os.Args[1:]
	for _, url := range urls {
		statusCode, body, err := client.Get(body, url)
		if err != nil {
			log.Fatalf("Error when loading foobar page through local proxy: %s", err)
		}
		if statusCode != fasthttp.StatusOK {
			log.Fatalf("Unexpected status code: %d. Expecting %d", statusCode, fasthttp.StatusOK)
		}
		log.Printf("body: %s\n", string(body))
	}
}
