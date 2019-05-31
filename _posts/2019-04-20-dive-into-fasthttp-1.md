---
layout: post
title:  "fasthttp高性能之道（一）"
date:   2019-04-20 22:00:00 +0800
categories: "GoodCode"
---
本篇是fasthttp高性能知道的第一篇，我计划用三篇博客来分析一下fasthttp这个库，也帮助自己更好地了解高并发HTTP服务器的设计思路。

2017年第一次接触fasthttp，用它构建了一个http服务，替换之前跑的python服务，在使用过程中对其设计思想和优化思路非常欣赏。2018年，在做压测工具的时候，采用它构建了http的客户端，同样充满惊喜，尤其是内存占用控制的很好。使用中间断断续续地看过一些它的代码，这里做一个总结，代码不多，当是学习一下高并发http框架的实现思路。

## 总览

[fasthttp](https://github.com/valyala/fasthttp)是一个纯golang编写的，以高并发、高性能为目的的HTTP library，包括客户端和服务端。相比`net/http`，fasthttp在大部分情况下都拥有更好的性能以及更低的内存占用，但是相比`net/http`，它不支持HTTP/2（对于fasthttp的HTTP/2支持，见 https://github.com/dgrr/http2 项目）。

## 依赖

fasthttp的依赖项很少，只有4个直接依赖项，他们是：
* github.com/klauspost/compress
* github.com/valyala/bytebufferpool
* github.com/valyala/tcplisten
* golang.org/x/net

这里面`klausport/compress`模块是标准库`compress`的替代品，主要目的是提高压缩性能和减少压缩时候的内存占用。

`bytebufferpool`维护了一个数据对象池，比如当接收到Request的时候，就会从该对象池中获取一个数据对象填充，用完再还回去，这种方式可以显著减少GC，关于它，我们后面会专门介绍。

`tcplisten`则是为了在启动HTTP Server时，添加对几个TCP options的支持，它们是：
* SO_REUSEPORT          多个进程可以绑定到同一个端口，内核层面进行负载均衡
* TCP_DEFER_ACCEPT      三次握手之后，服务端不马上accept，必须等到客户端数据到来之后才accept，可以减少惊群效应
* TCP_FASTOPEN          一个简化三次握手手续的扩展，提高两端点间连接的打开速度

具体代码见`fasthttp.reuseport`。因为这些设置非常常用，fasthttp就内置支持了创建一个建立在reuseport的TCP连接上的HTTP server，见下面的示例代码。

## Server

从一个使用示例开始
```go
package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"log"
)

func defaultHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNoContent)
}

func onReportV1(ctx *fasthttp.RequestCtx) {
	content := ctx.PostBody()
	fields := make(map[string]interface{})
	if err := json.Unmarshal(content, &fields); err != nil {
		panic("invalid content")
	}
	log.Printf("%v", fields)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func main() {
	listener, err := reuseport.Listen("tcp4", "0.0.0.0:1234")
	if err != nil {
		log.Fatal(err)
	}
	router := fasthttprouter.New()
	router.POST("/v1/report", onReportV1)
	router.NotFound = defaultHandler
	fasthttp.Serve(listener, router.Handler)
}
```

这是一个非常简单的HTTP的例子，查看main函数，前3行用reuseport创建了一个TCP listener，bind到1234端口。4 ~ 6行，为了简化handler的编写，引入了fasthttprouter模块，并添加了到`POST /v1/report`接口的handler，即onReportV1，同时对于没有定义的handler，使用defaultHandler。第7行通过fasthttp.Serve方法启动HTTP服务。整个代码非常简洁明了。

而对于每一个handler，都只有唯一的一个`*RequestCtx`类型的参数，通过该参数可以获取请求数据，如`ctx.PostBody()`，设置返回状态码（见`ctx.SetStatusCode(fasthttp.StatusOK)`。相对于标准库`net/http`的request，response两个参数的形式，fasthttp作了一定的简化。

`fasthttp.Serve`方法只是对`Server.Serve`的封装，因此我们先关注一下`Server`的结构。

```go
type Server struct {
	noCopy noCopy

	// Handler for processing incoming requests.
	//
	// Take into account that no `panic` recovery is done by `fasthttp` (thus any `panic` will take down the entire server).
	// Instead the user should use `recover` to handle these situations.
	Handler RequestHandler

	// ErrorHandler for returning a response in case of an error while receiving or parsing the request.
	//
	// The following is a non-exhaustive list of errors that can be expected as argument:
	//   * io.EOF
	//   * io.ErrUnexpectedEOF
	//   * ErrGetOnly
	//   * ErrSmallBuffer
	//   * ErrBodyTooLarge
	//   * ErrBrokenChunks
	ErrorHandler func(ctx *RequestCtx, err error)

	// Server name for sending in response headers.
	//
	// Default server name is used if left blank.
	Name string

	// The maximum number of concurrent connections the server may serve.
	//
	// DefaultConcurrency is used if not set.
	Concurrency int

	// Whether to disable keep-alive connections.
	//
	// The server will close all the incoming connections after sending
	// the first response to client if this option is set to true.
	//
	// By default keep-alive connections are enabled.
	DisableKeepalive bool

	// Per-connection buffer size for requests' reading.
	// This also limits the maximum header size.
	//
	// Increase this buffer if your clients send multi-KB RequestURIs
	// and/or multi-KB headers (for example, BIG cookies).
	//
	// Default buffer size is used if not set.
	ReadBufferSize int

	// Per-connection buffer size for responses' writing.
	//
	// Default buffer size is used if not set.
	WriteBufferSize int

	// ReadTimeout is the amount of time allowed to read
	// the full request including body. The connection's read
	// deadline is reset when the connection opens, or for
	// keep-alive connections after the first byte has been read.
	//
	// By default request read timeout is unlimited.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset after the request handler
	// has returned.
	//
	// By default response write timeout is unlimited.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alive is enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used.
	IdleTimeout time.Duration

	// Maximum number of concurrent client connections allowed per IP.
	//
	// By default unlimited number of concurrent connections
	// may be established to the server from a single IP address.
	MaxConnsPerIP int

	// Maximum number of requests served per connection.
	//
	// The server closes connection after the last request.
	// 'Connection: close' header is added to the last response.
	//
	// By default unlimited number of requests may be served per connection.
	MaxRequestsPerConn int

	// MaxKeepaliveDuration is a no-op and only left here for backwards compatibility.
	// Deprecated: Use IdleTimeout instead.
	MaxKeepaliveDuration time.Duration

	// Whether to enable tcp keep-alive connections.
	//
	// Whether the operating system should send tcp keep-alive messages on the tcp connection.
	//
	// By default tcp keep-alive connections are disabled.
	TCPKeepalive bool

	// Period between tcp keep-alive messages.
	//
	// TCP keep-alive period is determined by operation system by default.
	TCPKeepalivePeriod time.Duration

	// Maximum request body size.
	//
	// The server rejects requests with bodies exceeding this limit.
	//
	// Request body size is limited by DefaultMaxRequestBodySize by default.
	MaxRequestBodySize int

	// Aggressively reduces memory usage at the cost of higher CPU usage
	// if set to true.
	//
	// Try enabling this option only if the server consumes too much memory
	// serving mostly idle keep-alive connections. This may reduce memory
	// usage by more than 50%.
	//
	// Aggressive memory usage reduction is disabled by default.
	ReduceMemoryUsage bool

	// Rejects all non-GET requests if set to true.
	//
	// This option is useful as anti-DoS protection for servers
	// accepting only GET requests. The request size is limited
	// by ReadBufferSize if GetOnly is set.
	//
	// Server accepts all the requests by default.
	GetOnly bool

	// Logs all errors, including the most frequent
	// 'connection reset by peer', 'broken pipe' and 'connection timeout'
	// errors. Such errors are common in production serving real-world
	// clients.
	//
	// By default the most frequent errors such as
	// 'connection reset by peer', 'broken pipe' and 'connection timeout'
	// are suppressed in order to limit output log traffic.
	LogAllErrors bool

	// Header names are passed as-is without normalization
	// if this option is set.
	//
	// Disabled header names' normalization may be useful only for proxying
	// incoming requests to other servers expecting case-sensitive
	// header names. See https://github.com/valyala/fasthttp/issues/57
	// for details.
	//
	// By default request and response header names are normalized, i.e.
	// The first letter and the first letters following dashes
	// are uppercased, while all the other letters are lowercased.
	// Examples:
	//
	//     * HOST -> Host
	//     * content-type -> Content-Type
	//     * cONTENT-lenGTH -> Content-Length
	DisableHeaderNamesNormalizing bool

	// SleepWhenConcurrencyLimitsExceeded is a duration to be slept of if
	// the concurrency limit in exceeded (default [when is 0]: don't sleep
	// and accept new connections immidiatelly).
	SleepWhenConcurrencyLimitsExceeded time.Duration

	// NoDefaultServerHeader, when set to true, causes the default Server header
	// to be excluded from the Response.
	//
	// The default Server header value is the value of the Name field or an
	// internal default value in its absence. With this option set to true,
	// the only time a Server header will be sent is if a non-zero length
	// value is explicitly provided during a request.
	NoDefaultServerHeader bool

	// NoDefaultContentType, when set to true, causes the default Content-Type
	// header to be excluded from the Response.
	//
	// The default Content-Type header value is the internal default value. When
	// set to true, the Content-Type will not be present.
	NoDefaultContentType bool

	// ConnState specifies an optional callback function that is
	// called when a client connection changes state. See the
	// ConnState type and associated constants for details.
	ConnState func(net.Conn, ConnState)

	// Logger, which is used by RequestCtx.Logger().
	//
	// By default standard logger from log package is used.
	Logger Logger

	// KeepHijackedConns is an opt-in disable of connection
	// close by fasthttp after connections' HijackHandler returns.
	// This allows to save goroutines, e.g. when fasthttp used to upgrade
	// http connections to WS and connection goes to another handler,
	// which will close it when needed.
	KeepHijackedConns bool

	tlsConfig  *tls.Config
	nextProtos map[string]ServeHandler

	concurrency      uint32
	concurrencyCh    chan struct{}
	perIPConnCounter perIPConnCounter
	serverName       atomic.Value

	ctxPool        sync.Pool
	readerPool     sync.Pool
	writerPool     sync.Pool
	hijackConnPool sync.Pool

	// We need to know our listener so we can close it in Shutdown().
	ln net.Listener

	mu   sync.Mutex
	open int32
	stop int32
	done chan struct{}
}
```
注释写得很详细，这里只对部分字段说明一下。

* noCopy    fasthttp里面大量要求类型是不可复制的，即使用Server对象时，必须用指针，而不能直接拷贝，不然`go vet`会提示错误
* Concurrency、MaxConnsPerIP、MaxRequestsPerConn这几个都是限流字段，防止流量大的时候服务器被冲垮，或者做简单的DOOS防护
* perIPConnCounter是一个针对每个IP连接统计的计数，用于防止DDOS
* ctxPool, readerPool, writerPool, hijackConnPool都是对象池，以尽量减少GC

现在，让我们深入Serve方法来一窥究竟。

```go
func (s *Server) Serve(ln net.Listener) error {
	maxWorkersCount := s.getConcurrency()
	s.concurrencyCh = make(chan struct{}, maxWorkersCount)
	wp := &workerPool{
		WorkerFunc:      s.serveConn,
		MaxWorkersCount: maxWorkersCount,
		LogAllErrors:    s.LogAllErrors,
		Logger:          s.logger(),
		connState:       s.setState,
	}
	wp.Start()

	atomic.AddInt32(&s.open, 1)
	defer atomic.AddInt32(&s.open, -1)

	for {
		if c, err = acceptConn(s, ln, &lastPerIPErrorTime); err != nil {
			wp.Stop()
			if err == io.EOF {
				return nil
			}
			return err
		}
		s.setState(c, StateNew)
		atomic.AddInt32(&s.open, 1)
		if !wp.Serve(c) {
			atomic.AddInt32(&s.open, -1)
			s.writeFastError(c, StatusServiceUnavailable,
				"The connection cannot be served because Server.Concurrency limit exceeded")
			c.Close()
			s.setState(c, StateClosed)
			if time.Since(lastOverflowErrorTime) > time.Minute {
				s.logger().Printf("The incoming connection cannot be served, because %d concurrent connections are served. "+
					"Try increasing Server.Concurrency", maxWorkersCount)
				lastOverflowErrorTime = time.Now()
			}

			if s.SleepWhenConcurrencyLimitsExceeded > 0 {
				time.Sleep(s.SleepWhenConcurrencyLimitsExceeded)
			}
		}
		c = nil
	}
}
```

* 前两行代码生成了一个默认长度为256*1024的channel，也就是说fasthttp最多同时处理262144个用户的请求，实现上很简单，就是在开始处理的时候写一个进channel，完成的时候从channel读一个出来。
* 第3 ~ 10行代码创建了一个workerPool，它的最大worker数与上面的channel容量一样，默认都是256*1024，它的特点是FILO，也就是最后一个完成的worker将会被调度来执行下一个请求，这种设计据说能尽可能利用CPU hot cache。
* 第12、13行对server的open连接进行计数，这样server就可以知道当前有多少连接还在干活，在shutdown的时候等待open为0就可以了。你可能注意到，在golang里面，监控数据这样的频繁增减的需求，都会通过atomic来实现，以减少资源消耗，避免锁的开销。
* 接下来就是一个for循环，不断接受用户连接，进行处理
* wp.Serve就是具体的请求处理方法，如果无法处理，比如超过并发数什么的，就会返回一个错误的结果给客户端，并且会更新过载时间，如果过载时间超过了1分钟，则写一条log。
* 最后那个SleeWhenConcurrencyLimitsExceeded比较有意思，为了避免洪泛的时候，客户端快速的重试，故意在处理失败后将主循环等待，也就是说客满了，这段时间不接新客了。

## 限流及对抗DDOS

通过对上面的Serve方法的分析，我们可以看出，限流及对抗DDOS的思想在fasthttp框架中几乎是深入骨髓的。这里对其用到的手段做一下简单总结：

1. 通过concurrency channel限制最大的worker数，避免连接数过多
1. 通过FILO（先进后出）的workerPool来调度请求，提高CPU hot cache的利用率
1. 通过MaxConnsPerIP等在accept阶段就进行限流，如果同一IP连接数过多，直接返回错误
1. 当并发数超过系统最大负载时，记录发生的时间，并在超过1分钟后，记录日志
1. 可以通过SleeWhenConcurrencyLimitsExceeded在服务以及满负荷的情况下，让Server优先处理已有的用户，等一段时间再接受新的请求

## Hijacking（劫持）

`*RequestCtx`有一个方法Hijack，它接收一个handler函数，用于对请求的劫持，它是作用在net.Conn对象上，也就是TCP层上的。
一开始我不明白这玩意儿有什么用，反正都是HTTP请求，一来一回多简单，后来搜索看到有人将其用在gRPC和websocket上，前者不支持HTTP/1.1，通过hijack，可以退化到TCP层实现双向的RPC通信，后者由HTTP切换到TCP协议，也可以通过hijack作为中间人搞点事情。

注意：hijackHandler是在一个单独的goroutine执行的，所以你完全在里面多做点事情，比如做类似long polling这样的工作。

看一下官方的例子： 
```go
func ExampleRequestCtx_Hijack() {
	// hijackHandler is called on hijacked connection.
	hijackHandler := func(c net.Conn) {
		fmt.Fprintf(c, "This message is sent over a hijacked connection to the client %s\n", c.RemoteAddr())
		fmt.Fprintf(c, "Send me something and I'll echo it to you\n")
		var buf [1]byte
		for {
			if _, err := c.Read(buf[:]); err != nil {
				log.Printf("error when reading from hijacked connection: %s", err)
				return
			}
			fmt.Fprintf(c, "You sent me %q. Waiting for new data\n", buf[:])
		}
	}

	// requestHandler is called for each incoming request.
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		path := ctx.Path()
		switch {
		case string(path) == "/hijack":
			// Note that the connection is hijacked only after
			// returning from requestHandler and sending http response.
			ctx.Hijack(hijackHandler)

			// The connection will be hijacked after sending this response.
			fmt.Fprintf(ctx, "Hijacked the connection!")
		case string(path) == "/":
			fmt.Fprintf(ctx, "Root directory requested")
		default:
			fmt.Fprintf(ctx, "Requested path is %q", path)
		}
	}

	if err := fasthttp.ListenAndServe(":80", requestHandler); err != nil {
		log.Fatalf("error in ListenAndServe: %s", err)
	}
}
```

通过一个简单的Python脚本来向这个服务发送POST数据：

```python
import requests

s = requests.session()

print(s.post('http://127.0.0.1', data=b'hello').content)  # b'Hijacked the connection!
```

通过wireshark抓包，可以发现服务端事实上回的内容为，如果这时候通过上述session的TCP连接去receive数据，就能得到后面的数据了。

```
POST /hijack HTTP/1.1
Host: 127.0.0.1:80
User-Agent: python-requests/2.20.1
Accept-Encoding: gzip, deflate
Accept: */*
Connection: keep-alive
Content-Length: 5

helloHTTP/1.1 200 OK
Server: fasthttp
Date: Fri, 31 May 2019 07:18:20 GMT
Content-Type: text/plain; charset=utf-8
Content-Length: 24

Hijacked the connection!This message is sent over a hijacked connection to the client 127.0.0.1:57582
Send me something and I'll echo it to you
```

关于Hijack的更多信息，可以看标准库对应的代码：https://golang.org/pkg/net/http/#Hijacker
由于HTTP/2已经支持Server Push，对于HTTP/2协议，Hijack是不支持的。

## NEXT...

关于服务端的部分，暂时介绍这么多，下一篇会分析下Client部分，然后再对于里面用到的各种优化手段作一番剖析。
