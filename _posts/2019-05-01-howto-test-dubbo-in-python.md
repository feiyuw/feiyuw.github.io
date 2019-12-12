---
layout: post
title:  "利用Python测试标准dubbo provider和consumer"
date:   2019-05-01 22:00:00 +0800
categories: "Python"
---
最近工作中经常与各类dubbo服务打交道，而在测试框架选型的时候考虑到开发效率和团队背景，选择了Python作为基础语言，这样就带来一个问题，如何用Python来测试dubbo服务？以及如何用Python来模拟dubbo服务被其它dubbo消费者调用？

dubbo默认采用hessian2的编解码方式，没有找到现成的库，而转换为其它编码方式如jsonrpc等，在有些场景下也不可行，于是就决定自己撸一个，也就有了[dubbo-py](https://github.com/feiyuw/dubbo-py)这个项目。

### 项目进展

当前这个项目还处于初期，但已经能满足我当前的需求了，如果您遇到问题，欢迎在github上给我提[issue](https://github.com/feiyuw/dubbo-py/issues)。

主要功能的实现情况：
- [x] 大部分hessian2协议的编解码（我用到的那些吧）
- [x] Python DubboClient调用Java的Dubbo服务
- [x] Python DubboService可以注册到zookeeper上作为Dubbo Provider
- [x] 代码重构，让Python的rpc函数只需要关心调用参数与返回值（和普通函数一样），其它都交给library实现
- [ ] Python DubboClient注册到zookeeper上作为Dubbo Consumer
- [ ] 添加AsyncDubboService以支持asyncio
- [ ] 更多地单测覆盖
- [ ] 支持Python2.7

### 安装

```sh
# python >= 3.6
pip3 install dubbo-py
```

### 示例

```python
from dubbo.codec.hessian2 import DubboResponse, JavaList
from dubbo.server import DubboService
from dubbo.client import DubboClient


def remote_max(*args):
    return max(args)


def remote_sum(lst):
    return sum(lst)


service = DubboService(12358, 'demo')
service.add_method('com.myservice.math', 'max', remote_max)
service.add_method('com.myservice.math', 'sum', remote_sum)
# service.register('127.0.0.1:2181', '1.0.0')  # register to zookeeper
service.start()  # service run in a daemon thread

client = DubboClient('127.0.0.1', 12358)
resp = client.send_request_and_return_response(service_name='com.myservice.math', method_name='max', args=[1, 2, 3, 4])
print(resp.data)  # 4
resp2 = client.send_request_and_return_response(service_name='com.myservice.math', method_name='sum', args=[JavaList([1, 2, 3, 4])])
print(resp2.data)  # 10
```

### 代码结构

dubbo-py的代码很简单，目前也没来得及好好整理，它要求Python3.6以上环境（对Python3.5和Python2.7的支持在计划中），依赖[kazoo](https://kazoo.readthedocs.io/en/latest/)，借助它实现与zookeeper的交互。

* dubbo
    * server    实现了DubboService类，用于模拟Dubbo Provider，感觉后续名字改成provider更合适
    * client    实现了DubboClient类，用于模拟Dubbo Consumer，感觉后续名字改成consumer更好
    * codec     协议编解码模块，目前只实现了hessian2协议的编解码，好像也只需要这个就够了
* tst    单元测试，采用pytest框架
