---
layout: post
title:  "Python设计模式 - 单例模式"
date:   2016-02-16 15:23:00
categories: Python
---
[Python](https://www.python.org)语言当然也是有设计模式的，有人说，设计模式是为语言的设计打补丁，因为语言本身无法直接在语法层面支持。而谈设计模式，多数都采用JAVA或者C++语言来描述，为免引起口水战，笔者不打算比较Python与这些语言的特点，而是从我自己的实践来聊聊我所用到和了解的Python设计模式。

今天先聊聊**单例模式**，所谓单例模式，在维基百科上它的定义如下：

>单例模式，也叫单子模式，是一种常用的软件设计模式。在应用这个模式时，单例对象的类必须保证只有一个实例存在。许多时候整个系统只需要拥有一个的全局对象，这样有利于我们协调系统整体的行为。比如在某个服务器程序中，该服务器的配置信息存放在一个文件中，这些配置数据由一个单例对象统一读取，然后服务进程中的其他对象再通过这个单例对象获取这些配置信息。这种方式简化了在复杂环境下的配置管理。

---

## import

在Python语言里面，module的import就是一种典型的单例模式。先看一个例子：

```python
# myglobal.py
COUNT = 0
```

写一个多线程的测试程序，每个线程都把COUNT加一，并且打印出来。

```
In [1]: from multiprocessing.dummy import Pool

In [2]: p = Pool(10)

In [3]: def count(x):
   ...:     import myglobal
   ...:     myglobal.COUNT += 1
   ...:     print myglobal.COUNT
   ...:

In [4]: p.map(count, range(10))
1
 3
 5
2
6
7
 9
8
10
4
Out[4]: [None, None, None, None, None, None, None, None, None, None]

In [5]: import myglobal

In [6]: myglobal.COUNT
Out[6]: 10
```

可以看到，在每个线程中都把COUNT加了一，得到不同的值。在最后也按照期望得到了COUNT = 10。

我们再看一下每个线程里面myglobal的内存地址，就可以明白为什么说它是单例模式了。

```
In [1]: from multiprocessing.dummy import Pool

In [2]: p = Pool(10)

In [3]: def imp(x):
    import myglobal
    print id(myglobal)
   ...:

In [4]: p.map(imp, range(10))
140519408948848
140519408948848
 140519408948848
140519408948848
 140519408948848
 140519408948848
140519408948848
140519408948848
140519408948848
140519408948848
Out[4]: [None, None, None, None, None, None, None, None, None, None]
```

可以看到，print出来的值是相同的。

现在让我们想想import的用法，就可以明白为什么它要设计成这样了。

我们都知道，import一个模块之后，会在程序的不同地方调用这个模块，可能在不同的函数，也可能在不同的线程中。我们希望无论在什么地方，这个模块的行为是一致的，也就是说，我们希望它是唯一的，而不是一个独立的个体。

---

## logging

`logging`是单例模式的一个非常典型的应用，因为我们的应用内部一般不去关心logging的细节，比如写到数据库还是文件，时间戳的格式等等。在应用内部，我们只关心logging的分类，内容和级别，是app的还是db的log，是warning，error还是debug。

而logging的配置，则是全局的，一个进程中任何地方对logging配置的修改（如增加输出到文件），都会影响到所有使用logging的地方。那么它是怎么做到的呢？

看一个例子：

```python
import logging
# 全局的配置
logging.basicConfig(format="%(asctime)s %(levelname)s %(name)s %(message)s", level='DEBUG')
# 在函数中如下使用
logging.debug('hello')

# 在其他module中，sub.py
import logging
logging.info('this is the log in sub module')
```

打印出来的log如下：

```
2016-02-17 10:26:33,100 DEBUG root hello
2016-02-17 10:29:18,274 INFO root this is the log in sub module
```

可见两个log的行为和格式是一样的。

让我们看一下logging模块的相关代码，比如logging.info。

```python
root = RootLogger(WARNING)

def info(msg, *args, **kwargs):
    """
    Log a message with severity 'INFO' on the root logger.
    """
    if len(root.handlers) == 0:
        basicConfig()
    root.info(msg, *args, **kwargs)
```

可以看到具体的执行都是有root这个实例来进行的，由于root是一个module层面的实例，并不在info函数内部实例化，所以所有的logging.info都是同一个实例来进行的。

---

## 共享的instance

在很多讲述Python单例模式的文章中都能看到一段代码，我把它摘录如下：

```python
# https://github.com/faif/python-patterns/blob/master/borg.py
#!/usr/bin/env python
# -*- coding: utf-8 -*-


class Borg:
    __shared_state = {}

    def __init__(self):
        self.__dict__ = self.__shared_state
        self.state = 'Init'

    def __str__(self):
        return self.state


class YourBorg(Borg):
    pass

if __name__ == '__main__':
    rm1 = Borg()
    rm2 = Borg()

    rm1.state = 'Idle'
    rm2.state = 'Running'

    print('rm1: {0}'.format(rm1))
    print('rm2: {0}'.format(rm2))

    rm2.state = 'Zombie'

    print('rm1: {0}'.format(rm1))
    print('rm2: {0}'.format(rm2))

    print('rm1 id: {0}'.format(id(rm1)))
    print('rm2 id: {0}'.format(id(rm2)))

    rm3 = YourBorg()

    print('rm1: {0}'.format(rm1))
    print('rm2: {0}'.format(rm2))
    print('rm3: {0}'.format(rm3))

### OUTPUT ###
# rm1: Running
# rm2: Running
# rm1: Zombie
# rm2: Zombie
# rm1 id: 140732837899224
# rm2 id: 140732837899296
# rm1: Init
# rm2: Init
# rm3: Init
```

可以看到这些rm1和rm2的id并不一样，但是他们的state却是一样的，达到这个效果的关键就在Borg的那个`__init__`函数中。

```python
    def __init__(self):
        self.__dict__ = self.__shared_state
```

通过对`__dict__`的重载，它指向了一个类变量，这样所有的instance的`__dict__`都指向同一个变量，即`__shared_state`，通过这种方式在各个instance间共享了同样的数据。

---

## 小结

在实际工作中，我几乎从未使用过第三种方法，因为它容易让人疑惑，大部分情况下类logging的方法都是不错的解决方案。
