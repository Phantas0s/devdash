package internal

// TODO to test all of that!

type blogConfig struct {
	Keyfile string
	Address string
	ViewID  string
}

type githubConfig struct {
	Token string
	Owner string
	Repo  string
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
		bc.ViewID = "# view_id: The value of the viewID is required"
	}

	return bc
}

func CreateGitHubProjectConfig(token string, owner string, repo string) githubConfig {
	gc := githubConfig{}
	if token != "" {
		gc.Token = "token: " + token
	} else {
		gc.Token = "# token: env var DEVDASH_GITHUB_TOKEN"
	}

	if owner != "" {
		gc.Owner = "owner: " + owner
	} else {
		gc.Owner = "# owner: your github username"
	}

	if repo != "" {
		gc.Repo = "repository: " + repo
	} else {
		gc.Repo = "# repository: name of your GitHub repository"
	}

	return gc
}

func Blog() string {
	return `---
general:
  refresh: 600
  keys:
    quit: "C-c"
	reload: "C-r"
	edit: "C-e"

projects:
  - name: Your Blog
    title_options:
      border_color: default
      text_color: default
      size: XXL
      bold: true
    services:
      google_analytics:
        {{ .Keyfile }}
        {{ .ViewID }}
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

func GitHubProject() string {
	return `---
projects:
  - name: Your GitHub Project
    services:
      github:
        {{ .Token }}
        {{ .Owner }}
        {{ .Repo }}
    themes:
      bar:
        # Everything is yellow except the title color / bar color.
        color: yellow
        title_color: red
        bar_color: green
        bar_gap: 1
      table:
        border_color: green
        row_limit: 10
      ocean:
        border_color: blue
        num_color: black
        bar_color: cyan
        title_color: magenta
        bar_gap: 1
    widgets:
      - row:
          - col:
              size: 12
              elements:
                # The widget is of type "bar": the theme bar is applied.
                - name: github.bar_stars
      - row:
          - col:
              size: 6
              elements:
                - name: github.bar_views
                  # The theme "ocean" override the theme "bar".
                  theme: ocean
                  options:
                    height: 23
                    bar_gap: 5
                    bar_width: 6
          - col:
              size: 6
              elements:
                # The theme table is applied
                - name: github.table_issues
                  options:
                    bar_gap: 1
                    bar_width: 6
      - row:
          - col:
              size: 12
              elements:
                # The theme bar is applied
                - name: github.bar_commits
                  options:
                    start_date: 35_weeks_ago
`
}
