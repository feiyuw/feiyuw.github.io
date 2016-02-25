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

另一个非常常见的例子是在Python中调用系统命令，比如调用ping、ls等命令。在Python中调用系统命令的代码一般像下面这样：

```python
import subprocess

def ping(*args):
    cmd = ['ping'] + [str(arg) for arg in args]

    return subprocess.call(cmd)
```

由于执行命令的代码基本上都差不多，我们不希望重复去写这些逻辑。另外，每次都调用subprocess.call，在阅读上不够明了。我们真正希望的调用方式是这样的：

```python
from myshell import ping

ping('-c', 5, '127.0.0.1')
```

在这里，`myshell`模块就是我们的proxy。具体怎么做呢？我们需要借助`__getattr__`和module重定义。

```python
# myshell.py
import sys

class _ShellProxy(object):
    def __getattr__(self, cmd):
        def run_cmd(*args):
            import subprocess
            return subprocess.call(['ping'] + [str(arg) for arg in args])

        return run_cmd

sys.modules[__name__] = _ShellProxy()

# usage.py
from myshell import ping
from myshell import ls
ping('-c', 5, '127.0.0.1')
ls('/')
```

在Python中**一切皆是对象**，包括module的import，所以在from myshell import xxx的时候实际上类似于访问myshell.xxx，也就是getattr(myshell, 'xxx')。所以如果myshell模块是一个类的实例，我们就可以借助`__getattr__`来实现系统命令的路由了。

要达到这个目的，我们需要将myshell模块动态替换为一个类的实例，而最后一句`sys.modules[__name__] = _ShellProxy()`就是把myshell模块重写为_ShellProxy的实例。


### 参考

* [\_\_getattr\_\_ in Python](https://docs.python.org/2/reference/datamodel.html#object.__getattr__)
* [sh module of Python](https://github.com/amoffat/sh)
