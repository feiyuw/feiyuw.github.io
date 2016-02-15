---
layout: post
title:  "自己动手编写持续交付系统（二）"
date:   2016-02-15 07:00:00
categories: Programming
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

---

## 1. 定义一个Plan

一个**plan**需要包含哪些信息呢？在上一章中我们简单描述了一下，这里，我们继续探讨这个问题。

[TODO]

---

## 2. 管理测试用例

[TODO]

---

## 3. 管理测试资源

[TODO]

---

## 4. 管理执行步骤

[TODO]

---

## 5. 解决冲突

[TODO]

