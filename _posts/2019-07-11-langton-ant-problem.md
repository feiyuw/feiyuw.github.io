---
layout: post
title:  "细胞自动机 - 兰顿蚂蚁"
date:   2019-07-11 23:59:00 +0800
categories: "Karta"
---

改完[细胞自动机 - 康威生命游戏](/karta/2019/07/11/game-of-life/)觉得不过瘾，再写一个细胞自动机玩玩。

## 简介

兰顿蚂蚁是另一个简单又有名的细胞自动机，诞生于1986年，同样是二维的，规则非常简单：

> 在平面上的正方形格被填上黑色或白色。在其中一格正方形有一只“蚂蚁”。它的头部朝向上下左右其中一方。

> 若蚂蚁在白格，右转90度，将该格改为黑格，向前移一步；
> 若蚂蚁在黑格，左转90度，将该格改为白格，向前移一步。

## 运行示例

借用生命游戏d3版本的界面，实现了一个简单的版本，用红色代表蚂蚁，绿色表示存活，白色表示死亡。

下面是一个执行示例，你也可以访问[Online Demo](/langtonant/index.html)。

<div>
  <div style="float:right"><label>Steps: </label><span id="steps">0</span><div>
  <div id="board" width='100%'></div>
  <style>
    svg {
      width: 100%;
    }
    circle[data="2"] {
      fill: red;
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
  <script src='/langtonant/index.d3.js'></script>
  <script>
    const board = new Board('#board')
    const steps = document.getElementById("steps")
    let stepCount = 0

    const rows = 40
    const cols = 80
    const delay = 500
    const game = new LangtonAnt(rows, cols, 0.0)

    game.initBoard()
    board.render(game.getLives())

    const handler = () => {
      game.nextRound()
      stepCount++
      board.render(game.getLives())
      steps.innerText = stepCount
      intervalEvt = setTimeout(handler, delay)
    }

    handler()
  </script>
</div>
