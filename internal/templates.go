package internal

type blogConfig struct {
	Keyfile string
	Address string
	ViewID  string
}

func CreateBlogConfig(keyfile string, viewID string, address string) blogConfig {
	bc := blogConfig{}
	if keyfile != "" {
		bc.Keyfile = "keyfile: " + keyfile
	} else {
		bc.Keyfile = "# keyfile: env variable DEVDASH_GA_KEYFILE"
	}

	if address != "" {
		bc.Address = "address: " + address
	} else {
		bc.Address = "# address: The value of the address is required"
	}

	if viewID != "" {
		bc.ViewID = "view_id: " + viewID
	} else {
		bc.ViewID = "# view_id: The value of the address is required"
	}

	return bc
}

func GA() string {
	return `---
projects:
  - name: https://thevaluable.dev - General
    title_options:
      border_color: default
      text_color: default
      size: XXL
      bold: true
    services:
      google_analytics:
        {{ .Keyfile }}
        {{ .ViewID }}
        view_id: 89379071
      feedly:
        {{ .Address }}
      monitor:
        {{ .Address }}
      google_search_console:
        {{ .Keyfile }}
        # keyfile: env variable DEVDASH_GSC_KEYFILE
        {{ .Address }}
    themes:
      table:
        color: blue
      box:
        color: green
        text_color: default

    widgets:
      - row:
          - col:
              size: "M"
              elements:
                - name: ga.bar_sessions
                  options:
                    start_date: "10_days_ago"
                    end_date: "today"
                    color: yellow
                    num_color: black
                    bar_gap: 1
                    bar_width: 6
                - name: ga.bar_sessions
                  options:
                    start_date: "12_months_ago"
                    end_date: "this_month"
                    time_period: month
                    color: red
                    bar_width: 8
                    num_color: black
                    bar_gap: 1
          - col:
              size: "S"
              elements:
                - name: ga.bar_new_returning
                  options:
                    start_date: "5_months_ago"
                    end_date: "this_month"
                    time_period: month
                    metric: "users"
                    title_color: blue
                    bar_width: 8
                    num_color: black
                    bar_gap: 2
                    height: 20
          - col:
              size: "XS"
              elements:
                - name: ga.box_real_time
                  options:
                - name: mon.box_availability
                  options:
                    num_color: default
                - name: ga.box_total
                  options:
                    title: "Average session duration this month"
                    metric: "ga:avgSessionDuration"
                    start_date: this_month
                    end_date: this_month
                - name: ga.box_total
                  options:
                    title: "sessions/users 2 weeks ago"
                    metric: "ga:sessionsPerUser"
                    start_date: 2_weeks_ago
                    end_date: 2_weeks_ago
                - name: ga.box_total
                  options:
                    title: "sessions/users 1 week ago"
                    metric: "ga:sessionsPerUser"
                    start_date: last_week
                    end_date: last_week
                - name: ga.box_total
                  options:
                    title: "Total Users From Beginning"
                    metric: "users"
                    start_date: 30_months_ago
                    end_date: today
                - name: ga.box_total
                  options:
                    title: "Bounce Rate (%) this month"
                    metric: "ga:bounceRate"
                    start_date: this_month
                    end_date: this_month
                - name: feedly.box_subscribers
                  options:
                    title: "Feedly"
      - row:
          - col:
              size: "S"
              elements:
                - name: ga.table_traffic_sources
                  options:
                    title: " Sources | Today "
                    start_date: "today"
                    end_date: "today"
                    row_limit: 17
          - col:
              size: "S"
              elements:
                - name: ga.table
                  options:
                    title: " Referrers "
                    dimension: "ga:fullReferrer"
                    metrics: "sessions"
                    start_date: "today"
                    end_date: "today"
                    row_limit: 100
                    character_limit: 65
                    `
}
