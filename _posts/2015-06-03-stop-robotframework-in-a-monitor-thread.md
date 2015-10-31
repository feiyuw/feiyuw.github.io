---
layout: post
title:  "Stop RobotFramework in a monitor thread"
date:   2015-06-03 13:21:00
categories: TA
---

There is a common scenario in a TA solution, during the execution of a series of cases, there is a backend monitoring job that check some error logs or specified events of the system, when it detects some critical errors, it will stop the whole TA execution.

In `RobotFramework`, there is no built-in support for it, but we can write our own.

### First try - failed

At the first time I consider this problem (it is raised by Ju Fei), I think execute the specified keyword like `Fatal Error` in the monitoring thread can do this task. And I think writing it as a listener is a good choice.
The first version is like below:

```python
# listener.py
import thread
import os
from robot.running import EXECUTION_CONTEXTS
from robot.running import Keyword

class RaiseFatalErrorInThread:
    ROBOT_LISTENER_API_VERSION = 2

    def start_test(self, name, attrs):
        thread.start_new_thread(self._stop_execution, ())

    def _stop_execution(self):
        Keyword('Log', ('stopped by user', )).run(EXECUTION_CONTEXTS.current)
        Keyword('fatal error', ()).run(EXECUTION_CONTEXTS.current)
```

And the testcase is as below:

```robotframework
# test.robot
*** Test Cases ***
test listener in thread
    Just Test Step

*** Keywords ***
Just Test Step
    Do Error Log Monitoring
    Sleep   5min
    Log     Second step
```

Run the test with command `pybot --listener listener.RaiseFatalErrorInThread test.robot`, but it failed to generate the log.html, and the case did not stop as expected.

### Using signal - OK
Checking the source code of `robot.running.EXECUTION_CONTEXTS`, there is no thread switch, so it is impossible for this issue. But at the same time, I think about the `Stop Gracefully` feature in `RobotFramework`, it use `signal` to interrupt the execution. So I decide to follow it and send `SIGINT` signal to main process. The code is as below:

```python
# listener.py
import thread
import os
import signal
from robot.api import logger

class RaiseFatalErrorInThread:
    ROBOT_LISTENER_API_VERSION = 2

    def start_test(self, name, attrs):
        logger.info('start a new test')
        thread.start_new_thread(self._stop_execution, ())

    def _stop_execution(self):
        logger.librarylogger.warn('stopped by user')
        os.kill(os.getpid(), signal.SIGINT)
```

It is OK, but the problem is we cannot detect it is a user interrupt log or monitor interrupt log, there is no difference between them. And as a listener, it is not easy to use for a tester.

### Different signal in library - OK

The third try is putting the code in a keyword, and use `SIGUSR1`, as `SIGINT`, `SIGTERM` and `SIGALRM` are used in `RobotFramework` already. The code is as below and on gist https://gist.github.com/feiyuw/76faf6cfdf087a9a04a2

```python
# lib.py
import signal
from robot.running import EXECUTION_CONTEXTS
from robot.running import Keyword
import thread
import os

def do_error_log_monitoring():
    def _stop_execution(signum, frame):
        Keyword('fatal error', ()).run(EXECUTION_CONTEXTS.current)
    def _monitor_log():
        import time
        time.sleep(5)
        os.kill(os.getpid(), signal.SIGUSR1)
    signal.signal(signal.SIGUSR1, _stop_execution)
    thread.start_new_thread(_monitor_log, ())
```

Add library to the test case.

```robotframework
# test.robot
*** Settings ***
Library     lib.py

*** Test Cases ***
Stop test by signal
    Just Test Step

*** Keywords ***
Just Test Step
    Do Error Log Monitoring
    Sleep   5min
    Log     Second step
```

Run the case with command `pybot test.robot`, it seems to work well now.

