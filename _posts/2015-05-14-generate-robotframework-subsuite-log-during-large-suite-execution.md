---
layout: post
title:  "Generate RobotFramework sub-suite log during large suite execution"
date:   2015-05-14 15:02:00
categories: TA
---
[RobotFramework](http://robotframework.org) is an great Test Automation Framework, it has good flexbility and extensibility. But logging mechinisam is not very good, especially when you execute a large suite that contain hundreds of sub suites, you cannot view the executed suite log before the whole one finished.

But the good news is RobotFramework is open source, and it has good listener support, it may be possible to get it yourself.
Yesterday, Ju Fei said he wish to do it to me, although I thought about it at least months ago, I didn't write even one byte of code of it.
It's his push to let me start this work and write this blog, he also send the first version of code to me. Thanks him very much!

As we know, during the execution of RobotFramework suite, it will write the log into `output.xml`, when the execution finished, it will generate `log.html` and `report.html` based on that `output.xml`. The easiest way may be generating a suite specified `output.xml` before it is started, and generate `log.html` after it finished.

Let's find out the solution of `output.xml` logging in RobotFramework firstly.

In the module `robot.output.xmllogger`, we find the class `XmlLogger`. The code is as below:

```python
class XmlLogger(ResultVisitor):

    def __init__(self, path, log_level='TRACE', generator='Robot'):
        self._log_message_is_logged = IsLogged(log_level)
        self._error_message_is_logged = IsLogged('WARN')
        self._writer = self._get_writer(path, generator)
        self._errors = []

    def _get_writer(self, path, generator):
        if not path:
            return NullMarkupWriter()
        try:
            writer = XmlWriter(path, encoding='UTF-8')
        except EnvironmentError as err:
            raise DataError("Opening output file '%s' failed: %s" %
                            (path, err.strerror))
        writer.start('robot', {'generator': get_full_version(generator),
                               'generated': get_timestamp()})
        return writer

    def close(self):
        self.start_errors()
        for msg in self._errors:
            self._write_message(msg)
        self.end_errors()
        self._writer.end('robot')
        self._writer.close()

    def set_log_level(self, level):
        return self._log_message_is_logged.set_level(level)

    def message(self, msg):
        if self._error_message_is_logged(msg.level):
            self._errors.append(msg)

    def log_message(self, msg):
        if self._log_message_is_logged(msg.level):
            self._write_message(msg)

    def _write_message(self, msg):
        attrs = {'timestamp': msg.timestamp or 'N/A', 'level': msg.level}
        if msg.html:
            attrs['html'] = 'yes'
        self._writer.element('msg', msg.message, attrs)

    def start_keyword(self, kw):
        attrs = {'name': kw.name, 'type': kw.type}
        if kw.timeout:
            attrs['timeout'] = unicode(kw.timeout)
        self._writer.start('kw', attrs)
        self._writer.element('doc', kw.doc)
        self._write_list('arguments', 'arg', (unic(a) for a in kw.args))

    def end_keyword(self, kw):
        self._write_status(kw)
        self._writer.end('kw')

    def start_test(self, test):
        attrs = {'id': test.id, 'name': test.name}
        if test.timeout:
            attrs['timeout'] = unicode(test.timeout)
        self._writer.start('test', attrs)

    def end_test(self, test):
        self._writer.element('doc', test.doc)
        self._write_list('tags', 'tag', test.tags)
        self._write_status(test, {'critical': 'yes' if test.critical else 'no'})
        self._writer.end('test')

    def start_suite(self, suite):
        attrs = {'id': suite.id, 'name': suite.name}
        if suite.source:
            attrs['source'] = suite.source
        self._writer.start('suite', attrs)

    def end_suite(self, suite):
        self._writer.element('doc', suite.doc)
        self._writer.start('metadata')
        for name, value in suite.metadata.items():
            self._writer.element('item', value, {'name': name})
        self._writer.end('metadata')
        self._write_status(suite)
        self._writer.end('suite')

    def start_statistics(self, stats):
        self._writer.start('statistics')

    def end_statistics(self, stats):
        self._writer.end('statistics')

    def start_total_statistics(self, total_stats):
        self._writer.start('total')

    def end_total_statistics(self, total_stats):
        self._writer.end('total')

    def start_tag_statistics(self, tag_stats):
        self._writer.start('tag')

    def end_tag_statistics(self, tag_stats):
        self._writer.end('tag')

    def start_suite_statistics(self, tag_stats):
        self._writer.start('suite')

    def end_suite_statistics(self, tag_stats):
        self._writer.end('suite')

    def visit_stat(self, stat):
        self._writer.element('stat', stat.name,
                             stat.get_attributes(values_as_strings=True))

    def start_errors(self, errors=None):
        self._writer.start('errors')

    def end_errors(self, errors=None):
        self._writer.end('errors')

    def _write_list(self, container_tag, item_tag, items):
        self._writer.start(container_tag)
        for item in items:
            self._writer.element(item_tag, item)
        self._writer.end(container_tag)

    def _write_status(self, item, extra_attrs=None):
        attrs = {'status': item.status, 'starttime': item.starttime or 'N/A',
                 'endtime': item.endtime or 'N/A'}
        if not (item.starttime and item.endtime):
            attrs['elapsedtime'] = str(item.elapsedtime)
        if extra_attrs:
            attrs.update(extra_attrs)
        self._writer.element('status', item.message, attrs)
```
What do you find from the code? If you are puzzled, let's read the document of listener interfaces.

| *Method*      |   *Arguments*    |   *Attributes / Explanation*                     |
| ------------- | -----------------| -------------------------------------------------|
| start_suite   | name, attributes | Keys in the attribute dictionary:                |
|               |                  |                                                  |
|               |                  | * id: suite id. 's1' for top level suite, 's1-s1'|
|               |                  |   for its first child suite, 's1-s2' for second  |
|               |                  |   child, and so on. (new in 2.8.5)               |
|               |                  | * longname: suite name including parent suites   |
|               |                  | * doc: test suite documentation                  |
|               |                  | * metadata: dictionary/map containing `free test |
|               |                  |   suite metadata`_                               |
|               |                  | * source: absolute path of the file/directory    |
|               |                  |   test suite was created from (new in 2.7)       |
|               |                  | * suites: names of suites directly in this suite |
|               |                  |   as a list of strings                           |
|               |                  | * tests: names of tests directly in this suite   |
|               |                  |   as a list of strings                           |
|               |                  | * totaltests: total number of tests in this suite|
|               |                  |   and all its sub-suites as an integer           |
|               |                  | * starttime: execution start time                |
+---------------+------------------+--------------------------------------------------+
| end_suite     | name, attributes | Keys in the attribute dictionary:                |
|               |                  |                                                  |
|               |                  | * id: suite id. 's1' for top level suite, 's1-s1'|
|               |                  |   for its first child suite, 's1-s2' for second  |
|               |                  |   child, and so on. (new in 2.8.5)               |
|               |                  | * longname: test suite name including parents    |
|               |                  | * doc: test suite documentation                  |
|               |                  | * metadata: dictionary/map containing `free test |
|               |                  |   suite metadata`_                               |
|               |                  | * source: absolute path of the file/directory    |
|               |                  |   test suite was created from (new in 2.7)       |
|               |                  | * starttime: execution start time                |
|               |                  | * endtime: execution end time                    |
|               |                  | * elapsedtime: execution time in milliseconds    |
|               |                  |   as an integer                                  |
|               |                  | * status: either `PASS` or `FAIL`                |
|               |                  | * statistics: suite statistics (number of passed |
|               |                  |   and failed tests in the suite) as a string     |
|               |                  | * message: error message if the suite setup or   |
|               |                  |   teardown has failed, empty otherwise           |
+---------------+------------------+--------------------------------------------------+
| start_test    | name, attributes | Keys in the attribute dictionary:                |
|               |                  |                                                  |
|               |                  | * id: test id in format like 's1-s2-t2', where   |
|               |                  |   beginning is parent suite id and last part     |
|               |                  |   shows test index in that suite (new in 2.8.5)  |
|               |                  | * longname: test name including parent suites    |
|               |                  | * doc: test case documentation                   |
|               |                  | * tags: test case tags as a list of strings      |
|               |                  | * critical: `yes` or `no` depending              |
|               |                  |   is test considered critical or not             |
|               |                  | * template: contains the name of the template    |
|               |                  |   used for the test. If the test is not templated|
|               |                  |   it will be an empty string                     |
|               |                  | * starttime: execution start time                |
+---------------+------------------+--------------------------------------------------+
| end_test      | name, attributes | Keys in the attribute dictionary:                |
|               |                  |                                                  |
|               |                  | * id: test id in format like 's1-s2-t2', where   |
|               |                  |   beginning is parent suite id and last part     |
|               |                  |   shows test index in that suite (new in 2.8.5)  |
|               |                  | * longname: test name including parent suites    |
|               |                  | * doc: test case documentation                   |
|               |                  | * tags: test case tags as a list of strings      |
|               |                  | * critical: `yes` or `no` depending              |
|               |                  |   is test considered critical or not             |
|               |                  | * template: contains the name of the template    |
|               |                  |   used for the test. If the test is not templated|
|               |                  |   it will be an empty string                     |
|               |                  | * starttime: execution start time                |
|               |                  | * endtime: execution end time                    |
|               |                  | * elapsedtime: execution time in milliseconds    |
|               |                  |   as an integer                                  |
|               |                  | * status: either `PASS` or `FAIL`                |
|               |                  | * message: status message, normally an error     |
|               |                  |   message or an empty string                     |
+---------------+------------------+--------------------------------------------------+
| start_keyword | name, attributes | Keys in the attribute dictionary:                |
|               |                  |                                                  |
|               |                  | * type: string `Keyword` for normal              |
|               |                  |   keywords and `Test Setup`, `Test               |
|               |                  |   Teardown`, `Suite Setup` or `Suite             |
|               |                  |   Teardown` for keywords used in suite/test      |
|               |                  |   setup/teardown                                 |
|               |                  | * doc: keyword documentation                     |
|               |                  | * args: keyword's arguments as a list of strings |
|               |                  | * starttime: execution start time                |
+---------------+------------------+--------------------------------------------------+
| end_keyword   | name, attributes | Keys in the attribute dictionary:                |
|               |                  |                                                  |
|               |                  | * type: same as with `start_keyword`             |
|               |                  | * doc: keyword documentation                     |
|               |                  | * args: keyword's arguments as a list of strings |
|               |                  | * starttime: execution start time                |
|               |                  | * endtime: execution end time                    |
|               |                  | * elapsedtime: execution time in milliseconds    |
|               |                  |   as an integer                                  |
|               |                  | * status: either `PASS` or `FAIL`                |
+---------------+------------------+--------------------------------------------------+
| log_message   | message          | Called when an executed keyword writes a log     |
|               |                  | message. `message` is a dictionary with          |
|               |                  | the following keys:                              |
|               |                  |                                                  |
|               |                  | * message: the content of the message            |
|               |                  | * level: `log level`_ used in logging the message|
|               |                  | * timestamp: message creation time, format is    |
|               |                  |   `YYYY-MM-DD hh:mm:ss.mil`                      |
|               |                  | * html: string `yes` or `no` denoting            |
|               |                  |   whether the message should be interpreted as   |
|               |                  |   HTML or not                                    |
+---------------+------------------+--------------------------------------------------+
| message       | message          | Called when the framework itself writes a syslog_|
|               |                  | message. `message` is a dictionary with          |
|               |                  | same keys as with `log_message` method.          |
+---------------+------------------+--------------------------------------------------+
| output_file   | path             | Called when writing to an output file is         |
|               |                  | finished. The path is an absolute path to the    |
|               |                  | file.                                            |
+---------------+------------------+--------------------------------------------------+
| log_file      | path             | Called when writing to a log file is             |
|               |                  | finished. The path is an absolute path to the    |
|               |                  | file.                                            |
+---------------+------------------+--------------------------------------------------+
| report_file   | path             | Called when writing to a report file is          |
|               |                  | finished. The path is an absolute path to the    |
|               |                  | file.                                            |
+---------------+------------------+--------------------------------------------------+
| debug_file    | path             | Called when writing to a debug file is           |
|               |                  | finished. The path is an absolute path to the    |
|               |                  | file.                                            |
+---------------+------------------+--------------------------------------------------+
| close         |                  | Called after all test suites, and test cases in  |
|               |                  | them, have been executed.                        |
+---------------+------------------+--------------------------------------------------+

Have you noticed they have the same function names, like `start_suite`, `end_suite`, `start_test` and `end_test`? That's interesting, what's the reason?

As we know, from RobotFramework 2.8, the execution and logging are all inherited from `robot.model.visitor.SuiteVisitor`. And let's check `XmlLogger`.

* robot.output.xmllogger.XmlLogger
    - robot.result.visitor.ResultVisitor
        - robot.model.visitor.SuiteVisitor

After finding out this, I think it should be possible to implement a listener to generate log file for each sub suite. In the `start_suite`, a `XmlLogger` instance will be created, and it will be used in other interfaces like `start_test`, `start_keyword`, and in the `end_suite`, the xml file should be done and closed. So based on this, I wrote the first version.

```python
from robot.output import XmlLogger

class SuiteLogger:
    ROBOT_LISTENER_API_VERSION = 2

    def start_suite(self, name, attributes):
        self._logger = XmlLogger(name + '.output.xml')

    def end_suite(self, name, attributes):
        self._logger.end_suite(_DictObj(attributes))
        self._logger.close()

    def start_test(self, name, attributes):
        self._logger.start_test(_DictObj(attributes))

    def end_test(self, name, attributes):
        self._logger.end_test(_DictObj(attributes))

    def start_keyword(self, name, attributes):
        self._logger.start_keyword(_DictObj(attributes))

    def end_keyword(self, name, attributes):
        self._logger.end_keyword(_DictObj(attributes))

    def log_message(self, message):
        self._logger.log_message(_DictObj(message))

    def message(self, message):
        self._logger.message(_DictObj(message))

    def set_log_level(self, level):
        self._logger.set_log_level(level)


class _DictObj(object):
    def __init__(self, attributes):
        self._attrs = attributes

    def __getattr__(self, attr):
        if attr in self._attrs:
            return self._attrs.[attr]
        raise AttributeError
```
Test it, unfortunately it failed. I got many error messages like below:
```
[ ERROR ] Calling listener method 'start_test' of listener 'listener.SuiteLogger' failed: AttributeError
```

Checking the code, I found `XmlLogger` always need a object has `name` field, but the listener interface does not contain it.
Add `attributes['name'] = name` before calling `self._logger`, most errors dispeared, but still AttirbuteError with `message` and `timeout`.
This time, update `_DictObj`, instead of raising AttributeError, return None.
Test again, most errors were removed, and `output.xml` was generated. But there was no any log message in it, only the structure. As executing `message` and `log_message`, the `self._logger` is always `None`.

Thanks to RobotFramework, it provides a global context for us. It is `robot.running.EXECUTION_CONTEXTS`. So I wrote the first workable version.

```python
from robot.output import XmlLogger
from robot.running import EXECUTION_CONTEXTS

class SuiteLogger:
    ROBOT_LISTENER_API_VERSION = 2

    def start_suite(self, name, attributes):
        attributes['name'] = name
        self._get_logger().start_suite(_DictObj(attributes))

    def end_suite(self, name, attributes):
        attributes['name'] = name
        self._get_logger().end_suite(_DictObj(attributes))
        self._get_logger().close()

    def start_test(self, name, attributes):
        attributes['name'] = name
        self._get_logger().start_test(_DictObj(attributes))

    def end_test(self, name, attributes):
        attributes['name'] = name
        self._get_logger().end_test(_DictObj(attributes))

    def start_keyword(self, name, attributes):
        attributes['name'] = name
        attributes['type'] = 'kw'
        self._get_logger().start_keyword(_DictObj(attributes))

    def end_keyword(self, name, attributes):
        attributes['name'] = name
        attributes['type'] = 'kw'
        self._get_logger().end_keyword(_DictObj(attributes))

    def log_message(self, message):
        if self._get_logger():
            self._get_logger().log_message(_DictObj(message))

    def message(self, message):
        if self._get_logger():
            self._get_logger().message(_DictObj(message))

    def set_log_level(self, level):
        self._get_logger().set_log_level(level)

    def _get_logger(self):
        current = EXECUTION_CONTEXTS.current
        if not current:
            return None
        if hasattr(current, 'suite_logger'):
            return current.suite_logger
        current.suite_logger = XmlLogger(current.suite.name + '.output.xml')
        return current.suite_logger

class _DictObj(object):
    def __init__(self, attributes):
        self._attrs = attributes

    def __getattr__(self, attr):
        return self._attrs.get(attr, None)
```

I have `output.xml` of sub-suite generated correctly, and then I can add code the generate `log.html` and `report.html`. The generated path is another issue, it should base on the settings of complete one.

The latest code is at `https://github.com/feiyuw/idiomatic-robotframework/blob/master/examples/listener.py`.
