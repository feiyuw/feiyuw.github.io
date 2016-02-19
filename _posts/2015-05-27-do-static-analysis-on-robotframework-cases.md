---
layout: post
title:  "Do static analysis on RobotFramework cases"
date:   2015-05-27 21:12:00
categories: TA
---
As we know, "static analysis" is valuable and important in developing work. Many genious tools like `lint`, `klocwork`, `pyflakes` and `jshint` help the developer to find the problems earlier.

`Test Automation` is another style of developing work, but unfortunately we always ignore its quality. The result is, in most scenarios, the Test Automation asserts are hard to maintained and unefficiency.

I wrote a simple static analysis script for [RobotFramework](http://robotframework.org/) during the Sprint Festival, tried to add the lost part of `Test Automation`. In this blog, I want to discuss how I think and implement.

As the cases of RobotFramework may be saved as a HTML, a txt or csv file, we care about the logical structure instead of the plain text. So in the first step, I decided to rely on `robot.api`, and listed the checkpoints as below:

1. show bad suite/case/keyword/variable name warning
1. show too little test cases in one suite warning
1. show too many test cases in one suite warning
1. show some dangerous keywords like (Fatal Error) usage warning (or error)
1. show too many steps in one test case warning
1. show too little steps in one test case warning
1. show too many steps in one keyword warning
1. show mandentory tags missing warning
1. show too many arguments in one keyword warning
1. show set suite variable/set test variable invalid usage warning
1. show set global variable usage warning
1. show deprecated keywords used warning
1. show hard coded warning
1. show case duplication warning
1. show dry-run warning or errors
1. show no \_\_init\_\_ file warning (neither setup nor teardown defined)
1. use "Run Keyword and Ignore Error" but no return value used warning
1. performance issue warning like using sleep
1. complexity checking
1. dependency checking
1. recursive calling checking

In the same time, I found the operation I did was similar as a robot listener. And after reading the code of `RobotFramework`, I found it has provided a class called `robot.result.visitor.SuiteVisitor`, it helps to do the work to walk through the whole suite.

Based on this, the first checkpoint was ready, it is like below:

```python
import os
from robot.model import SuiteVisitor
from robot import get_version
if get_version() < '2.8.0':
    raise RuntimeError('RobotFramework 2.8+ required!')
from robot.api import TestSuiteBuilder, TestSuite

class NamingChecker(SuiteVisitor):
    def start_suite(self, _suite):
        if '.' in os.path.basename(_suite.source):
            print 'suite name should not contain "."'

if __name__ == '__main__':
    paths = sys.argv[1:]
    if '-h' in paths or '--help' in paths:
        print __doc__
        print 'usage:\n\npython rfexplain.py <path 1> <path 2> ...\n'
        sys.exit(0)
    suite = TestSuiteBuilder().build(*paths)
    suite.visit(NamingChecker())
```

In this script, it will check the name of suite, if it contains ".", it will print a message. It is simple, but it can work.

At this moment, I wrote some test robot cases to verify the script, the developing was TDD, and agile. In that two days, the first workable version was finished, including most of the checkpoints. See the gist below:

{% gist 35723b58d238a67234d8 %}

If you have interest on this script, welcome to use and modify it. If you have good ideas about it, don't forget to share with me :)
