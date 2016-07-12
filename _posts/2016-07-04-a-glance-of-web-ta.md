---
layout: post
title:  "A glance of web test automation"
date:   2016-07-04 15:00:00 +0800
categories: TA
---
I'm not sure how many people in the world are out of Internet. But Internet did great success in past 20 years, what ever the content and the technology. That is impossible for me to summarize the technologies used in the past 20 years, I can only list the technology I have used or heard, and the related test automation solutions to this technologies.

It's a home work I got from one community of my company, great chanllenge but also a good chance to do a retrospective of myself.

## A glance of web development technology

This is the first chapter, although "impossible", I still try to draw a draft of the web development technology used in the past 20 years.

### Big Events

There are many big events of web development during past ~20 years, including web server, browser, protocol, programming language and framework.

* *Feb 1995* `Apache` Group created
* *May 1995* `Javascript` created
* *May 1996* `HTTP/1.0` published
* *1996* `iframe` tag introduced
* *1998* `JSP` introduced
* *June 1998* `PHP 3.0` released
* *June 1999* `HTTP/1.1` published
* *2002* `Nginx` created
* *July 2004* `Ruby on Rails` released
* *November 2004* `Firefox` 1.0 released
* *Feb 2005* `Ajax` publicly stated
* *April 2006* first draft specification for the `XMLHttpRequest` object released
* *2008* `HTML5` published
* *December 2008* `Chrome` released
* *2009* `NodeJS` created
* *2009* `AngularJS` released
* *May 2013* `ReactJS` open-sourced
* *March 2015* `React Native` open-sourced
* *May 2015* `HTTP/2` published

If we summarize the core technology stack of these years, we may find, in past 20 years, there are great progress in web development.

### Technology stack

