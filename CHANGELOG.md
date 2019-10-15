# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [MASTER BRANCH - NOT RELEASED]

### Added

* Add Tracis CI service
    * Add Travis builds widget
* Add Feedly service
    * Add Feedly subscribers widget

## [0.2.0] - 2019-10-05

### Added

* Add Github widgets
  * Display count stars overtime
  * Display count commits overtime
  * Display issues
  * Display repositories in table with information
  * Display last week traffic on Github page

* Add `color` options to have same color for border, title and everything color related for one widget
* Add themes to simplify the configuration - possibility to use same options defined once, for multiple widgets

* Add possibility to hot reload any dashboard via a keystroke - no need to restart DevDash when changing a dashboard configuration

* Create the [official DevDash website](https://thedevdash.com)

### Updated 

* Replace `title_options` by `name_options` for project's config (breaking change)

## [0.1.1] - 2019-07-21

### Added

* use goreleaser for relases

## [0.1.0] - 2019-05-28

### Added

* Write README documentation
* Add Github widgets
  * github.box_stars
  * github.box_watchers
  * github.box_open_issues
  * github.table_branches
  * github.table_issues
  * github.table_repositories
* Add Github API
* Google Search Console widgets
  * gsc.table_pages
  * gsc.table_queries
  * gsc.table
* Add Google Search Console API
* Create ToTime library
* Google Analytics widgets:
  * ga.box_real_time
  * ga.box_total
  * ga.bar_sessions
  * ga.bar_bounces
  * ga.bar_users
  * ga.bar_returning
  * ga.bar_pages
  * ga.bar
  * ga.bar_new_returning
  * ga.table_pages
  * ga.table_traffic_sources
  * ga.table
* Google Analytics service
* Google Analytics API
* Monitoring service
* Dashboard refreshing system 
* Projects / Services / Widget system
* Display and Grid system
* YAML configuration system

[0.1.1]: https://github.com/Phantas0s/devdash/releases/tag/v0.1.1
[0.1.0]: https://github.com/Phantas0s/devdash/releases/tag/v0.1.0
