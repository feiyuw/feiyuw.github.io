---
layout: post
title:  "细胞自动机 - 布莱恩大脑"
date:   2019-09-17 23:59:00 +0800
categories: "Karta"
---

继[康威生命游戏](/karta/2019/07/11/game-of-life/)和[兰顿蚂蚁](/karta/2019/07/11/langton-ant-problem/)之后的第三个细胞自动机。

## 简介

Brians' brain可以看作生命游戏的一个扩展，它也是一个二维的细胞自动机，在生命游戏的基础上引入了第三个状态，规则如下：

1. 每个细胞有三种状态：
    * ready
    * firing
    * refactory
1. 细胞下一个迭代的状态由这一个迭代其自身状态和它的八个邻居的状态决定
1. 如果当前细胞为ready状态，且它的八个邻居中有两个为firing状态，则细胞变为firing状态，不然保持ready状态不变
1. 如果当前细胞为firing状态，则其下个迭代变为refactory状态
1. 如果当前细胞为refactory状态，则其下个迭代变为ready状态

## 参考代码

[brian's brain](https://github.com/feiyuw/brianbrain)

## 运行示例

借用生命游戏d3版本的界面，实现了一个简单的版本，白色表示ready、绿色表示firing、橙色表示refactory状态。

下面是一个执行示例，你也可以访问[Online Demo](/brianbrain/index.html)。

<div>
  <div id="board" width='100%'></div>
  <style>
    svg {
      width: 100%;
    }
    circle[data="2"] {
      fill: orange;
    }
    circle[data="1"] {
      fill: green;
    }
    circle[data="0"] {
      fill: white;
    }
  </style>
  <script src="//cdnjs.cloudflare.com/ajax/libs/lodash.js/4.13.1/lodash.min.js"></script>
  <script src="//cdnjs.cloudflare.com/ajax/libs/d3/4.1.1/d3.min.js"></script>
  <script src='/brianbrain/index.d3.js'></script>
  <script>
    const board = new Board('#board')
    const rows = 30
    const cols = 60
    const delay = 500
    const game = new BrianBrain(rows, cols)
    game.initBoard()
    board.render(game.getLives())
    const handler = () => {
      game.nextRound()
      board.render(game.getLives())
      setTimeout(handler, delay)
    }
    handler()
  </script>
</div>

## 终端运行示例

以下为一个golang版本在终端运行的示例：

<script id="asciicast-f1SPnBZdVPVPZAWbgnNYHOstX" src="https://asciinema.org/a/f1SPnBZdVPVPZAWbgnNYHOstX.js" async></script>
