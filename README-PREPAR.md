Devdash

DevDash is a terminal based dashboard for developers who wants to have all the up-to-date information they need about their projects, at one place.

# What DevDash can bring you

* All the metrics you need at one place
* Maximum flexibility via config file
  * Choose your widgets
  * Place your widgets easily thanks to a grid system
  * Many options available for each widget
* Data from Google Analytics and Google Search Console available. More services to come!
* Unlimited amount of different dashboards with different configs
* Widgets are refreshed on a chosen time period to have the most up-to-date data
* It's in the terminal, so it's cool.

# Menu

# Installation

## Linux

You can simply grab the [latest released binary file]() and put it wherever you want.

Here's a simple way to download and put DevDash in `/usr/local/bin`, which should be part of your path.

```shell
cd /usr/local/bin && sudo curl -LO https://github.com/Phantas0s/devdash/releases/download/v0.0.1/devdash && sudo chmod 755 devdash && cd -
```

# Getting started

For now, DevDash has three services: 

* A monitoring service
* Google Api
* Google Search Console.

To see DevDash in action and get familiar with the config file, you can easily configure the monitoring service. 
It has one widget: `avalability`. Here's what it looks like:

![img](example/img/monitor.png)

Here the config to create this dashboard:

```yml
---
general:
  refresh: 600
  keys:
    quit: "C-c"

projects:
  - name: Quickstart
    services:
      monitor:
        address: "https://www.web-techno.net"
    widgets:
      - row:
          - col:
              size: "M"
              elements:
                - name: mon.box_availability
                  options:
                    border_color: green
```

You can: 
* Copy past the config in a new file (`monitoring.yml` for example) (you can choose another name, of course)
* Run DevDash: `devdash -config <config_filename>.yml`

Congratulation! You just created your first dashboard. DevDash will simply send a request every `600` seconds to `https://www.web-techno.net` and display the response's status code.

Our dashboard looks a bit empty and boring for now, but it won't last long.

# Structure of the config file

Since a diagram is better than a wall of text, here we go:

