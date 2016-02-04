---
layout: post
title:  "自己动手编写持续交付软件"
date:   2016-02-04 17:32:00
categories: Programming
---

[持续交付系统](https://en.wikipedia.org/wiki/Continuous_delivery)或者[持续集成系统](https://en.wikipedia.org/wiki/Continuous_integration)是现代软件开发不可或缺的基础设施，尤其在大型项目中，**持续交付系统**的质量和效率往往极大影响开发的质量和效率。许多大型软件项目甚至有一个或多个专门的团队来开发和维护这类系统。另外，**持续交付系统**包含软件研发中的大量一手数据，很好地挖掘这些数据，也可以得到很多有价值的信息。

市面上有很多商业的和开源的持续交付系统的解决方案，如[Bambook](https://www.atlassian.com/software/bamboo/)，[Jenkins](http://jenkins-ci.org)，与github深度集成的[Travis-CI](http://travis-ci.org)等，另外在8.0+版本的[gitlab](http://gitlab.org)里也将gitlab-ci集成了进去，作为内置的持续交付平台。

既然已经有这么多现成的工具了，那么为什么我们还要自己写一套持续交付系统呢？

1. 写一个持续交付系统很有趣
1. 可以完美适配自己的开发流程
1. 可以整合其他既有的资源，如报表系统等

那么，一个**持续交付系统**需要包含哪些基本功能呢？

1. 与版本控制工具如[subversion](http://subversion.apache.org/)、[git](http://git-scm.com/)的集成
1. 执行节点的管理
1. 执行内容的可定制化
1. 结果保存和查看
1. 反馈与通知

所以，一个**持续交付系统**至少具有：

* 一个可视化的界面（通常是Web）
* 一个数据库或文件系统（用于保存配置和执行结果）
* 一个定时任务工具（用于检测代码版本的变化和发送通知等）。

这里，笔者希望通过一个具体的例子来展示一下如何开发一个持续交付系统，这个进行中的项目的地址在[github](http://github.com/lybicat/)。

## 从界面开始

[TODO]

## 数据的保存

[TODO]

## 执行日志的存储

[TODO]

## 实时查看执行的结果

[TODO]

## 如何管理执行节点

[TODO]

## 如何集成现有的版本控制工具

[TODO]

## 如何优雅地部署和升级

[TODO]

