---
layout: post
title:  "使用pandas进行数据处理"
date:   2016-02-29 15:00:00
categories: Python
---
最近开始学习pandas用来作为数据分析的入门，这里将最近的学习所得记录在这里，以作小结。

## pandas简介

[pandas](http://pandas.pydata.org/)是一个开源的数据结构化和分析工具。它的出现是为了解决Python语言在数据分析和建模方面的缺失，让工程师可以在数据采集和分析方面采用同样的语言，而不需要切换到专门的分析语言如R。

pandas依赖于[numpy](http://www.numpy.org/)，提供了一个非常有效的数据结构**DataFrame**，可以把**DataFrame**想象成Excel的工作表，它可以提供非常有效的组织和索引功能。

另外，pandas提供了API用于读写多种类型数据，包括：

* CSV
* Excel
* SQL database
* HDF5
* JSON
* 剪贴板
* ...

而借助于[Jupyter notebook](http://jupyter.org/)，pandas可以方便地进行数据的分析和可视化，如我学习用的[notebook](/assets/learn-pandas.ipynb)。

## iPython简介

[ipython](http://ipython.org/)最初是作为一个增强型的交互式python shell创建起来的，现在它已经成为python数据分析和可视化的一个不可或缺的工具。在日常的Python开发中，它的使用频率已经远超过IDE，对于笔者来说，ipython+vim作为Python开发环境非常高校。

### 安装

ipython的安装非常简单，通过pip就可以了。
通常，我们在安装完ipython后，还会安装jupyter，因为我们需要用到notebook和qtconsole之类的功能。这个中间有些模块需要另外安装，根据错误提示来就可以了。

* 安装ipython：`pip install ipython`
* 安装jupyter：`pip install jupyter`

笔者的工作电脑是Linux，如果你用的是Windows，通常你还需要安装一下pyreadline，通过`pip install pyreadline`就可以安装了。

### 使用

在一个Terminal工具（如gnome-terminal）里执行`ipython`就可以进入ipython。

执行`jupyter qtconsole`或`jupyter notebook`就可以打开qtconsole和web notebook了。比较推荐notebook的模式，因为在notebook中可以混合markdown、代码、执行结果和生成的图表，本身就是一份活的文档了。

### 特点

ipython相比默认的python shell有很多优点，至少包括：

* 自动补全，输入几个字符按Tab
* 文档查看，在模块或变量名前加？执行，如`?os`
* 执行系统命令，在命令钱加！，如`！ls`
* 执行文件，如`%run test.py`
* 内嵌显示图表（需qtconsole或notebook）

jupyter notebook作为交互式的web notebook，可以极大地提高效率，因为：

* 将代码和执行结果融合到一起
* 将生成的图表嵌入到notebook中
* 可以添加markdown用于说明和文档
* web形式便于分享

### 性能分析

ipython有内置的性能分析方法，可以方便地分析函数的性能，包括：

* %time  ==>  获得程序运行的时间
* %timeit  ==>  持续运行100万次，分析函数所用时间
* %prun  ==>  以cProfile的方式运行函数进行性能分析

除此意外，我们还可以使用[line_profiler](https://github.com/rkern/line_profiler)

* 安装：pip install line_profiler
* 在ipython中启用：在~/.ipython/profile_default/ipython_config.py中加入`c.TerminalIPythonApp.extensions = [ 'line_profiler', ]`
* 在ipython中使用`%lprun`

## pandas入门

在[利用Python进行数据分析](https://book.douban.com/subject/25779298/)这本书的引言部分有三个例子，我在pandas 0.17.1版本上把他们都实现了一下，通过这三个例子，我们可以一窥pandas数据分析的门径。

### 例子-1：分析网页请求数据

在http://1usagov.measuredvoice.com/2013/上可以下载到网页请求数据，我们选取其中一天的数据文件**usagov_bitly_data2013-05-17-1368832207**来进行分析。

这个文件的每一行都是一个json字符串，因此我们可以很方便地把它按行转换成一个list。

```python
import json
# 数据来源: http://1usagov.measuredvoice.com/2013/
with open('usagov_bitly_data2013-05-17-1368832207') as fp:
    records = map(json.loads, fp)
```

#### 使用pandas来过滤数据

我们会将上面得到的records转换成pandas的DataFrame对象，以进行后续的处理。这里我们希望通过简单的方法可以过滤出我们需要的数据。

```python
from pandas import DataFrame

data = DataFrame(records) # 以frame形式使用数据
data[(data['tz'] == '') & (data['al'] == 'en')] # filter data
# 也可以使用
data[(data.tz == '') & (data.al == 'en')]
```

上面的代码将时区（tz）为空，并且语言（al）为en的数据过滤出来。

#### 数据清洗

很多时候我们的数据集中的某些数据会存在一些字段的缺失或异常，pandas可以很方便地进行这方面的清洗和处理工作。

```python
clean_tz = data['tz'].fillna('Missing')
clean_tz[clean_tz == ''] = 'Unknown'
```

`fillna`方法将所有没有tz字段的数据的该字段值设置为**Missing**，而第二行代码则将所有tz为空的数据改为**Unknown**。

#### 按值聚合并排序

上面的clean_tz记录了所有这些数据的时区信息，我们希望知道每一个时区的访问次数并排序，这个时候我们可以采用`value_counts`方法。

```python
%matplotlib inline
clean_tz.value_counts()[:15].plot(kind='barh', figsize=(12, 5))
```

首先将matplotlib以inline方式显示，然后画图，这个图里面会显示最少访问的15个时区，如下图：

![timezone](/assets/tz15.png)

### 例子-2：使用pandas来分析电影评分数据

[TODO]

>Pivot Table(数据透视表)是一种交互式的表，可以进行某些计算，如求和与计数等。 所进行的计算与数据跟数据透视表中的排列有关。 之所以称为数据透视表，是因为可以动态地改变它们的版面布置，以便按照不同方式分析数据，也可以重新安排行号、列标和页字段。 每一次改变版面布置时，数据透视表会立即按照新的布置重新计算数据。

### 例子-3：使用pandas来分析新生儿姓名数据

[TODO]

## 参考