* Static web with few inputs
    * Before 2000, most websites are **READ ONLY**, with few inputs like `search` or `shopping cart`.
    * See yahoo.com homepage of [1996](https://web.archive.org/web/19961220154510/http://www.yahoo.com/).
    * Many news web sites had background jobs to generate html files during night.
* "Dynamic Website"
    * Most PHP/ASP/JSP driven websites were treated as dynamic websites, usually they connected to a database like MySQL or Oracle, hosted on Apache/IIS.
* Ajax and rich client
    * `Gmail` and `Google Maps` helped everyone to know the power of AJAX.
    * From that time, the website contained more Javascript code, and frontend engineer became more and more important.
* Single Page Website
    * `AngularJS`, `ReactJ` and many other frameworks made `SPA` a fasion technology.
    * Many DOM elements are not from HTML, but created by Javascript dynamically.
    * The Javascript code become much more complex than before.
    * Packager tool and JS compiler widely used.
* Mobile web
    * We need to think about how to make the UI smooth on bad network environment.
    * We also need to optimize the traffic size.
    * Touch/drag/zoom events replace the mouse events.
    * Performance of UI rendering become a big issue.
    * Huge amount of device models should be covered.

## Only for the web data (API, "static" page)

The ROI(return of investment) decrease from small tests on a single function to the large test on a whole application.

![ROI](/assets/web/roi.png)

In most time, implementing the testing automation for huge test scenario cost too high, we will choose the medium one. In web, the easiest and most important test scenario should be the API, the second is the page with less communications.

### API

Now, `JSON` may be the most popular data structure used in API service, below are two examples written in `NodeJS` and `Python`.

#### `restify` module in NodeJS

```javascript
const restify = require('restify');

let client = restify.createJsonClient('http://127.0.0.1');
client.get('/api/actions', (err, req, res, obj) => {console.log(obj)});
```

#### Python related modules

##### urllib

```python
import urllib
import json

opener = urllib.urlopen('http://127.0.0.1/api/actions')
json.loads(opener.read())
```

##### requests

```python
import requests

s = requests.session()
res = s.get('http://127.0.0.1/api/actions')
res.json()
```

### "static" page

If the data you need to fetch or verify is from a "static" page, which means there is no need to run `javascript` engine to generate the DOM, use `requests` module in Python may be a good solution.


```python
import requests

res = requests.get('http://www.google.com')
res.status_code # 200
res.headers
res.cookies.get_dict()
res.elapsed # datetime.timedelta(0, 0, 208763)
res.content
```

## Functional testing related

Before you start to use a `functional testing automation` framework, make sure you really need it and the scope of using it.

There is an interesting phenomenon that the Project Manager or the Test Manager always ask the testing should be end to end, although the engineer may tell him/her the cost is too high, and writing more medium test cases are better.
And this decision always made in `big` companies.

Why?

This is my suppose, SW development is full of risks, if a critical bug found after software release, most companies will do the root cause analysis. If the root cause is lacking of user similar testing coverage, the related person will take the resposibilities. But if the policy is what ever the risk and investment is, just do the end to end testing. That will lead to the high cost of testing, but less responsibilities, and in big companies, the resource is always not a problem.

### Use "webdriver" to test a web site

W3C plan to introduce webdriver as a standard, see https://www.w3.org/TR/webdriver/.
For most websites, [webdirver](http://www.seleniumhq.org/) is a good functional testing frmeworks.

As it can:

* browser independent
* platform independent
* use browser native API, and will become a W3C standard
* programming support of Python, Java, Javascript ...
* execute javascript code

#### Example

In this example, we will open a web page, then visit its all links, to make sure each link can be visited, and no critical logs occurred.

```python
from selenium import webdriver

class MyWebLib(object):
    def __init__(self, browser='Chrome'):
        _browser_name = browser.capitalize()
        if not hasattr(webdriver, _browser_name):
            raise RuntimeError('webdriver does not support browser "%s"!' % _browser_name)
        self._browser = getattr(webdriver, _browser_name)()

    def __getattr__(self, attr):
        if not attr.startswith('_') and hasattr(self._browser, attr):
            return getattr(self._browser, attr)

    def no_critical_errors(self):
        error_logs = filter(lambda log: log['level'] in ('ERROR', 'SEVERE'), self.console_logs)
        if error_logs:
            raise RuntimeError(error_logs)

    @property
    def console_logs(self):
        return self.get_log('browser')

    def get_all_links(self):
        links = self._browser.find_elements_by_tag_name('a')
        return list(set([e.get_attribute('href').strip('#/') for e in links]))


if __name__ == '__main__':
    client = MyWebLib('chrome')
    client.get('http://127.0.0.1')
    client.no_critical_errors()
    for link in client.get_all_links():
        if link not in ('javascript:void(0)', 'http://127.0.0.1', 'http://127.0.0.1/index.html'):
            client.get(link)
            client.no_critical_errors()
```

As you see, `webdriver` is very easy to use, you can write the code in the REPL like `ipython`, after everything is OK, put them together.

Unlike `requests` or `urllib`, `webdriver` will let your test act like a real user. That means not only the response of the HTTP request you fire, but also the resources like css, images and javascript files, and the `onReady` javascript code will be executed as also.

**But** the abstract level is still [DOM](//www.w3.org/DOM/), you can find a button, click it, and fill some fields, submit the form, etc. Everythink you do is based on the `DOM`.


### Test an "AngularJS" or "ReactJS" application

After [AngularJS](//angularjs.org/) and [ReactJS](//facebook.github.io/react/) released, the TA solution based on `DOM` become more and more unconvinient.
As the modularization of APP, we need to communicate with one `Component`, not one `DOM`.

So the legacy solution may have some issues:

* hard to know the page has finished rendering
* hard to know the page has finished re-rendering some components after the state update
* hard to know the AJAX calls are fired and finished

Using wait is a dirty solution, which we prefer something better.

For `AngularJS`, [protractor](//www.protractortest.org/) is a good solution.
For `ReactJS`, [ReactTestUtils](//facebook.github.io/react/docs/test-utils.html) is good for unit testing, but for end to end testing, seems no good choice.

### Headless solution

Running testing on a real Display sometimes is impossible, like in a continuous integration agent, as the agent is always a Linux virtual machine without X.

> [xvfb](//www.x.org/archive/X11R7.6/doc/man/man1/Xvfb.1.xhtml) is an X server that can run on machines with no display hardware and no physical input devices.
> It emulates a dumb framebuffer using virtual memory.

Based on `xvfb` and its python wrapper [pyvirtualdisplay](//github.com/ponty/pyvirtualdisplay), we can build a headless solution on a real browser like Chrome and Firefox.

```python
from pyvirtualdisplay import Display

vx = Display(visible=0, size=(1024, 768))
vx.start()

# do testing here

vx.stop()
```

### Integrate with ATDD/BDD frameworks like "RobotFramework"

Changing your python common code to a RobotFramework library is almost **0** effort. Or you can use exist library [Selenium2Library](//github.com/robotframework/Selenium2Library).

I recommanded you writing your own one, as the interface of webdriver is very simple, and it's quite easy for you to write your own.

## Performance testing related

`Performance Testing` is not just a `Testing` work, before you do the testing, you should know:

* Used technologies
* Deployment (docker is a good option)
* The critical part (UI rendering or the concurrent user amount)
* Acceptance criteria

Different application may have different solutions, in most web application, response time and status in large amount of users are the critical criterias.

There are many free or commercial related tools, but I prefer using [locust](//locust.io). As it:

* is open-sourced
* support 10k users based on gevent
* write scenarios in Python which is a real and good programming language
* has good command line support
* master/slave mode support

#### Example

```python
from locust import HttpLocust, TaskSet, task
from locust.exception import StopLocust
import random
import json
import sys


class SubTask(TaskSet):
    def get(self, url):
        self.client.get(url)

    def post(self, url, data):
        self.client.post(url, data=json.dumps(data),
            headers={'content-type': 'application/json'})

class UserBehavior(TaskSet):
    @task(50)
    class _Index(SubTask):
        @task(10)
        def index(self):
            self.get('/')

        @task(20)
        def search_testline(self):
            filter_exprs = (' ', 'f', 'fSM', 'FSMF', 'T22', 'FSMF T22 ')
            self.post('/api/testconf', {'filter': filter_exprs[random.randrange(len(filter_exprs))]})

    @task(10)
    class _ExpertPool(SubTask):
        def on_start(self):
            self.post('/api/auth', {'uid': 'atest', 'passwd': 'nopassword'})

        @task(5)
        def expertpool(self):
            self.get('/expertpool.html')

        @task(20)
        def search_expertpool(self):
            filter_exprs = (' ', '', 'sta', 'start', '@', 'atest', '@atest')
            self.post('/api/pool', {'filter': filter_exprs[random.randrange(len(filter_exprs))]})

        @task(3)
        def view_expertpool(self):
            self.post('/api/pool', {'id': 'startup@atest'})


class WebsiteUser(HttpLocust):
    task_set = UserBehavior
    min_wait = 5000
    max_wait = 9000

    def run(self):
        self.task_set(self).run()
```

Above is a sample test suites with contain two pages with some executions.

To run the test, use `locust` command:

```sh
locust -f perf_test.py --no-web -c 1000 -r 100 -n 10000 --host=http://127.0.0.1
```


