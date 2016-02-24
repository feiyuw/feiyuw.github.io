---
layout: post
title:  "Python设计模式 - 代理模式"
date:   2016-02-24 15:00:00
categories: Python
---
**代理模式**在实际开发中应用极为广泛，通过它能把原本复杂而耗时的操作优雅地隐藏起来，提供简洁的接口，比如对与远程RESTful API的调用，对于系统命令的操作等。

在维基百科上，**代理模式**的定义是这样的：

>代理模式（英语：Proxy Pattern）是程序设计中的一种设计模式。
>所谓的代理者是指一个类可以作为其它东西的接口。代理者可以作任何东西的接口：网络连接、存储器中的大对象、文件或其它昂贵或无法复制的资源。
>著名的代理模式例子为引用计数（英语：reference counting）指针对象。

**代理模式**的关键就在代理两个字上，我们以两个例子来看看怎么实现一个代理。

### RESTful API

有很多项目Python一边用作Web的后端，提供REST服务，另外又充当client处理脚本，用于同步数据等。

在client上，我们调用REST API通常可能是这样的：

```python
import urllib
import json

def fetch_resource(resource_id):
    opener = urllib.urlopen('http://remote.server/api/resource/' + resource_id)
    if opener.code != 200:
        raise RuntimeError('invalid return code!')
    content = opener.read()
    try:
        return json.loads(content)
    except ValueError:
        return content
```

对于每一个REST操作，我们都会写一段类似的代码，这些代码基本一样，差别的地方可能就在API的地址和HTTP method（POST、GET、PUT等）上。而且后续我们还可能会对这个URL操作进行重试，加入cache等，这样频繁修改不同地方的类似代码实在是太罗嗦了。

为了减少痛苦，我们引入一个Proxy，所有痛苦的事情都交给它来做。

```python
import urllib
import json

class GetProxy(object):
    def __getattr__(self, api_path):
        def _rest_fetch(*paras):
            opener = urllib.urlopen('http://remote.server/api/' + api_path + '/' + '/'.join(resource_id))
            if opener.code != 200:
                raise RuntimeError('invalid return code!')
            content = opener.read()
            try:
                return json.loads(content)
            except ValueError:
                return content

        return _rest_fetch

proxy = GetProxy()

# 调用API
proxy.user(123) # http://remote.server/api/user/123
proxy.resource('switch', 456) # http://remote.server/api/resource/switch/456
```

### 调用系统命令

[TODO]

