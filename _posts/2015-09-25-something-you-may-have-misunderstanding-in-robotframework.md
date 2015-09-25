---
layout: post
title:  "Something you may have misunderstanding in robotframework"
date:   2015-09-25 10:26:00
categories: TA
---

There are many [RobotFramework](http://robotframework.org) users, with huge amount of libraries, keywords, resources and test suites. But I guess there should be something you misunderstand of this framework. I listed some of them here, and if you have others, that should be very helpful if you post them in the comment.

### 1. Keyword `Run Keyword and Ignore Error` does not ignore all errors

Sometimes we need to ignore the keyword execution error, such as an extra IP configuration. In most scenario, we will use BuiltIn keyword `run keyword and ignore error`. But does it ignore all errors?

The answer is **No**. Below is the source code of this keyword.

```python
try:
    return 'PASS', self.run_keyword(name, *args)
except ExecutionFailed as err:
    if err.dont_continue:
        raise
    return 'FAIL', unicode(err)
```

We can see, when `err.dont_continue` is `True`, this keyword will fail. And let's see the code of `ExecutionFailed`.

```python
def dont_continue(self):
    return self.timeout or self.syntax or self.exit
```

When `timeout` or `syntax` or `exit`, it will return `True`.

OK, until now, we know there are three scenarios `run keyword and ignore error` will fail: Timeout, Syntax Error, Fatal Exception.


### 2. Keyword `Wait Until Keyword Succeeds` will fail before timeout

Similar as `run keyword and ignore error`, `wait until keyword succeeds` is not always reliable. Sometimes it will fail before timeout reached.

Below is the while loop of this keyword.

```python
while True:
    try:
        return self.run_keyword(name, *args)
    except ExecutionFailed as err:
        if err.dont_continue:
            raise
        count -= 1
        if time.time() > maxtime > 0 or count == 0:
            raise AssertionError("Keyword '%s' failed after retrying "
                                 "%s. The last error was: %s"
                                 % (name, message, err))
        self._sleep_in_parts(retry_interval)
```

It raise in `err.dont_continue` too. Same as the previous one.


### 3. RobotFramework will continue do the execution even if some keywords fail in teardown

RobotFramework does not fail the keyword immediately after running it. Actually, it catch the exception, and do different actions on different scenarios.

Let's see the code.

```python
# robot.running.keywordrunner.KeywordRunner

def run_keywords(self, keywords):
    errors = []
    for kw in keywords:
        try:
            self.run_keyword(kw)
        except ExecutionPassed as exception:
            exception.set_earlier_failures(errors)
            raise exception
        except ExecutionFailed as exception:
            errors.extend(exception.get_errors())
            if not exception.can_continue(self._context.in_teardown,
                                          self._templated,
                                          self._context.dry_run):
                break
    if errors:
        raise ExecutionFailures(errors)
```

Here we can see, when the `exception.can_continue(self._context.in_teardown)` is `True`, the execution will not be interruptted. And `can_continue` is as below:

```python
def can_continue(self, teardown=False, templated=False, dry_run=False):
    if dry_run:
        return True
    if self.dont_continue and not (teardown and self.syntax):
        return False
    if teardown or templated:
        return True
    return self.continue_on_failure
```

So everything is clear, when the keyword is in the teardown, it will not fail immediately.

### 4. `Stop Gracefully` in `jybot` will not stop immediately

`Stop Gracefully` is an amazing feature of `RobotFramework`. But it does not stop current keyword immediately in `jybot`, why?

Let's see the code.

```python
# robot.running.signalhandler._StopSignalMonitor

def __call__(self, signum, frame):
    self._signal_count += 1
    LOGGER.info('Received signal: %s.' % signum)
    if self._signal_count > 1:
        sys.__stderr__.write('Execution forcefully stopped.\n')
        raise SystemExit()
    sys.__stderr__.write('Second signal will force exit.\n')
    if self._running_keyword and not sys.platform.startswith('java'):
        self._stop_execution_gracefully()
```

Obviously, when the platform is `java`, it will not stop.

