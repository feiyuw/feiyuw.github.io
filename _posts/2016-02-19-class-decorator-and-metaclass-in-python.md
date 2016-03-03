---
layout: post
title:  "在Python中使用class decorator和metaclass"
date:   2016-02-19 15:00:00
categories: Python
---
在Python语言中class层面的decorator和metaclass可能是绝大部分Pythoner都没有使用过的黑魔法了。笔者自问，到目前为止还没有在任何一个部署的代码中使用过这两个语言特性，对metaclass的需求也仅仅一次，最后还是用别的手段解决了。

但这两个语言特性是否就是多余的呢？自省一下，笔者不禁觉得当时对metaclass的那个需求最好的手段可能还是用metaclass，而之所以没有用，是因为自己对这个特性很陌生，怕驾驭不了。

那么我们常见的Python库中有没有这两者的使用案例呢？

## class decorator

关于常见的Python库中使用class decorator的例子我没有找到，如果有人知道哪个库里这么使用了，麻烦留言告诉我，谢谢！～

我们使用函数层面的decorator很常见，它以一个函数作为输入，返回一个新的函数，而class层面的decorator其实也类似，它以一个类作为输入，返回一个新的类。见下面的例子：

```python
def add_len_method(cls):
    class NewCls(cls):
        @property
        def length(self):
            if hasattr(self, '__len__'):
                return len(self)
            return 0

    return NewCls

@add_len_method
class User(object):
    pass

u = User()
u.length # 0
```

* 这个decorator给类加上了一个length属性
* 如果这个类已经有`__len__`方法，就用`__len__`方法的结果
* 如果`__len__`方法没有定义，就返回0

然并卵，对于这种需求，写个基类用mixin模式不是更好吗？

## metaclass

使用metaclass的一个著名的例子就是[Django](https://www.djangoproject.com/)ORM的ModelBase，让我们看看它是怎么做的。

当我们使用django定义model的时候，我们一般是这么写的：

```python
from django.db import models

class User(models.Model):
    name = models.CharField(max_length=32)
    age = models.IntegerField()
    addr = models.CharField(max_length=512)
```

而在使用的时候，我们是这样的：

```python
user = User()

user.name = 'zhang3'
user.age = 32
user.addr = 'Xihu, Hangzhou, China'

user.save()
```

我们并没有在User的`__init__`里面定义这些attribute，像user.name这样也不是models.CharField类型，那么显然在真正生成class的时候它被改变了。

那么，metaclass究竟是什么呢？简单来说，就是用来生成class的class。当一个class定义了`__metaclass__`之后，生成这个class的时候就会使用`__metaclass__`，不然的话就用用它的父类的`__metaclass__`或者module的`__metaclass__`，直至type。

这个有点绕，关于metaclass的原理不太理解的，可以看**参考**的两个帖子。这里只是简单说一下：

>Python里面所有的东西都是对象，包括整数，字符串等等，它们都来源于type，你可以通过__class__属性来看到这点。这里当然也包括class本身，所以class本身也是一个对象，生成它的那个对象就是metaclass。

关于int，string都是对象的，可以看下面的代码：

```python
i = 5
s = 'abc'
i.__class__ # int
s.__class__ # str
i.__class__.__class__ # type
s.__class__.__class__ # type
```

怎么实现一个metaclass呢？我们用一个例子来说明：

* 我们希望所有的类的函数名字都是小写的
* 不能改变这些类的现有的实现
* 原来的不规范的函数名字不能被调用

```python
class Meta(type):
    def __new__(cls, name, bases, attrs):
        return type(name, bases, {k.lower(): v for k,v in attrs.items()})

class User(object):
    __metaclass__ = Meta

    @property
    def Age(self):
        return 22

    @property
    def Name(self):
        return 'Tom Jerry'

u = User()
hasattr(u, 'Age') # False
hasattr(u, 'age') # True
hasattr(u, 'Name') # False
hasattr(u, 'name') # True
```

可以看到我们把User类篡改了，不再有`Age`属性，而是变成了`age`。

现在，我们来看一下`django.db.models.Model`是怎么实现的。

```python
# django/db/models/base.py

class Model(six.with_metaclass(ModelBase)):
    _deferred = False
    #...

class ModelBase(type):
    """
    Metaclass for all models.
    """
    def __new__(cls, name, bases, attrs):
        super_new = super(ModelBase, cls).__new__

        # Also ensure initialization is only performed for subclasses of Model
        # (excluding Model class itself).
        parents = [b for b in bases if isinstance(b, ModelBase)]
        # ...

# django/utils/six.py

def with_metaclass(meta, *bases):
    """Create a base class with a metaclass."""
    # This requires a bit of explanation: the basic idea is to make a dummy
    # metaclass for one level of class instantiation that replaces itself with
    # the actual metaclass.
    class metaclass(meta):

        def __new__(cls, name, this_bases, d):
            return meta(name, bases, d)
    return type.__new__(metaclass, 'temporary_class', (), {})
```

## 小结

99%的情况下你都不需要用class decorator和metaclass，除非你真有这个需求，而且明确自己在做什么。所以，哪怕不知道这些，其实也没有什么。。。（我这是在写的啥？）

送上Tim Peters（Zen Of Python的作者）的一段话：

>Metaclasses are deeper magic that 99% of users should never worry about. If you wonder whether you need them, you don't (the people who actually need them know with certainty that they need them, and don't need an explanation about why).

*Python Guru Tim Peters*

## 参考

* [what is a metaclass in python (stackoverflow)](http://stackoverflow.com/questions/100003/what-is-a-metaclass-in-python)
* [元编程](http://pycon.b0.upaiyun.com/ppt/shell909090-meta-class.html)