![Alt text](https://g.gravizo.com/source/custom_mark13?https%3A%2F%2Fraw.githubusercontent.com%2Phantas0s%2Fdevdash%2Fmaster%2FREADME-PREPAR.md)
<details> 
<summary></summary>
custom_mark13
@startuml;

general : Global configuration of your dashboard;
projects : List of your project;
services: Configurations of every services you want to use;
widgets: List of widgets you want to display;
row: Create a row which contains columns;
col: Create a column which contains widgets;
size: Size of the column (T-shirt sizes or number 0-12);
elements: Your actual widgets and their configuration;


general-->projects;
projects-->services;
projects--->widgets;
widgets-->row;
row-->col;
col-->size;
col--->elements;

@enduml;

# Widget displays

There are three way to display your widgets:

* box 
* bar (bar diagram)
* table

The widget's type and the service you need are directly in the name of the widget itself. For example, the widget `mon.box_availability` will display a box widget and you need the `monitor` service correctly configured (with every required fields) for it to work.

# Services available

## Google Services

You can pull data from your Google Analytics and Google Search Console accounts directly to DevDash.
In order to do so, you need to authorize DevDash to pull the data from there.

I wrote tutorials how to do exactly that:

* [Google Analytics]()
* [Google Search Console]()

### Google Analytics

#### Service configuration

```yml
    services:
      google_analytics:
        keyfile: goanalytics-abc123.json
        view_id: 456789123
```

#### Widgets available

| Name                       | Description                                                                       |
| -------------------------- | --------------------------------------------------------------------------------- |
| ga.box_real_time           | Number of users on the website now                                                |
| ga.box_total               | Total of any metric on a given time period                                        |
| ga.bar_sessions            | Count of sessions overtime                                                        |
| ga.bar_bounces             | Count of bounce sessions overtime                                                 |
| ga.bar_users               | Count of users overtime                                                           |
| ga.bar_returning           | Count of returning users overtime                                                 |
| ga.bar_new_returning       | Count of new and returning users overtime                                         |
| ga.bar_pages               | Count of sessions (or any other metric like users) on specific page(s) overtime   |
| ga.bar                     | Count of theoretically any metrics from Google Analytics overtime                 |
| ga.table_pages             | Display choosed data about pages on a given time period                           |
| ga.table_traffic_sources   | Display Data about traffic sources on a given time period                         |
| ga.table                   | Display theoretically any metrics from Google Analytics on a given time period    |

#### Options available

##### Bar widgets

###### Data Options

| Name             | Description                                                                     | Required       | Default value     | Examples                             | Not used by   |
|------------------|---------------------------------------------------------------------------------|----------------|-------------------|--------------------------------------|---------------|
| start_date       | Start date of time period                                                       | no             | `7_days_ago`      | `2018-01-01`, `2_weeks_ago`          |               |
| end_date         | End date of time period                                                         | no             | `today`           | `2018-01-31`, `2_weeks_ago`          |               |
| time_period      | Time period represented by a bar (days, months, years)                          | no             | `days`            | `days`, `months`, `years`            |               |
| metric           | Google analytics metric                                                         | no             | `sessions`        | `page_views`, `bounces`, `entrances` |               |
| dimensions       | Google analytics dimensions. Multiple value possible separated with a comma (,) | no             | no dimension      | `2018-01-31`, `2_weeks_ago`          |               |
| filters          | Query filter. `-` can be used in front to exclude instead of include            | no             |                   | `value`, `-value`                    |               |

###### Display Options

| Name             | Description                | Default value             | Examples                        | Not used by   |
| ---------------- | -------------------------- | -----------------         | ------------------------------- | ------------- |
| title            | Widget title               | `Depending on the widget` | `Users `                        |               |
| border_color     | Border color of the widget | `Default color`           | `yellow`, `red` (see colors)    |               |
| height           | Widget height              | `10`                      | `5`                             |               |
| text_color       | Text color                 | `Default color`           | `yellow`, `red` (see colors)    |               |
| num_color        | Color of numerical data    | `Default color`           | `yellow`, `red` (see colors)    |               |
| bar_color        | Bar color                  | `Default color`           | `yellow`, `red` (see colors)    |               |
| bar_gap          | Space size between the bar | `0`                       | `5`                             |               |
| bar_width        | Bar width                  | `6`                       | `5`                             |               |

##### Table widgets

###### Data Options

| Name                 | Description                                                                         | Default value                                       | Examples                                 | Not used by       |
| -------------------- | ----------------------------------------------------------------------------------- | --------------------------------------------------- | ---------------------------------------- | ----------------- |
| start_date           | Start date of time period.                                                          | `7_days_ago`                                        | `2018-01-01`, `2_weeks_ago`              |                   |
| end_date             | End date of time period.                                                            | `today`                                             | `2018-01-31`, `2_weeks_ago`              |                   |
| metrics              | Google analytics metrics. Multiple values possible separated with a comma.          | "sessions,page_views,entrances,unique_page_views"   | "bounces,sessions"                       |                   |
| dimensions           | Google analytics dimensions. Multiple value possible separated with a comma.        | no dimension                                        | `2018-01-31`, `2_weeks_ago`              |                   |
| filters              | Query filter. Include by default. `-` can be used in front for exclusion.           |                                                     | `value`, `-value`                        |                   |

#### Examples 

Here are some examples. Click on the screenshot to see the config for each of them:

[<img src="./example/img/ga-1.png" alt="monitor_widget" type="image/png" >](./example/ga-1.yml)
[<img src="./example/img/ga-2.png" alt="monitor_widget" type="image/png" >](./example/ga-2.yml)

### Google Search Console

#### Service configuration

```yml
    services:
      google_search_console:
        keyfile: goanalytics-abc123.json
```

#### Widgets available

 | Name                       | Description                                                                       |
 | -------------------------- | --------------------------------------------------------------------------------- |
 | gsc.table_pages            | Display clicks, impressions, ctr, position for pages                              |
 | gsc.table_queries          | Display clicks, impressions, ctr, position for queries                              |




Authorize google search console api

https://console.developers.google.com/apis/api/webmasters.googleapis.com/overview?project=goanalytics-213713

Add user associated to the json config file to be able to see some properties
https://search.google.com/search-console/users

Example:
https://search.google.com/search-console/users?resource_id=https%3A%2F%2Fweb-techno.net%2F
