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

If we summarize the core technology stack of these years, we may find, in past 20 years,

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

## Only for the web data (API, "static" HTML)

* `restify` module in NodeJS
    ```python
    const restify = require('restify');

    let client = restify.createJsonClient('http://10.69.81.34');
    client.get('/api/actions', (err, req, res, obj) => {console.log(obj)});
    ```
* Python related modules
    * urllib
        ```python
        import urllib
        import json

        opener = urllib.urlopen('http://10.69.81.34/api/actions')
        json.loads(opener.read())
        ```
    * requests
        ```python
        import requests

        s = requests.session()
        res = s.get('http://10.69.81.34/api/actions')
        res.json()
        ```

## Functional testing related

### selenium-webdriver + RobotFramework

### Headless solution

* xvfb + pyvirtualdisplay
* phantomjs

### CSS/UI assertion

## Performance testing related

* HTTP protocol
* Web application architecture
* Interface
    * jmeter
    * locust
* UI

