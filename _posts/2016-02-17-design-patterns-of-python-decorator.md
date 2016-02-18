---
layout: post
title:  "Python设计模式 - 装饰器模式"
date:   2016-02-17 15:00:00
categories: Python
---
>修饰模式，是面向对象编程领域中，一种动态地往一个类中添加新的行为的设计模式。就功能而言，修饰模式相比生成子类更为灵活，这样可以给某个对象而不是整个类添加一些功能。
>通过使用修饰模式，可以在运行时扩充一个类的功能。原理是：增加一个修饰类包裹原来的类，包裹的方式一般是通过在将原来的对象作为修饰类的构造函数的参数。装饰类实现新的功能，但是，在不需要用到新功能的地方，它可以直接调用原来的类中的方法。修饰类必须和原来的类有相同的接口。

装饰器（Decorator）模式是我使用率很高的一种模式，粗看它时有点像黑魔法，一旦熟悉了，相信很多人都会喜欢上它的。

### 从一个问题开始

让我们从一个具体的问题开始：

```python
# 编写一个函数，获取一个URL的内容
import urllib

def fetch(url):
    return urllib.urlopen(url).read()
```

上面这段代码很简单，获取一个URL的内容。但这时我们遇到了一个问题，由于网络状况或者网站的负载等原因，有些情况下访问会失败，但是经过一些重试之后就可以成功。这个时候，我们就需要把代码做一些修改，见下面的代码：

```python
import urllib
import time

def fetch(url):
    for _ in xrange(5):
        try:
            return urllib.urlopen(url).read()
        except:
            time.sleep(1)
    else:
        raise RuntimeError
```

在这个改进版的fetch函数中，遇到访问失败，我们会进行重试，每次等待间隔1秒，重试5次。

这个函数明显臃肿了很多，更大的问题是，多出来的这些代码与这个函数的目的**fetch**是无关的，它们只是为了处理重试而存在，而重试这个需求在很多地方都是通用的，它可以被拿出来。

这个时候，就该装饰器上场了。

### 什么是装饰器

装饰器是一个函数，它的使用方法是在被装饰的函数前加上`@<decorator>`，从原理上说，装饰器函数是一个接收函数作为输入参数，返回一个新的函数的函数。见如下例子：

```python
def decorator(f): # 输入参数为函数
    def wrapper(*args, **kwargs): # wrapper函数，用于替代f
        print 'start'
        result = f(*args, **kwargs)
        print 'end'
        return result

    return wrapper # 返回wrapper函数

# 被装饰函数
def test(name):
    print name

# 不用装饰器语法
wrapper_test = decorator(test)
wrapper_test('hello')
# OUTPUT
# start
# hello
# end
```

我们在上面定义了一个装饰器函数decorator，一个被装饰函数test，当我们需要为test函数增加功能的时候，通过decorator(test)生成一个新的函数来实现。可以看到，这样的代码有点罗嗦，并且不太清楚，如果test有多个装饰器，这里会生成多个类似wrapper_test的新函数，所以我们渴望有一个好用的语法糖，而Python提供了。

```python
def decorator(f):
    def wrapper(*args, **kwargs):
        print 'start'
        result = f(*args, **kwargs)
        print 'end'
        return result

    return wrapper

@decorator
def test(name):
    print name

test('hello')
# OUTPUT
# start
# hello
# end
```

可以看到通过@decorator装饰后的代码与之前效果一样，但是语法简洁了很多。

### 使用装饰器重构上面的代码

让我们使用装饰器来重构上面的代码，把retry的部分独立成一个装饰器函数，它有两个参数：重试次数和重试间隔。所以我们需要写一个接受参数的装饰器，它的返回值是一个新的装饰器函数，代码如下：

```python
import urllib
import time
import functools

def retry(times, interval):
    def _retry(f):
        @functools.wraps(f)
        def wrapper(*args, **kwds):
            for _ in xrange(times):
                try:
                    return f(*args, **kwds)
                except:
                    time.sleep(interval)
            else:
                raise RuntimeError

        return wrapper

    return _retry

@retry(5, 1)
def fetch(url):
    return urllib.urlopen(url).read()
```

这时候，fetch有回到了最初的样子，简单明确，而retry作为一个独立的函数则可以被很多其他地方复用，我们成功地把两者解藕了。

又来了一个需求，为了提高性能，需要对fetch函数做cache，5秒内访问过的url无法再次请求，如果采用装饰器模式，我们的代码应该是这样的：

```python
@cache(5)
@retry(5, 1)
def fetch(url):
    return urllib.urlopen(url).read()
```

这就是装饰器模式，在很多地方都有它的应用，比如最常见的property。

```python
class MyObj(object):
    @property
    def name(self):
        return 'MyObj', self.__hash__()
```

在一些library如[bottle](http://bottlepy.org/)里面也大量使用了装饰器，如：

```python
from bottle import route, run, template

@route('/hello/<name>')
def index(name):
    return template('<b>Hello {{name}}</b>!', name=name)

run(host='localhost', port=8080)
```

### 小结

装饰器模式是让应用解藕的一个非常好用的模式，对于认证、缓存、重试等需求，用该模式可以在不改变现有代码逻辑的情况下添加增强功能。

但是，也需要注意的是，不是什么代码都适合放在装饰器里面的，如果那本来就是函数逻辑的一部分，那还是放在函数内部吧，另外在做单元测试的时候，我们通常也会把装饰器都mock掉，以方便测试。

### 参考

* [PEP 0318 -- Decorators for Functions and Methods](https://www.python.org/dev/peps/pep-0318/)
* [PythonDecorators](https://wiki.python.org/moin/PythonDecorators)

