---
layout: post
title:  "Debug RobotFramework case in iPython"
date:   2015-06-02 18:26:00
categories: TA
---
[RobotFramework](http://robotframework.org) is an popular open source Test Automation Framework, it's written in Python and can provide a good ability of extention. It's a good choice for ATDD.

[iPython](http://ipython.org) is an enhanced shell of Python, and provide many great features like browser-based notebook, data visualization and scientific computing. It can improve your efficiency greatly.

Now it is time to combine them together.

## With a suite file
In most scenario, you have a suite file already, and based on it, you may add steps, check the result and so on.

Go through the `run` method of RobotFramework, it is as below:
```python
if not settings:
    settings = RobotSettings(options)
    LOGGER.register_console_logger(**settings.console_logger_config)
with pyloggingconf.robot_handler_enabled(settings.log_level):
    with STOP_SIGNAL_MONITOR:
        IMPORTER.reset()
        init_global_variables(settings)
        output = Output(settings)
        runner = Runner(output, settings)
        self.visit(runner)
    output.close(runner.result)
return runner.result
```

Debugging is based on the `keyword` steps, we want to execute some `keyword`, and may we will insert some from library file or iPython console.

Let's start!

### 1. Build runner
```python
from robot.api import TestSuiteBuilder

suite = TestSuiteBuilder().build('kw-driven.robot') # use existed robot suite file here

```

## Without a suite file
If you start to implement a new test suite, you can follow the following steps.

[TODO]...
