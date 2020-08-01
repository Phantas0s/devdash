![travis CLI](https://travis-ci.org/Phantas0s/devdash.svg?branch=master&style=for-the-badge) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/ec1e19b08f3b40d19f3acaf93e3e186b)](https://www.codacy.com/app/Phantas0s/devdash?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=Phantas0s/devdash&amp;utm_campaign=Badge_Grade)  [![Go Report Card](https://goreportcard.com/badge/github.com/Phantas0s/devdash)](https://goreportcard.com/report/github.com/Phantas0s/devdash) [![Hits-of-Code](https://hitsofcode.com/github/phantas0s/devdash)](https://hitsofcode.com/view/github/phantas0s/devdash) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) ![Documentation](https://img.shields.io/website?url=https%3A%2F%2Fthedevdash.com&label=documentation)
![logo of devdash with a gopher](./doc/img/logo.jpg) 
[![ko-fi](https://www.ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/T6T4W5K0) [![Tweet](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/intent/tweet?text=DevDash%20-%20Highly%20Configurable%20Terminal%20Dashboard%20For%20Developers:&url=https%3A%2F%2Fgithub.com%2Fphantas0s%2Fdevdash&hashtags=developers,dashboard,terminal,CLI,golang)

DevDash is a highly configurable terminal dashboard for developers, who want to choose and display the most up-to-date metrics they need, at one place.

[![google analytics example DevDash configuration](./example/img/mix-1.png)](https://raw.githubusercontent.com/Phantas0s/devdash/master/example/img/mix-1.png)

# Why using DevDash?

* Pull the data and display it in cool diagrams (or widgets) using YAML / JSON config, from: 
    * Your own computer. From your own scripts / command lines too!
    * A remote computer via SSH
    * Github
    * Travis
    * Google Analytics 
    * Google Search Console
    * Feedly
* A huge amount of flexibility compared to other terminal dashboards:
  * Choose the widgets you want.
  * Place your widgets where you want.
  * Choose the data you want to display, the colors you want to use, and a lot of other things for each widget.
  * Don't want to personalize everything? Don't overwrite the defaults, then.
* Unlimited amount of different dashboards with different configurations.
* Data refreshed automatically via time ticks, or via a keyboard shortcut (Ctrl + r by default).

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

You need to clone this repository and build the binary: `go build devdash.go`.

# How Does It Work?

In a nutshell:

* If you run DevDash without giving a dashboard configuration, it will create and display a default dashboard (`default.yml`) located in `$XDG_CONFIG_HOME/devdash` or `$HOME/.config/devdash`.
* To get used to dashboard' configurations, there are many [examples here](https://thedevdash.com/getting-started/examples/). They can help you getting started.
* To run a dashboard created in the two filepaths mentioned above, you just need to execute `dashboard -c my-super-dashboard`, if your configuration file is called `my-super-dashboard.yml`. You can use JSON as well!
* You can as well run any dashboard from anywhere if you give an absolute or relative path.
* I'm thriving to make DevDash easier to configure, yet very flexible and customizable. The next updates will go in that direction.

# Documentation

[The complete DevDash documentation is here.](https://thedevdash.com).

You'll find:

* [Installation / getting started](https://thedevdash.com/getting-started/installation/)
* [Simple examples](https://thedevdash.com/getting-started/examples/) and [real use cases](https://thedevdash.com/getting-started/use-cases/devdash/)
* [Complete reference for configuring whatever you want](https://thedevdash.com/reference/).

# Acknowledgement

Thanks to [MariaLetta](https://github.com/MariaLetta/free-gophers-pack) for the awesome and beautiful Gopher pack she made! I used it for my logo on top.

DevDash was inspired from other open source projects:

* [wtf](https://github.com/wtfutil/wtf)
* [tdash](https://github.com/jessfraz/tdash)

# Bugs and Ideas

I would be happy to read about new ideas and to fix bugs. Opening an issue is the way to go.

# Contribute

First of all, thanks a lot if you want to contribute to DevDash!

If you want to implement a new feature, let's speak about it first and decide if it fits DevDash scope.

# Making Of

For anybody interested how I managed to develop DevDash on side of a full time job, and how I organized my time and kept my motivation, [I wrote an article about that on my blog](https://thevaluable.dev/programming-side-project-example-devdash/).

# Licence

[Apache Licence 2.0](https://choosealicense.com/licenses/apache-2.0/)

# Showcase

![google analytics example DevDash configuration](./example/img/thevaluabledev-2.png)
-------
![google analytics example DevDash configuration](./example/img/thevaluabledev-3.png)
-------
![github example DevDash configuration](./example/img/devdash-1.png)

