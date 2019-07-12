---
layout: post
title:  "细胞自动机 - 康威生命游戏"
date:   2019-07-11 22:00:00 +0800
categories: "Karta"
---

近来整理之前的笔记，无意中看到了早期写的game of life程序，发现其代码大都已经老旧，如Python版本之前还是基于2.7实现的，现在已普遍切换到3.7了，ReactJS版本的依赖版过劳，也有一些问题，于是干脆重新清理了一下。在这个过程中，再次玩了下经典的自动细胞机 - 康威生命游戏。

## 细胞自动机

细胞自动机的各类实验出于一种对寻找大自然终极简单规律的理想。就像很多优美的数学公式那样简单。最早的细胞自动机有冯.诺依曼在20世纪50年代左右提出，用来模拟生物细胞的自我复制。20世纪70年代，英国数学家康威发明了生命游戏，从此细胞自动机开始被很多人所认识。

细胞自动机的特点是，在一个N维的世界，每个格子代表一个生命，它有有限种状态，每个格子时间t的状态由时间t-1的状态唯一决定，所有的格子按照统一的规则演进。

## 生命游戏

康威定义的规则非常简单，在一个二维的世界里，每个细胞只有生和死两种状态，它下一轮的状态取决于这一轮周围的8个邻居的状态，简单规则如下：

* 当前细胞为存活状态时，当周围的存活细胞低于2个时（不包含2个），该细胞变成死亡状态。（模拟生命数量稀少）
* 当前细胞为存活状态时，当周围有2个或3个存活细胞时，该细胞保持原样。
* 当前细胞为存活状态时，当周围有超过3个存活细胞时，该细胞变成死亡状态。（模拟生命数量过多）
* 当前细胞为死亡状态时，当周围有3个存活细胞时，该细胞变成存活状态。（模拟繁殖）

## 运行示例

下面是一个执行示例，你也可以访问[Online Demo](/gameoflife/index.d3.html)。

<div>
  <div id="board" width='100%'></div>
  <style>
    circle[data="1"] {
      fill: green;
    }
    circle[data="0"] {
      fill: white;
    }
  </style>
  <script src="//cdnjs.cloudflare.com/ajax/libs/lodash.js/4.13.1/lodash.min.js"></script>
  <script src="//cdnjs.cloudflare.com/ajax/libs/d3/4.1.1/d3.min.js"></script>
  <script src='/gameoflife/gol.d3.js'></script>
  <script>
    const board = new Board('#board')
    const rows = 30
    const cols = 60
    const delay = 500
    const game = new GameOfLife(rows, cols)

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
