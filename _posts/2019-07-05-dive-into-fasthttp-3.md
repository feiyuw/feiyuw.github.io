---
layout: post
title:  "fasthttp高性能之道（三）"
date:   2019-07-05 22:00:00 +0800
categories: "GoodCode"
---

在上两篇[fasthttp高性能之道（一）](/goodcode/2019/04/20/dive-into-fasthttp-1/)和[fasthttp高性能之道（二）](/goodcode/2019/05/25/dive-into-fasthttp-2/)中我们分别介绍了fasthttp在HTTP server和client两个方面的一些实现特点，以及它在应对高并发时候的一些策略。这一次，我们将深入到它的一些基础模块，看看为了提高并发，降低内存，它都在标准库的基础上做了哪些改进。

## bytebufferpool


## tcpdialer


