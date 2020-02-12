---
layout: post
title:  "gRPC测试模拟实践"
date:   2020-02-12 22:00:00 +0800
categories: "DevOps"
---
**status: draft**

[grpc](https://grpc.io)是一个被广泛使用的RPC协议，由于它高性能、跨语言的特点，被很多基于微服务架构的产品所采用。而随着产品中微服务数量的增加，针对它的测试需求也逐渐显露出来，包括：

1. 怎么测试一个gRPC服务
2. 怎么测试一个依赖gRPC服务的应用

在开始具体的实践之前，我们先来认识一下grpc协议。

## grpc简介

### 协议简述

grpc的协议栈如下：

![grpc protocol stack]({{ site.url }}/assets/grpc/protocol_stack.png)

我们用[grpc examples](https://github.com/grpc/grpc/tree/master/examples)里面那个最简单的helloworld程序为例，它的proto文件定义如下：

```proto
syntax = "proto3";

option java_multiple_files = true;
option java_package = "io.grpc.examples.helloworld";
option java_outer_classname = "HelloWorldProto";
option objc_class_prefix = "HLW";

package helloworld;

// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```

通过抓包，可以看到一次grpc的request报文如下：

![grpc request]({{ site.url }}/assets/grpc/grpc_request.png)

这里可以看出，grpc协议跑在HTTP2上，一次RPC调用就是一次POST method，URI的第一部分是service名（这里是helloworld.Greeter），第二部分就是method名（这里是SayHello）。而请求内容则是protobuf序列化的具体数据you。

相应的，一个response的报文为：

![grpc response]({{ site.url }}/assets/grpc/grpc_response.png)

响应也是普通的HTTP2报文，可以看到这次的response的status code为200，返回的数据为protobuf序列化的Hello, you!。

> 注：为了便于分析协议，我们让grpc工作在非加密通道上。

### 定义一个服务

我们以go语言为例，来看一下如何定义一个gRPC服务。

```go
package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```

其中，`google.golang.org/grpc/examples/helloworld/helloworld`是用protoc工具从helloworld.proto文件创建出来的。我们打开[helloworld.pb.go](https://github.com/grpc/grpc-go/blob/master/examples/helloworld/helloworld/helloworld.pb.go)，可以看到里面关键的一个接口定义为：

```go
type GreeterServer interface {
	// Sends a greeting
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
}
```

所以，在我们自己的程序中，我们具体实现了SayHello，作为接口暴露出来。

### 使用grpcurl工具访问服务

相比官方的grpc_cli，go语言编写的[grpcurl](https://github.com/fullstorydev/grpcurl)安装使用都更为方便一些。

1. 查看服务或者proto文件提供的grpc接口

   ![image-20200211152833205](/Users/zhang/Library/Application Support/typora-user-images/image-20200211152833205.png)

   这里第二个服务开启了[ServerReflection](https://github.com/grpc/grpc/blob/master/doc/server-reflection.md)。

2. 查看某个接口的方法

   ![image-20200211153022160](/Users/zhang/Library/Application Support/typora-user-images/image-20200211153022160.png)

3. 访问某个方法

   ![image-20200211153259688](/Users/zhang/Library/Application Support/typora-user-images/image-20200211153259688.png)

## 测试grpc服务

### 单元测试

一个grpc方法就是一个普通的函数，因此针对它的单元测试，跟其它函数的单元测试类似，这里给个上面SayHello的例子，不再赘述。

```go
package main

import (
  "testing"

	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

func TestSayHello(t *testing.T) {
  in := &pb.HelloRequest{Name: "world"}
  s := &server{}
  out, err := s.SayHello(nil, in)
  if err != nil {
    t.Error(err)
  }
  if out.GetMessage() != "Hello world" {
    t.Errorf("invalid response %s", out.GetMessage())
  }
}
```

### 模块测试

我们这里关注下怎么对一个grpc服务进行功能性的测试。我们先看一下一个grpc client通常是怎么实现的，还是以上面那个helloworld为例。

```go
package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
```

可以看到，作为client，它也需要引用protoc生成的代码。这样，当我们的测试应用或者接口描述改变的时候，测试代码也需要相应的更改，显然，直接参照官方client的实现思路来进行client测试不是一个好的办法。

其实，使用过grpcurl就会有些疑惑，我们并没有帮它生成proto语言定义文件，它是怎么发送grpc消息的呢？翻阅它的代码，我们发现，有一个叫“github.com/jhump/protoreflect”的模块帮助我们实现了动态的grpc消息，借助它，grpcurl实现了动态组装和解析grpc消息的功能。

所以，简单地通过封装grpcurl，我们就可以对grpc服务的进行测试了。例如，在[simgo](https://github.com/feiyuw/simgo/)项目中，通过web界面连接到grpc服务，进行测试。

![simgo client]({{ site.url }}/assets/grpc/simgo_client.png)

## 测试依赖grpc服务的应用

微服务化改造带来的一大难题，就是分布式应用的测试问题。见下面这张图：

![grpc arch]({{ site.url }}/assets/grpc/grpc_arch.png)

假设我们要对APP进行测试，由于它依赖了六个gRPC服务，我们需要把他们完整地部署起来才行。而如果这中间有一个服务拖了后腿，测试就得延后，这种情况显然是我们不愿意看到的。所以，为了对APP进行高效测试，我们通常采取契约测试的方法，通过模拟其他依赖服务的方式来进行验证。

simgo同样提供了对grpc server的模拟支持，以单元测试为例：

```go
import (
  "testing"
  
	"github.com/feiyuw/simgo/protocols"
)

func TestAppWithSimGrpc(t *testing.T) {
  // create simulated grpc server
  s, _ := NewGrpcServer(":4999", []string{"echo.proto", "helloworld.proto"})
	s.SetMethodHandler("helloworld.Greeter.SayHello", func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error {
		out.SetFieldByName("message", in.GetFieldByName("name"))
		return nil
	})
	s.Start()
  defer s.Stop()
  
  // test App
}
```

可以看到，通过提供接口定义的proto文件和服务端口，simgo就可以生成一个模拟的grpc服务，通过SetMethodHandler方法可以对某个方法的行为进行模拟。

在实际工程实践中，建议将grpc协议描述独立于具体服务保存，并版本化，保证契约测试的有效性。

## 下一步做什么

至此，我们基本能在单元测试中完成测试一个gRPC服务和测试一个依赖gRPC服务应用。接下来就是借助一些好的实践方法来提升测试效率。

### 与测试框架集成

首先，是与测试框架的集成，笔者比较推荐BDD风格的测试框架，比如[gauge](https://gauge.org)。如果你用golang编写测试代码，那可以将simgo作为go module直接使用。如果更偏向于python等语言来编写测试代码，可以通过调用simgo提供的RESTful接口来实现类似的目的。

### 与持续交付流水线集成

测试的投入产出比随着执行频率的增加而提高，因此，尽早将你的测试在持续交付流水线中跑起来吧。可能除了grpc server的模拟，还有其他依赖需要处理，我个人的建议是如果成本不太高，就用真实的服务（如MySQL数据库），然后再考虑Mock和simulator等手段。

## 结语

当前simgo这个项目还处于非常初级的状态，欢迎提交Issue和PR。

## 参考

* [A short introduction to Channelz](https://grpc.io/blog/a_short_introduction_to_channelz/)
* [bloomrpc，一个NodeJS实现的grpc client](https://github.com/uw-labs/bloomrpc)
* [grpcurl，一个golang实现的grpc client](https://github.com/fullstorydev/grpcurl)
* [grpc-ecosystem](https://github.com/grpc-ecosystem)
* [grpc binary encoding](https://youtu.be/VTdkDu4OGbE)
* [protobuf encoding](https://developers.google.com/protocol-buffers/docs/encoding)
* [awesome grpc](https://github.com/grpc-ecosystem/awesome-grpc)
* [grpc-protobuf的动态加载及类型反射实战]([http://xiaorui.cc/2019/04/01/grpc-protobuf%E7%9A%84%E5%8A%A8%E6%80%81%E5%8A%A0%E8%BD%BD%E5%8F%8A%E7%B1%BB%E5%9E%8B%E5%8F%8D%E5%B0%84%E5%AE%9E%E6%88%98/](http://xiaorui.cc/2019/04/01/grpc-protobuf的动态加载及类型反射实战/))
* [grpc-reflection-and-grpcurl](https://about.sourcegraph.com/go/gophercon-2018-grpc-reflection-and-grpcurl/)
* [grpc_cli](https://github.com/grpc/grpc/blob/master/doc/command_line_tool.md)


