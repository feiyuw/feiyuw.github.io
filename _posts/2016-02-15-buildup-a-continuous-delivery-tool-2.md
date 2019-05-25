---
layout: post
title:  "自己动手编写持续交付系统（二）"
date:   2016-02-15 15:00:00
categories: "DevOps"
---
在[自己动手编写持续交付系统（一）](http://blog.zhangyu.so/programming/2016/02/04/buildup-a-continuous-delivery-tool-1/)中，我们从基本的功能需求出发，构建了一个拥有基本功能的持续交付系统的骨架。

这一章的一开始，我们先看一下这个应用系统的简单架构。

![architecture]({{ site.url }}/assets/cd-arch.png)

从上图可以看到：

* 整个系统基本可以分为以下几部分：
    * Web UI前端
    * RESTful服务端
    * socket.io服务端
    * 一个类cron的轮询系统
    * Agent服务（作为一个socket.io客户端）
    * HDFS Log服务
* Web UI与RESTful服务端通过HTTP进行通讯
* Agent与socket.io服务端通过websocket协议进行通讯
* Agent与HDFS Log服务通过HTTP通讯

骨架已经搭好，是时候添砖加瓦了！

## 1. 定义一个Plan

一个**plan**需要包含哪些信息呢？在上一章中我们简单描述了一下，这里，我们继续探讨这个问题。

作为一个Plan，某种程度上类似于Jenkins的Job，它至少需要有执行步骤，如果什么都不做，这个Plan干什么呢？所以，基本信息包括：

* Plan的名字（大部分情况下是唯一的）
* Plan的执行步骤（我们称之为Action）
* Plan的所有者
* Plan的创建和修改时间

另外，大部分情况下，我们需要知道这个Plan是针对哪个Repository的，以及它的测试方式是什么，如果是像gitlab或者travis-ci那样针对一个具体的repository的前期测试的话，可能提供一个yaml文件是不错的办法。

Lybica暂时不打算重复做这些事情，我们针对的是更高维度的测试和集成，因此在这里可以选择的变成了：

* 测试用例的Repository
* 测试用例的过滤规则
* 用到的Resource
* 额外的变量Variable

另外，除了这些，还有可能包括：

* 成功/失败后的通知内容和配置
* 触发源（Trigger）

## 2. 管理测试用例

大多数情况下，测试用例和代码一样，都是存在于SCM仓库里的，如Git，Subversion等。而对于高维度的测试，checkout这些测试用例可能要花费比较多的时间，以笔者经历过的一个项目为例，完整地checkout一份测试用例Repo需要超过30分钟的时间，如果每次都这么做，开销有点太大了。

另外，在高维度的测试里面，关注的是用户功能，所以一个测试用例可能会跨多个软件模块，而且本身有层次关系，比如稳定性测试用例就会等功能测试都OK了才开始。作为一个Plan来说，是不需要运行所有的测试用例的。这样，就需要有一个过滤机制来过滤出与Plan相关的测试用例。

因此在界面配置上，需要提供：

* 测试用例Repository的选择
* 过滤条件的设置

而在Platform层面，则需要一个Repository上测试用例的过滤机制，背后可能还有一个从Repository到数据库的同步机制。

## 3. 管理测试资源

[TODO]

## 4. 管理执行步骤

[TODO]

## 5. 解决冲突

[TODO]

