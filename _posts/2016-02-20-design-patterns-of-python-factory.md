---
layout: post
title:  "Python设计模式 - 工厂模式"
date:   2016-02-20 15:00:00
categories: Python
---
**工厂模式**可能是最为著名的设计模式了，我记得自己在还没有看过任何有关设计模式的书或文章的时候，就对自己的代码进行过一些重构，而这些重构绝大多数用的都是类似工厂模式的方法。

维基百科上有关**工厂模式**的定义如下：

>工厂方法模式（英语：Factory method pattern）是一种实现了“工厂”概念的面向对象设计模式。就像其他创建型模式一样，它也是处理在不指定对象具体类型的情况下创建对象的问题。工厂方法模式的实质是“定义一个创建对象的接口，但让实现这个接口的类来决定实例化哪个类。工厂方法让类的实例化推迟到子类中进行。”

---

### 普通的代码

在Python里面，**工厂模式**通常以一个工厂方法的面貌出现，比如说一个parse方法，用来解析不同类型的文件，如果不用工厂模式，我们的代码可能是这样的：

```python
def parse_xml(file_path):
    pass
    # ...

def parse_csv(file_path):
    pass
    # ...

def parse_ims2(file_path):
    pass
    # ...

# ...
```

所有这些函数`parse_xml`，`parse_csv`，`parse_ims2`等都会暴露在外部，调用的时候，我们会根据文件类型调用不同的函数。往往我们还需要根据文件的类型进行if...else判断。

---

### 引入工厂方法

而事实上，我们所需要的只是一个文件解析方法，把文件丢给它，它把解析好的内容返回给我们就行了。所以，这个时候，我们引入一个新的函数parse：

```python
def parse(file_path):
    _, ext = os.path.splitext(file_path)
    if ext == '.xml':
        return _parse_xml(file_path)
    elif ext == '.csv':
        return _parse_csv(file_path)
    elif ext == '.ims2':
        return _parse_ims2(file_path)

    raise RuntimeError('unknown file type "%s"' % ext)
```

对于调用方来说，只需要调用parse函数就可以了，它成了唯一的对外接口，而对于不同文件格式的parse则隐藏在内部了。`parse`就是我们的工厂方法。当然我们没有像它的定义描述的那样返回一个parse的实例，而是直接返回parse的结果了，因为Python在语法层面要比Java等语言灵活得多，这里就可以简化了。

---

### 继续抽象

上面的代码有一个不足，每当我们需要增加一种类型的文件的解析的时候，我们都需要修改parse函数的代码，这是一个不太安全的设计。我们希望当有新的文件类型支持的时候，只需要添加文件解析的方法就可以了，而不用修改parse，简而言之，就是**增加而不修改**。

这里我们借助Python的globals方法来实现这一点：

```python
def parse(file_path):
    _, ext = os.path.splitext(file_path)
    parser_name = '_parse_' + ext[1:]
    if parser_name in globals():
        return globals()[parser_name](file_path)

    raise RuntimeError('unknown file type "%s"' % ext)
```

---

### 真实的例子

在[robotframework](https://github.com/robotframework/robotframework)的代码中有很多使用工厂模式的例子，比如：

```python
# robot.result.resultbuilder

def ExecutionResult(*sources, **options):
    """Factory method to constructs :class:`~.executionresult.Result` objects.

    :param sources: Path(s) to the XML output file(s).
    :param options: Configuration options.
        Using ``merge=True`` causes multiple results to be combined so that
        tests in the latter results replace the ones in the original. Other
        options are passed directly to the :class:`ExecutionResultBuilder`
        object used internally.
    :returns: :class:`~.executionresult.Result` instance.

    Should be imported by external code via the :mod:`robot.api` package.
    See the :mod:`robot.result` package for a usage example.
    """
    if not sources:
        raise DataError('One or more data source needed.')
    if options.pop('merge', False):
        return _merge_results(sources[0], sources[1:], options)
    if len(sources) > 1:
        return _combine_results(sources, options)
    return _single_result(sources[0], options)

# robot.parsing.model

def TestData(parent=None, source=None, include_suites=None,
             warn_on_skipped=False):
    """Parses a file or directory to a corresponding model object.

    :param parent: (optional) parent to be used in creation of the model object.
    :param source: path where test data is read from.
    :returns: :class:`~.model.TestDataDirectory`  if `source` is a directory,
        :class:`~.model.TestCaseFile` otherwise.
    """
    if os.path.isdir(source):
        return TestDataDirectory(parent, source).populate(include_suites,
                                                          warn_on_skipped)
    return TestCaseFile(parent, source).populate()
```

