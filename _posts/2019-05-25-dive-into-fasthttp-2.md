---
layout: post
title:  "fasthttp高性能之道（二）"
date:   2019-05-25 22:00:00 +0800
categories: "GoodCode"
---
在[fasthttp高性能知道（一）]()中我们简要介绍了fasthttp项目的特点，以及Server端的一些实现思路，本篇将会把关注点从Server端移到Client端，分析一下fasthttp在Client端的实现又有哪些比较有意思的地方。

fasthttp包含四种Client，分别是：

* Client
* HostClient
* PipelineClient
* LBClient

其中，默认的client即为一个Client对象，对于这几种Client的差异和用法，我们下面会逐一介绍。

首先，和介绍Server一样，我们也来看两个例子：

```go
package main

import (
	"log"

	"github.com/valyala/fasthttp"
)

func main() {
	status, body, err := fasthttp.Get(nil, "https://www.baidu.com")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("status: %v, body: %s", status, string(body))
}
```

这是一个使用默认Client的例子，这里我们直接调用`fasthttp.Get`就可以发起HTTP请求了，Get的第一个参数是保存body的byte数组切片，如果你希望重用这个对象，可以传递一个body数组的切片进去，这样可以减少GC。

如果希望更细粒度的控制各种参数，如超时、连接数限制等，可以看下面这个例子：

```go
import (
	"log"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	// HTTPClient global http client object
	client *fasthttp.Client = &fasthttp.Client{
		MaxConnsPerHost: 16384, // MaxConnsPerHost  default is 512, increase to 16384
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    5 * time.Second,
	}
)

func main() {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://127.0.0.1:29898/api/v1/report")
	req.Header.SetMethod("POST")
	req.Header.SetContentType("text/plain")
	req.SetBody([]byte("hello world"))

	resp := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseResponse(resp)
	defer fasthttp.ReleaseRequest(req)

	if err := client.Do(req, resp); err != nil {
		log.Fatal(err)
	}

	log.Println(resp)
}
```

上面的代码有两个地方需要注意：
1. MaxConnsPerHost是一个限流的参数，保证对一个Host最大的打开连接数，如果超过这个数字，则会直接拒绝，这里默认值是512，但如果你打算用来做压测之类的事情，需要增加这个值，比如这里我就增加到了16384。
1. AcquireRequest和AcquireResponse分别从requestPool和responsePool中获取对象，所以用完得记得调用ReleaseRequest和ReleaseResponse把他们还回去，另外需要注意，由于他们是从对象池中获取的，当release之后他们的值可能会被覆盖，相关的处理一定要在release之前进行。

## HostClient

我们先从HostClient来分析，它也是Client的基础。

```go
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
	body = make([]byte, 4096)
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Missing url")
	}
	urls := os.Args[1:]
	for _, url := range urls {
		statusCode, body, err := client.Get(body, url)
		if err != nil {
			log.Fatalf("Error when loading page %s through local proxy: %s", url, err)
		}
		if statusCode != fasthttp.StatusOK {
			log.Fatalf("Unexpected status code: %d. Expecting %d", statusCode, fasthttp.StatusOK)
		}
		log.Printf("body: %s\n", string(body))
	}
}
```

我们先通过上面的例子来了解HostClient一些有趣的特点，将上面的代码保存为hostclient.go，执行`go build hostclient.go`完成编译，然后找两个目录，分别执行`python3 -m http.server 19898`和`python3 -m http.server 29898`。

完成这些之后，我们执行`./hostclient http://localhost:19898/a.txt http://localhost:19898/b.txt`，观察两个python进程的请求日志，会发现第一个请求发送到了监听19898端口的服务，而第二个请求则发给了监听29898端口的服务。

接下来我们把请求地址改成别的，如`./hostclient https://www.baidu.com http://www.jd.com`，发现请求仍然是发送给了两个Python进程。

现在，让我们来总结一下：
* Addr只有一个地址，且请求的URL就是在这个地址上的话，与其它语言的HTTP client没区别
* Addr有多个地址，无论URL请求的是哪个，都会在这多个地址上轮转，即一定程度的load balance，所以可以基于此实现反向代理功能
* 请求URL与Addr不同的时候，Addr扮演了正向代理服务的角色

在深入到HostClient的实现内部之前，我们先来梳理一下HTTP Client的基本思路。
我们知道HTTP(s)协议是构建在TCP之上的，作为一个Client，如果我们请求的地址是固定的，我们一般希望保持一个长连接，然后在这个连接之上发送HTTP报文。那么完成一次HTTP请求需要哪些工作呢？简单罗列一下，它一般包括：

1. DNS请求，将目标域名翻译成IP地址
1. 建立一个到目标IP:PORT的TCP连接
1. 通过TCP连接发送HTTP请求报文
1. 接收HTTP响应报文
1. 重复步骤3~4
1. 结束请求，关闭连接

让那个我们先停下来思考一下，要实现一个高性能的HTTP Client，我们需要注意哪些问题呢
首先，DNS请求不能太过频繁，如果每次建立连接都要进行DNS解析的话，对DNS服务器的冲击和对请求建连的开销就有点大了。
其次，TCP连接是很昂贵的，我们除了要保证尽可能地复用之外，还需要在连接不需要时，及早将其清理掉。
第三，HTTP的请求和响应是很频繁的，对于Request和Response对象，每次都分配显然是太浪费了，对象池技术在这里非常有用。
第四，如果一个Client同时建立了海量到同一个服务器的连接，那对服务器的压力是很大的，我们应当做一些限制和防范。

## Reference

* [fasthttp client internals](https://youtu.be/fg3JPUswiek)
