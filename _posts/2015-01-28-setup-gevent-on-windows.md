---
layout: post
title:  "在windows平台上安装gevent和gevent-socketio的简单方法"
date:   2015-01-28 18:05:03
categories: Python
---
很久没有在windows平台工作了, 最近写的一个小工具要部署到windows平台上, 而我在代码里面使用了gevent和gevent-socketio. 这下问题来了.

在Linux平台上安装gevent就是几个命令的事情, 有了gcc之类的编译环境很快就可以搞定, 但是windows上要编译什么的就麻烦太多了, 笔者尝试了VS2012 community版本, 编译libevent还是有各种问题, 而我要想把这套方案部署到其他人的机器上, 难度就可想而知了.

在尝试了一天的编译之后, 笔者决定转换思路, 寻找已经编译好的代码, 最好是一键安装的安装包.

gevent依赖greenlet, 从源码安装的话, 还需要cython的支持和编译libev. 而greenlet在windows平台上则需要编译libevent. 因此我们至少需要:

* libevent的dll文件
* greenlet的免编译安装包
* gevent的免编译安装包

## libevent.dll

通过万能的google, 找到了http://www.dll-found.com/libevent-2-0-5.dll_download.html (注意: 该地址可能需要翻墙才能访问)
笔者的测试环境是windows 7 64bit, 安装了32bit的Python 2.7版本. 因此将下载到的libevent-2-0-5.dll复制到C:\Windows\SysWOW64目录下即可.

## greenlet

在另一个神奇的网站http://www.lfd.uci.edu/~gohlke/pythonlibs/ 可以找到greenlet和gevent的whl包. 根据python版本下载下来, 我的是greenlet‑0.4.5‑cp27‑none‑win32.whl

通过`pip install greenlet‑0.4.5‑cp27‑none‑win32.whl`安装即可

## gevent

步骤同greenlet, 我下的是gevent‑1.0.1‑cp27‑none‑win32.whl

通过`pip install gevent‑1.0.1‑cp27‑none‑win32.whl`安装

## gevent-socketio

这个没有什么需要编译的地方, 直接`pip install gevent-socketio`安装就OK了


## 共享

提供一个我使用到的文件的压缩包, 在百度网盘上, 需要的人自己去下吧. [socket on windows](http://pan.baidu.com/s/1dDGRn49)
