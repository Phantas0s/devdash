Devdash

DevDash is a terminal-based dashboard, for developers who wants to have the up-to-date data they need about their projects, at one place.

# What DevDash can bring you

* All the metrics you need at one place
* Maximum flexibility via config file
  * Choose the widgets you want
  * Place your widgets easily thanks to a grid system
  * **Many** options available to configure each widget
* Pull data from Google Analytics and Google Search Console. More services to come!
* Unlimited amount of different dashboards with different configurations
* Widgets are refreshed on a chosen time period, to have the most up-to-date data
* It's in the terminal!

# Menu

1. Installation
2. Getting Started

# Installation

## Linux

You can simply grab the [latest released binary file]() and put it wherever you want.

Here's a simple way to download and put DevDash in `/usr/local/bin`, which should be part of your path.

```shell
cd /usr/local/bin && sudo curl -LO https://github.com/Phantas0s/devdash/releases/download/v0.0.1/devdash && sudo chmod 755 devdash && cd -
```

# Getting started

For now, DevDash has three services: 

* Google Analytics (ga)
* Google Search Console (gsc)
* A monitoring service (mon)

To see DevDash in action and get familiar with the config file, you can easily configure the monitoring service.

![img](example/img/monitor.png)

Here the config to create this (very simple) dashboard:

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

* Copy past the config in a new file (`monitoring.yml`, for example)
* Run DevDash: `devdash -config monitoring.yml`

Congratulation! You just created your first dashboard. DevDash will simply send a request every `600` seconds to `https://www.web-techno.net` and display the response's status code.

# Authorization and permissions

Some services need credentials and permissions for Devdash to pull the data into your shiny terminal. Here are the detailed step by step to create these permissions:

## Google Data

### Downloading the authorization JSON file 

