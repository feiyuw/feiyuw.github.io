---
layout: post
title:  "契约测试中的服务simulator"
date:   2019-12-06 12:00:00 +0800
categories: "DevOps"
---
随着微服务架构的兴起和系统复杂度的提高，针对单个服务的测试环境成本越来越高，究其原因，是因为各种rpc调用导致的服务间耦合，使得我们很难端到端的测试一个服务。

这给软件研发和测试带来了两个挑战：

1. 端到端测试成本急剧上升，在测试环境维护和问题分析追踪上花费大量时间，很难做到一次提交全面覆盖
2. 测试后置，导致软件发布的风险增加，由于测试环境依赖程度高，导致联调、单应用测试很难在早期开展，很多验证需要等待所有系统集成阶段进行，大大增加了软件发布风险

这与微服务架构的初衷是违背的，我们切换到微服务架构，是希望能够做到应用独立开发、独立测试、独立部署。

为了解决这个问题，业内出现了多种实践，常见的有两种：

1. 基础环境 + 测试环境，通过中间件隔离基础环境和测试环境，让单个应用在测试时不需要部署完整的测试环境，只需要部署自己，加上公共的基础环境，即可完成测试。这个方法的好处是，可以在不增加环境成本的情况下，做到方便的端到端测试。但是，由于需要做数据隔离和链路隔离，所以需要中间件的支持，一般只有上规模的企业才有能力维护。
2. 契约测试，通过模拟器来模拟应用依赖的服务，模拟的服务接口通过契约来定义。这个方法不需要对中间件的改造，用起来成本较低。但是，开发这套模拟器需要一定的成本，比如自定义协议的实现，如何在模拟器和真实服务间同步契约也需要所有人的努力。

## 如何开始

先抛出一个观点，无论采用的语言和技术栈是什么，实现服务模拟器，先从实现一个模拟服务的library开始，这个library需要支持：

1. 实现服务所用通信协议的协议栈，包括心跳保活、协议编解码等
2. 能够创建和启动一个模拟服务
3. 能够停止一个模拟服务
4. 能够控制模拟服务针对某个特定请求的响应动作，如：
   1. 回复特定的内容
   2. 延时回复
   3. 抛出特定的错误
   4. 不回复
   5. 其他可能的响应动作
5. 能够得到模拟服务接收到的（来自特定用户的）数据请求
6. 如果协议是双向的，能够通过函数控制模拟服务主动发送特定的数据到对端
7. 能够通过函数获得模拟服务的所有客户端信息（*）
8. 能够通过函数控制模拟服务的某个客户端的连接断开等操作（*）

从这个支持看，我们会发现即使是HTTP这样非常通用的服务模拟，很多现成的框架也不能直接使用，因为从设计角度讲，他们关注的是对外服务，模拟器关注的是对内控制。

下面我以一个案例来具体说明一下。

有一个采用Java语言和dubbo RPC框架的产品，其应用架构图简单描述如下：

![App Arch]({{ site.url }}/assets/simulator/app_arch.png)

> APP是我们的被测应用，Service是它所依赖的其他应用和服务，APP通过dubbo协议与它们通信。服务发现通过Zookeeper来完成，所有的dubbo服务都会把自己注册到Zookeeper上。

这里也对dubbo协议做下简单描述：

![protocol]({{ site.url }}/assets/simulator/protocol.png)

