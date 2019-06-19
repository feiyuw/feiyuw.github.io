package main

import (
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

var (
	lbc fasthttp.LBClient
)

func main() {
	servers := []string{
		"127.0.0.1:8888",
		"127.0.0.1:9999",
	}

	for _, addr := range servers {
		c := &fasthttp.HostClient{
			Addr: addr,
		}
		lbc.Clients = append(lbc.Clients, c)
	}

	var req fasthttp.Request
	var resp fasthttp.Response
	for i := 0; i < 10; i++ {
		url := fmt.Sprintf("http://abcedfg/foo/bar/%d", i)
		req.SetRequestURI(url)
		if err := lbc.Do(&req, &resp); err != nil {
			log.Printf("Error when sending request: %s", err)
			continue
		}
		if resp.StatusCode() != fasthttp.StatusOK {
			log.Printf("unexpected status code: %d. Expecting %d", resp.StatusCode(), fasthttp.StatusOK)
			continue
		}

		log.Println(resp.Body())
	}
}