1. Go to [https://console.developers.google.com/apis/api/webmasters.googleapis.com/credentials](Google APIs Credentials).
2. Select `Service account key`.
3. Create a new service account.
4. Select the role `Project -> Viewer` for a read access only.
5. Add a name.
6. Click on the button `create`.
7. Save the `Service account ID` somewhere. We will need it later.
8. Download the JSON file. **Its path need to be specified in the config of your DevDash dashboard.**

### Pulling data from Google Analytics

1. Go to your [Google Analytics account](https://search.google.com/search-console/users).
2. Select the application you want to use with Devdash in the property column (just below the column name `property`).
3. Click on the `+` button and add a user.
4. Enter the email address you saved before (the `Service account ID` of step 9).
5. Click on `View settings` on the `View` column and copy the `View ID` into your **DevDash configuration file**.

### Pulling data from Google Search Console

1. Go to your [Google Search Console account](https://search.google.com/search-console/)
2. Click the property you want to access with DevDash.
3. Click on `Settings` in the menu on the right.
4. Click on `Users and permissions` and add a user with the `Service Account ID` as email address.

# Configuration examples

Click on the screenshot to see the config of these examples

## Google Analytics

[<img src="./example/img/ga-1.png" alt="monitor_widget" type="image/png" >](./example/ga-1.yml)
[<img src="./example/img/ga-2.png" alt="monitor_widget" type="image/png" >](./example/ga-2.yml)

## Google Search Console

You can find more examples in the folder `example`.

## Mix services

[<img src="./example/img/mix-1.png" alt="monitor_widget" type="image/png" >](./example/mix-1.yml)

# Structure of the config file

Since a diagram is better than a wall of text, here we go:

@startuml

general : Global configuration of your dashboard
projects : List of your project
services: Configurations of every services you want to use
widgets: List of widgets you want to display.
row: Create a row which contains columns
col: Create a column which contains widgets
size: Size of the column (T-shirt sizes or number 0-12)
elements: Your actual widgets and their configuration


general-->projects
projects-->services
projects--->widgets
widgets-->row
row-->col
col-->size
col--->elements

@enduml

# Widget displays

There are three category of widgets:

* `box` - a single value in a box
* `bar` - a bar diagram with multiple values overtime
* `table` - data in a tabble


# Configuration reference

## Monitoring

### Service configuration

```yml
    services:
      monitor:
        address: "https://www.my-website.net"
```

### Widgets available

| Name                       | Description                                                                       |
| -------------------------- | --------------------------------------------------------------------------------- |
| mon.box_availability           | Send an HTTP request to the application and display the response's status code|

##### Data Options

None.

##### Display Options

| Name             | Description  | Default value           | Examples                                |
| ---------------- | -            | -----------------       | -------------------------------         |
| title            | Title        | Depending on the widget | ` Users `                               |
| height           | Height       | `10`                    | `5`                                     |
| title_color      | Title color  | `Default color`         | `yellow`, `red` (see [colors](#colors)) |
| border_color     | Border color | `Default color`         | `yellow`, `red` (see [colors](#colors)) |
| text_color       | Text color   | `Default color`         | `yellow`, `red` (see [colors](#colors)) |

## Google Analytics

### Service configuration

```yml
    services:
      google_analytics:
        keyfile: goanalytics-abc123.json
        view_id: 456789123
```

### Widgets available

Here's the list of every widgets and their different configuration. 

| Name                     | Description                                                                       |
| -                        | --------------------------------------------------------------------------------- |
| ga.box_real_time         | Number of users on the website now                                                |
| ga.box_total             | Total of any metric on a given time period                                        |
| ga.bar_sessions          | Count of sessions overtime                                                        |
| ga.bar_bounces           | Count of bounce sessions overtime                                                 |
| ga.bar_users             | Count of users overtime                                                           |
| ga.bar_returning         | Count of returning users overtime                                                 |
| ga.bar_pages             | Count of sessions (or any other metric like users) on specific page(s) overtime   |
| ga.bar                   | Count of theoretically any metrics from Google Analytics overtime                 |
| ga.bar_new_returning     | Count of new and returning users overtime                                         |
| ga.table_pages           | Display chosen data about pages on a given time period                            |
| ga.table_traffic_sources | Display Data about traffic sources on a given time period                         |
| ga.table                 | Display theoretically any metrics from Google Analytics on a given time period    |

### Widget Options

#### Bar Widgets

##### Data Options

| Name        | Description                                                                       | Default value       | Examples                               | Not available for                                    |
| -           | --------------------------------------------------------------------------------- | ------------------- | -------------------------------------- | ---------------                                      |
| start_date  | Start date of time period                                                         | `7_days_ago`        | `2018-01-01`, `2_weeks_ago`            |                                                      |
| end_date    | End date of time period                                                           | `today`             | `2018-01-31`, `2_weeks_ago`            |                                                      |
| time_period | Time period represented by a bar (days, months, years)                            | `days`              | `days`, `months`, `years`              |                                                      |
| metric      | Google analytics metric                                                           | `sessions`          | `page_views`, `bounces`, `entrances`   | `ga.bar_pages`, `ga.bar_returning`                   |
| dimensions  | Google analytics dimensions. Multiple value possible separated with a comma (,)   |                     | `page_path`, `user_types`              | `ga.bar_pages`, `ga.bar_bounces`, `ga.bar_returning` |
| filters     | Query filter. `-` can be used in front to exclude instead of include              |                     | `value`, `-value`                      |                                                      |

##### Display Options

| Name         | Description                | Default value             | Examples                        |
| -            | -------------------------- | -----------------         | ------------------------------- |
| title        | Widget title               | `Depending on the widget` | `Users `                        |
| border_color | Border color of the widget | `Default color`           | `yellow`, `red` (see colors)    |
| height       | Widget height              | `10`                      | `5`                             |
| title_color      | Title color  | `Default color`         | `yellow`, `red` (see [colors](#colors)) |
| text_color   | Text color                 | `Default color`           | `yellow`, `red` (see colors)    |
| num_color    | Color of numerical data    | `Default color`           | `yellow`, `red` (see colors)    |
| bar_color    | Bar color                  | `Default color`           | `yellow`, `red` (see colors)    |
| bar_gap      | Space size between the bar | `0`                       | `5`, `10`                       |
| bar_width    | Bar width                  | `6`                       | `5`, `10`                       |

#### Table widgets

##### Data Options

| Name            | Description                                                               | Default value                                     | Examples                                     | Not used by                |
| -               | ------------------------------------------------                          | ---                                               | ----------------------------------------     | -----------------          |
| start_date      | Start date of time period                                                 | `7_days_ago`                                      | `2018-01-01`, `2_weeks_ago`                  |                            |
| end_date        | End date of time period                                                   | `today`                                           | `2018-01-31`, `2_weeks_ago`                  |                            |
| metrics         | Google analytics metrics. Multiple values possible separated with a comma | `sessions,page_views,entrances,unique_page_views` | `bounces,sessions`, `entrances`              |                            |
| dimension       | Google analytics dimension                                                | `page_path`                                       | `2018-01-31`, `2_weeks_ago`                  | `ga.table_traffic_sources` |
| orders          | Order of the result. Multiple value possible separated with a comma       | `sessions desc`                                   | `sessions desc,page_views asc`. `page_views` |                            |
| filters         | Query filter (prefix `-` to exclude)                                      |                                                   | `value`, `-value`                            |                            |
| row_limit       | Limit the row number                                                      | 5                                                 | 5, 100                                       |                            |
| character_limit | Limit the number of character of the dimension                            | 1000                                              | 100, 200                                     |                            |

##### Display Options

| Name         | Description           | Default value             | Examples                     | Not used by   |
| -            | -                     | -----------------         | -                            | ------------- |
| title        | Widget's title        | `Depending on the widget` | `Users `                     |               |
| title_color      | Title color  | `Default color`         | `yellow`, `red` (see [colors](#colors)) |
| border_color | Widget's border color | `Default color`           | `yellow`, `red` (see colors) |               |
| text_color   | Text color            | `Default color`           | `yellow`, `red` (see colors) |               |

#### Box widgets

##### Data Options

| Name             | Description                                                                     | Default value     | Examples                             | 
|------------------| - |-------------------|--------------------------------------|
| start_date       | Start date of time period                                                       | `7_days_ago`      | `2018-01-01`, `2_weeks_ago`          |
| end_date         | End date of time period                                                         | `today`           | `2018-01-31`, `2_weeks_ago`          |
| metric           | Google analytics metric                                                         | `sessions`        | `page_views`, `bounces`, `entrances` |

##### Display Options

| Name             | Description                | Default value             | Examples                        |
| ---------------- | - | -----------------         | ------------------------------- |
| title            | Widget title               | `Depending on the widget` | `Users `                        |
| height           | Widget's height            | `10`                      | `5`                             |
| title_color      | Widget's title color       | `Default color`           | `yellow`, `red` (see colors)    |
| border_color     | Widget's border color      | `Default color`           | `yellow`, `red` (see colors)    |
| text_color       | Text color                 | `Default color`           | `yellow`, `red` (see colors)    |

### Examples 
## Google Search Console

### Service configuration

```yml
    services:
      google_search_console:
        keyfile: goanalytics-abc123.json
```

### Widgets available

 | Name              | Description                                                                           |
 | -                 | ---------------------------------------------------------------------------------     |
 | gsc.table_pages   | Display clicks, impressions, ctr, position for pages                                  |
 | gsc.table_queries | Display clicks, impressions, ctr, position for queries                                |
 | ga.table          | Display theoretically any dimension from Google Search Console on a given time period |

### Widget Options

#### Table widgets

##### Data Options

| Name            | Description                                                                   | Default value                     | Examples                                     | Not used by       |
| _               | --                                                                            | -                                 | -                                            | ----------------- |
| start_date      | Start date of time period                                                     | `7_days_ago`                      | `2018-01-01`, `2_weeks_ago`                  |                   |
| end_date        | End date of time period                                                       | `today`                           | `2018-01-31`, `2_weeks_ago`                  |                   |
| metrics         | Google Search Console metrics (multiple values possible separated with `,`)   | `clicks,impressions,ctr,position` | `query`, `page`                              | `gsc.table_pages` |
| dimension       | Google Search Console dimension (multiple values possible separated with `,`) | `quert`,                          | `2018-01-31`, `2_weeks_ago`                  |                   |
| orders          | Order of the result. (multiple values possible separated with `,`)            | `sessions desc`                   | `sessions desc,page_views asc`. `page_views` |                   |
| filters         | Filter the default dimension (prefix `-` to exclude)                          |                                   | `value`, `-super value`                      |                   |
| row_limit       | Limit the row number                                                          | 5                                 | 5, 100                                       |                   |
| character_limit | Limit the number of character of the dimension                                | 1000                              | 100, 200                                     |                   |

##### Display Options

| Name         | Description                | Default value             | Examples                        | Not used by   |
| -            | -------------------------- | -----------------         | ------------------------------- | ------------- |
| title        | Widget's title             | `Depending on the widget` | `Users `                        |               |
| border_color | Widget's border color      | `Default color`           | `yellow`, `red` (see colors)    |               |
| text_color   | Text color                 | `Default color`           | `yellow`, `red` (see colors)    |               |

# General references

## Options values

### Colors

The list of colors you can use with the widgets's options:


| Name    |
|---------|
| default |
| black   |
| red     |
| green   |
| yellow  |
| blue    |
| magenta |
| cyan    |
| white   |

These colors depend on the configuration of your terminal

### Size

Devdash is based on a 12 columns grid.

You can indicate the width of a widget in number of column, or using the equivalent t-shirt size as described below:

| Name | Number of columns |
| --   | --                |
| xxs  | 1                 |
| xs   | 2                 |
| s    | 4                 |
| m    | 6                 |
| l    | 8                 |
| xl   | 10                |
| xxl  | 12                |

# Contribute

# Licence

[Apache Licence 2.0](https://choosealicense.com/licenses/apache-2.0/)