> 注：上图是很简略的描述，具体的协议表示可以参阅[dubbo官网协议参考手册](https://dubbo.apache.org/zh-cn/docs/user/references/protocol/dubbo.html)。
>
> dubbo协议有多种序列化方式，比如jsonrpc、thrift、hessian等等，官方默认的是hessian2，采用TCP长连接形式，包括head和payload两部分组成。所以请注意这里的hessian2指的是数据序列化方式。

我希望能对APP单独进行测试，此时遇到的第一个问题是测试环境，由于APP依赖很多RPC service，而这些service又依赖更多的RPC service，同时，他们还可能依赖其它服务如DB、Cache、MQ等等，导致部署一套完整的测试环境需要多台服务器，启动几十个服务，成本较高，且维护起来比较麻烦。同时，由于问题追踪的链路较长，当出现问题时，对问题的定位比较耗时，导致回归效率低下。

第二个问题是当新的服务提出，或原有服务被修改（原则上禁止同版本服务不兼容老代码上线）的情况下，对服务的测试需要等到对应服务部署到测试环境才能进行，极大拖慢了开发进度。

于是，我开始着手实现一个针对dubbo的服务模拟器，并确定了要支持的第一个场景。

这个场景很简单，简单描述如下：

> 应用A提供一个Restful接口给终端用户使用，在处理用户请求的过程中，它会通过dubbo RPC调用服务B的一个method，构建一个只有应用A的测试环境，请求其提供的Restful接口，得到正确的返回。
>
> 为了简单起见，第一步我们假设服务B总是返回确定的数据。

当时功能测试框架采用的是[RobotFramework](https://robotframework.org/)，我采用Python语言来实现这个模拟library。

对于这个场景，首先要声明服务B和方法，并将其注册到Zookeeper上，然后启动模拟服务，代码大致如下：

```python
def handler():
    return None


service = DubboService(12358, 'demo')
service.add_method('com.myservice.first', 'hello', handler)
service.register('127.0.0.1:2181', '1.0.0')  # register to zookeeper
service.start()  # service run in a daemon thread

```

当模拟服务启动后，被测应用A会通过Zookeeper自动发现该服务，然后我们通过HTTP请求应用A的restful接口就行了，类似代码如下：

```python
# 假设使用requests模块
import requests


a_api_url = 'http://a.dns/api/v1/hello'

resp = requests.get(a_api_url)
assert resp.ok is True
# ...
```

从前面的协议栈，我们可以大致把这个library分为三个模块：

1. hessian2 codec
2. TCP service
3. Zookeeper handler

Hessian2 codec部分基本上是把Dubbo的序列化用Python重新实现了一遍。

TCP service关注于服务的start、stop、handler等，大致代码如下：

```python
class DubboService(object):
		# ...

    def start(self):
        # ...
        pass

    def stop(self):
        # ...
        pass

    def add_method(self, service, method, handler):
        # ...
        pass

```

Zookeeper handler包含将服务注册到Zookeeper和服务发现的功能。

关于这部分的实现细节有兴趣的可以参照[dubbo-py](https://github.com/feiyuw/dubbo-py)这个项目。

## 下一步做什么

一个简单的准则：技术要为业务服务。

对于模拟器的开发也是这样，我们要明确几点：

1. 什么协议的模拟器？
2. 什么服务的模拟器？
3. 当下有哪些应用场景？
4. 将来有哪些应用场景？

保证一个时间聚焦到一个问题，采用持续交付的方式，让模拟器的开发被需求（如功能自动化测试用例）所拉动，同时一旦需要的功能就绪，对应的需求（如功能自动化测试用例）立即进入下一流程（如进入每次提交的自动化验证）。

KISS（keep it simple, stupid)准则非常适用于模拟器的开发实践，模拟器通常应该只暴露非常少量的接口，关注于连接管理和消息处理，**不碰业务**。

那具体的业务层怎么办呢？比如我们要模拟10个dubbo服务，每个的方法和接口定义都不一样，这些在测试人员看来也是基本功能，总不能每次都要测试人员去写上面那些函数定义和注册代码吧？

大部分情况下，模拟器的使用者都是功能测试自动化，从这个角度讲，功能测试自动化用例就是它的用户，站在用户的角度，它希望的是自然语言描述的接口，假设以robotframework的语言描述，如：

```robot
*** Test Cases ***
App A can handle API correctly
    Start Mock Service B
    Service B should be registered into Zookeeper
    Start APP A
    Request to API /hello
    Response should be OK
```

可以看到用户角度的接口和我们前面模拟器提供的接口存在巨大的鸿沟，所以我们需要在两者中间加上一两层，来填平这个鸿沟。

还是以robotframework测试框架为例，描述一下分层策略

![TA Arch]({{ site.url }}/assets/simulator/ta_arch.png)

> 这里注意几点：
>
> 1. Simulator core和Simulated service就是一般的python module，不会引用测试框架的任何模块，做到足够的通用，方便后续支持不同测试框架。
> 2. Basic Robot Library和Functional Robot Library作为针对robotframework的library实现，会引用到robotframework的模块，但应该做到越少越好，并兼容不同的robotframework版本
> 3. 上层模块可以import下一层模块，反过来不允许
> 4. 不允许跨层import
> 5. 对于用户（功能自动化测试用例），它只需要关心Functional Robot Library和Basic Robot Library，这两个library会提供它在case里用到的像start_mock_service_B这样的接口。
> 6. Functional Library和Basic Library的区别是前者操作一个业务，后者仅操作一个接口。如点赞这个业务，可能涉及多个rpc请求，每个rpc接口对于Basic Librry就是一个函数。

这样分层之后，每个层可以由不同背景的人来维护，同时，由于core和service足够简单，单元测试变得非常容易，质量也可以得到保证。而上层的library涉及业务的就可以交由具体的开发和业务测试同学维护，做到及时更新。

建议在尽可能早的时候明确这个层次划分。

## 如何保证模拟器的质量

模拟器质量要求其实非常高，因为它是测试其他软件的基础软件，它出现问题，会导致我们在问题排查中浪费大量的时间，同时也会让研发和测试同学对自动化测试失去信任，所以对待模拟器的质量，需要给与足够的重视。

首先，坚持KISS原则，做得尽量简单，把业务模块和基础模块分离，就像上面做的那样。

其次，做到充分的测试覆盖，对代码质量有敬畏之心，从写下第一个接口开始，模拟器就要有单元测试覆盖，并且每发现一个模拟器的缺陷，都要做到有单元测试来验证。

第三，做到接口的向前兼容，确保用例不会因为模拟器版本的更新而失败。

第四，做到模拟器和测试library的持续交付，做到一天N次高质量的发布。

## 其他思考

做好一个产品很不容易，需要很多思考，做到高质量、快速、低成本地交付应该是所有参与者追寻的目标。

软件质量绝不仅是测试的事，应用本身要做到可测试、可观测、可监控、可运维，如果发现测试一个应用非常困难，需要停下来想一想，这个应用的拆分和架构是不是有问题，是不是需要改进？

抛砖引玉，希望能引起一些思考。

## 参考

* [contract testing - Martin Fowler](https://martinfowler.com/bliki/ContractTest.html)
* [pact](https://docs.pact.io/)
* [e2e vs contract based testing](https://techbeacon.com/app-dev-testing/end-end-vs-contract-based-testing-how-choose)

