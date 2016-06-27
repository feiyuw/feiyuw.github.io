---
layout: post
title:  "Inject a running robotframework process"
date:   2016-06-24 15:00:00 +0800
categories: TA
---
[Robotframework](https://robotframework.org) is a widely used ATDD framework. Sometime the test suite become very large and complicated. In my company, some suites contain thousands of execution steps and last over 12 hours.
I was always questioned on such scenario:

> Hi Zhang, could you help me? My robot hung for several hours, could you tell me what it is doing?

or

> Hi, there is a very very long execution case that already run for 12 hours, but just now I find one variable in teardown was set incorrectly, I don't want to rerun the whole suite, could you help me?

These requirements were hard to met, as I had no interfaces to interract with the running robot process.

---

Fortunatelly we have [pyrasite](http://pyrasite.com/). A tool that can inject code into running python process. If you still do not how it work, take some minutes to watch the video on its homepage.

Install `pyrasite` is a piece of cake `pip install pyrasite`.

To use `pyrasite`, you should enable [ptrace](http://man7.org/linux/man-pages/man2/ptrace.2.html) first.

```sh
sudo echo 0 | sudo tee /proc/sys/kernel/yama/ptrace_scope
```

Let's write a simple test suite as a demo.

```robotframework
*** Settings ***

*** Test Cases ***
Test One
    ${testvar}    Set Variable    0
    Log    Before Sleep
    Sleep   60s
    Log    ${testvar}
```

### 1. Inject to the running robot process

1. Run the suite with `pybot` command
1. Detect the process id with some command like `ps ax | grep robot`
1. Open a terminal, and execute `pyrasite-shell <pid>`, now it will open a shell, on that shell, you can inject code into running robot process

### 2. First try - print some message on terminal

Let's verify if our injection can work or not. On `pyrasite-shell`, let's input:

```python
import sys
sys.__stderr__.write('hello, you are injected!\n')
```

![show inject log](/assets/inject-stderr.png)

> In this example, we use `sys.__stderr__` instead of `sys.stderr` as robotframework will override the default `sys.stderr` and `sys.stdout`


### 3. Interact with robot EXECUTION_CONTEXT

`Robotframework` does not provide public API to interract with its context, but the interface `robot.running.EXECUTION_CONTEXTS` can be used too (at least from version 2.7 to 3.0).
Let's see how to use this interface to get the running variables right now.

Still in `pyrasite-shell`:

```python
from robot.running import EXECUTION_CONTEXTS
ctx = EXECUTION_CONTEXTS.current
print ctx.variables.as_dict()
```

![show variables](/assets/inject-variables.png)

We get all the variables in this scope, usually from the variables, we can guess which step it is right now.

And we are not only able to view the data, but also modify and update. For example, we can update the value of one variable.

```python
ctx.variables.set_keyword('${testvar}', '1')
```

Or insert a step.

```python
from robot.running import Keyword
Keyword(name='Log', args=('inserted step', )).run(ctx)
```

Read `pyrasite` document and the code of `RobotFramework` to find more usage.
