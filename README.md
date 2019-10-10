![travis CLI](https://travis-ci.org/Phantas0s/devdash.svg?branch=master&style=for-the-badge) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/ec1e19b08f3b40d19f3acaf93e3e186b)](https://www.codacy.com/app/Phantas0s/devdash?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=Phantas0s/devdash&amp;utm_campaign=Badge_Grade)  [![Go Report Card](https://goreportcard.com/badge/github.com/Phantas0s/devdash)](https://goreportcard.com/report/github.com/Phantas0s/devdash) [![Hits-of-Code](https://hitsofcode.com/github/phantas0s/devdash)](https://hitsofcode.com/view/github/phantas0s/devdash) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) 
![logo of devdash with a gopher](./doc/img/logo.jpg) 
[![Tweet](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/intent/tweet?text=DevDash%20-%20Highly%20Configurable%20Terminal%20Dashboard%20For%20Developers:&url=https%3A%2F%2Fgithub.com%2Fphantas0s%2Fdevdash&hashtags=developers,dashboard,terminal,CLI,golang) [![ko-fi](https://www.ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/T6T4W5K0)

DevDash is a highly configurable terminal dashboard for developers, who want to choose and display the most up-to-date metrics they need, at one place.

# Why using DevDash?

* *Choose* the metrics you specifically need.
* All the important data in your cosy terminal.
* Pull data from Github, Google Analytics or Google Search Console. More services to come!
* Unlimited amount of different dashboards with different configurations.
* Widgets' data refreshed automatically.
* A huge amount of flexibility compared to other terminal dashboards:
  * Choose the widgets you want.
  * Place your widgets where you want.
  * Choose the data you want to display, the colors you want to use, and a lot of other things for each widget.

# Menu

* [Installation](#installation)
* [Documentation](#documentation)
* [Acknowledgement](#acknowledgement)
* [Contribute](#contribute)
* [Licence](#licence)
* [Making of](#Making-of)
* [Showcase](#showcase)

# Installation

You can simply grab the [latest released binary file](https://github.com/Phantas0s/devdash/releases/latest) and download the version you need, depending on your OS.

## Linux script

Here's a simple way to download DevDash and move it in `/usr/local/bin`, in order to be able to use DevDash everywhere easily.

```shell
curl -LO https://raw.githubusercontent.com/Phantas0s/devdash/master/install/linux.sh && \
sh ./linux.sh && \
rm linux.sh
```

## Manual installation

You need Golang installed to compile DevDash.

You simply need to run `go get -u github.com/Phantas0s/devdash/cmd/devdash` in your terminal.

# Documentation

[The documentation is here.](https://thedevdash.com)

In there you will find:

* Installation / getting started
* Simple examples and real use cases
* Complete reference of dashboard configurations

# Acknowledgement

Thanks to [MariaLetta](https://github.com/MariaLetta/free-gophers-pack) for the awesome and beautiful Gopher pack she made! I used it for my logo on top.

DevDash was inspired from other open source projects:

* [wtf](https://github.com/wtfutil/wtf)
* [tdash](https://github.com/jessfraz/tdash)

# Contribute

First of all, thanks a lot if you want to contribute to DevDash!

I think the ["talk, then code"](https://dave.cheney.net/tag/contributing) practice is pretty good to avoid misunderstandings and hours of work for nothing.

Therefore:

"Every new feature or bug fix should be discussed with the maintainer(s) of the project before work commences. It’s fine to experiment privately, but do not send a change without discussing it first."

# Making Of

For anybody interested how I managed to develop DevDash on side of a full time job, and how I organized my time and kept my motivation, [I wrote an article about that on my blog](https://thevaluable.dev/programming-side-project-example-devdash/).

# Licence

[Apache Licence 2.0](https://choosealicense.com/licenses/apache-2.0/)

# Showcase

![google analytics example DevDash configuration](./example/img/mix-1.png)
-------
![google analytics example DevDash configuration](./example/img/thevaluabledev-2.png)
-------
![google analytics example DevDash configuration](./example/img/thevaluabledev-3.png)
-------
![github example DevDash configuration](./example/img/devdash-1.png)

