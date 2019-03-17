package plateform

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	gav3 "google.golang.org/api/analytics/v3"
	ga "google.golang.org/api/analyticsreporting/v4"
)

const gaPrefix = "ga:"

// earliest date google analytics accept for date ranges
const (
	earliestDate     = "2005-01-01"
	newVisitor       = "New Visitor"
	returningVisitor = "Returning Visitor"
)

var mappingMetrics = map[string]string{
	"sessions":                 "ga:sessions",
	"page_views":               "ga:pageViews",
	"bounces":                  "ga:bounces",
	"entrances":                "ga:entrances",
	"unique_page_views":        "ga:uniquePageviews",
	"users":                    "ga:users",
	"session_duration":         "ga:sessionDuration",
	"average_session_duration": "ga:avgSessionDuration",
	"bounce_rate":              "ga:bounceRate",
}

var mappingTimePeriod = map[string][]string{
	"day":   []string{"ga:month", "ga:day"},
	"month": []string{"ga:year", "ga:month"},
	"year":  []string{"ga:year"},
}

var mappingDimensions = map[string]string{
	"page_path":      "ga:pagePath",
	"traffic_source": "ga:source",
	"user_type":      "ga:userType",
}

var mappingHeader = map[string]string{
	"sessions":          "Sessions",
	"page_views":        "Page Views",
	"entrances":         "Entrances",
	"unique_page_views": "Unique Page Views",
}

var mappingOrder = map[string]string{
	"asc":  "ASCENDING",
	"desc": "DESCENDING",
}

type Analytics struct {
	config          *jwt.Config
	client          *http.Client
	servicev3       *gav3.Service
	realtimeService *gav3.DataRealtimeService
	service         *ga.Service
}

// New takes a keyfile for authentication and
// returns a new Google Analytics Reporting Client struct.
func NewAnalyticsClient(keyfile string) (*Analytics, error) {
	data, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, fmt.Errorf("reading keyfile %q failed: %v", keyfile, err)
	}

	an := &Analytics{}

	an.config, err = google.JWTConfigFromJSON(data, ga.AnalyticsReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("creating JWT config from json keyfile %q failed: %v", keyfile, err)
	}

	an.client = an.config.Client(context.Background())

	// analytics reporting v4 service
	an.service, err = ga.New(an.client)
	if err != nil {
		return nil, fmt.Errorf("creating the analytics reporting service v4 object failed: %v", err)
	}

	// analytics reporting v3 service object.
	an.servicev3, err = gav3.New(an.client)
	if err != nil {
		return nil, fmt.Errorf("creating the analytics reporting service v3 object failed: %v", err)
	}
	an.realtimeService = gav3.NewDataRealtimeService(an.servicev3)

	return an, nil

}

// SimpleMetric of one quantitative value.
func (c *Analytics) SimpleMetric(
	viewID,
	metric,
	startDate,
	endDate string,
	global bool,
) (string, error) {
	req := &ga.GetReportsRequest{
		ReportRequests: []*ga.ReportRequest{
			{
				ViewId: viewID,
				DateRanges: []*ga.DateRange{
					{StartDate: startDate, EndDate: endDate},
				},
				Metrics: []*ga.Metric{
					{Expression: mapMetric(metric)},
				},
				IncludeEmptyRows: true,
			},
		},
	}

	resp, err := c.service.Reports.BatchGet(req).Do()
	if err != nil {
		return "", errors.Wrapf(
			err,
			"can't get total metric data from google analytics with start data %s / end_date %s",
			startDate,
			endDate,
		)
	}

	if len(resp.Reports[0].Data.Rows) != 0 {
		return resp.Reports[0].Data.Rows[0].Metrics[0].Values[0], nil
	}

	return "0", nil
}

// BarMetric provides a bar display with one qualitative dimension and one quantitative value.
func (c *Analytics) BarMetric(
	viewID string,
	startDate string,
	endDate string,
	metric string,
	dimensions []string,
	timePeriod string,
	filters []string,
) ([]string, []int, error) {
	// Add the time dimension to the first two index of the slice (0, 1)
	tm := mapTimePeriod(timePeriod)
	dim := []*ga.Dimension{}
	for _, v := range tm {
		dim = append(dim, &ga.Dimension{Name: v})
	}

	formater := formatBar
	for _, v := range dimensions {
		if v == "user_returning" {
			formater = formatBarReturning
		}
		dim = append(dim, &ga.Dimension{Name: mapDimension(v)})
	}

	fi := []*ga.DimensionFilter{}
	if len(dimensions) == 1 {
		fi = append(fi, &ga.DimensionFilter{
			CaseSensitive: false,
			DimensionName: mapDimension(dimensions[0]),
			Expressions:   filters,
			Not:           false,
			Operator:      "PARTIAL",
		})
	}

	req := &ga.GetReportsRequest{
		ReportRequests: []*ga.ReportRequest{
			{
				ViewId: viewID,
				DateRanges: []*ga.DateRange{
					{StartDate: startDate, EndDate: endDate},
				},
				Metrics: []*ga.Metric{
					{Expression: mapMetric(metric)},
				},
				Dimensions: dim,
				OrderBys: []*ga.OrderBy{
					{
						FieldName: string(tm[0]),
						SortOrder: "ASCENDING",
					},
				},
				DimensionFilterClauses: []*ga.DimensionFilterClause{
					{
						Filters:  fi,
						Operator: "AND",
					},
				},
				IncludeEmptyRows: true,
			},
		},
	}

	resp, err := c.service.Reports.BatchGet(req).Do()
	// create fixture for tests
	// j, _ := json.Marshal(resp)
	// fmt.Println(string(j))

	if err != nil {
		return nil, nil, errors.Wrapf(
			err,
			"can't get users data from google analytics with start data %s / end_date %s",
			startDate,
			endDate,
		)
	}

	f := func(dim []string) string {
		if len(dim) >= 2 {
			return dim[0] + "-" + dim[1]
		}

		return dim[0]
	}

	return formater(resp.Reports, f)
}

// RealTimeUsers on the website (using Google api V3).
func (c *Analytics) RealTimeUsers(viewID string) (string, error) {
	metric := "rt:activeUsers"

	resp, err := c.realtimeService.Get(gaPrefix+viewID, metric).Do()
	if err != nil {
		return "", err
	}

	return resp.TotalsForAllResults[metric], nil
}

// Table shaped analytics.
// Headers on the first row are qualitative dimensions, value can be qualitative or quantitative.
func (c *Analytics) Table(
	viewID string,
	startDate string,
	endDate string,
	global bool,
	metrics []string,
	dimension string,
	orders []string,
	firstHeader string,
) (headers []string, dim []string, u [][]string, err error) {

	dateRange := []*ga.DateRange{
		{StartDate: startDate, EndDate: endDate},
	}
	if global {
		dateRange = []*ga.DateRange{
			{StartDate: earliestDate, EndDate: endDate},
		}
	}

	req := &ga.GetReportsRequest{
		ReportRequests: []*ga.ReportRequest{
			{
				ViewId:     viewID,
				DateRanges: dateRange,
				Metrics:    mapMetrics(metrics),
				Dimensions: []*ga.Dimension{
					{Name: mapDimension(dimension)},
				},
				OrderBys:         mapOrderBy(orders),
				IncludeEmptyRows: true,
			},
		},
	}

	resp, err := c.service.Reports.BatchGet(req).Do()
	if err != nil {
		return nil, nil, nil, errors.Wrapf(
			err,
			"can't get pages data from google analytics with start data %s / end_date %s",
			startDate,
			endDate,
		)
	}

	formater := func(dim []string) string {
		return dim[0]
	}

	headers = mapHeaders(firstHeader, metrics)
	dim, u = formatTable(resp.Reports, formater)
	return
}

// NewVsReturning users.
func (c *Analytics) StackedBar(
	viewID string,
	startDate string,
	endDate string,
	metric string,
	timePeriod string,
	dimensions []string,
) (dim []string, new []int, ret []int, err error) {
	// Add the time dimensions to the slice (index 1,2)
	d := []*ga.Dimension{}
	for _, v := range dimensions {
		d = append(d, &ga.Dimension{Name: mapDimension(v)})
	}
	tm := mapTimePeriod(timePeriod)
	for _, v := range tm {
		d = append(d, &ga.Dimension{Name: v})
	}

	req := &ga.GetReportsRequest{
		ReportRequests: []*ga.ReportRequest{
			{
				ViewId: viewID,
				DateRanges: []*ga.DateRange{
					{StartDate: startDate, EndDate: endDate},
				},
				Metrics: []*ga.Metric{
					{Expression: mapMetric(metric)},
				},
				Dimensions: d,
				OrderBys: []*ga.OrderBy{
					{
						FieldName: string(tm[0]),
						SortOrder: "ASCENDING",
					},
				},
				IncludeEmptyRows: true,
			},
		},
	}

	resp, err := c.service.Reports.BatchGet(req).Do()

	if err != nil {
		return nil, nil, nil, errors.Wrapf(
			err,
			"can't get pages data from google analytics with start data %s / end_date %s",
			startDate,
			endDate,
		)
	}

	formater := func(dim []string) string {
		return dim[1] + "-" + dim[2]
	}

	return formatNewReturning(resp.Reports, formater)
}

// formatBar vizualiser, to return one slice of dimension and one slice of values.
// The two slices are equals in size.
// If the same dimension is returned multiple time by the Google API, the data
// is aggregated not to have duplicated dimensions.
func formatBar(reps []*ga.Report, dimFormater func(dim []string) string) (dim []string, u []int, err error) {
	dimVal := map[string]int{}
	for _, v := range reps {
		for l := 0; l < len(v.Data.Rows); l++ {
			d := dimFormater(v.Data.Rows[l].Dimensions)

			// Add the dimension only if it was not added already.
			if _, ok := dimVal[d]; !ok {
				dim = append(dim, d)
			}

			for m := 0; m < len(v.Data.Rows[l].Metrics); m++ {
				value := v.Data.Rows[l].Metrics[m].Values[0]

				var vu int64
				if strings.Contains(value, ".") {
					f, _ := strconv.ParseFloat(value, 0)
					vu = int64(f)
				} else if vu, err = strconv.ParseInt(value, 0, 0); err != nil {
					return nil, nil, err
				}

				if _, ok := dimVal[d]; ok {
					dimVal[d] += int(vu)
					continue
				}

				dimVal[d] = int(vu)
				u = append(u, int(vu))
			}
		}
	}

	// Aggregate value with same dimension.
	for k, v := range dim {
		u[k] = dimVal[v]
	}

	return dim, u, nil
}

// formatBarReturning format the special case of new / returning users for a bar vizualisation.
func formatBarReturning(
	reps []*ga.Report,
	dimFormater func(dim []string) string,
) (dim []string, u []int, err error) {
	for _, v := range reps {
		for l := 0; l < len(v.Data.Rows); l++ {
			userType := v.Data.Rows[l].Dimensions[2]
			if userType == returningVisitor {
				dim = append(dim, dimFormater(v.Data.Rows[l].Dimensions))
			}

			for m := 0; m < len(v.Data.Rows[l].Metrics); m++ {
				value := v.Data.Rows[l].Metrics[m].Values[0]

				var vu int64
				if vu, err = strconv.ParseInt(value, 0, 0); err != nil {
					return nil, nil, err
				}

				if userType == returningVisitor {
					u = append(u, int(vu))
				}
			}
		}
	}

	return dim, u, nil
}

func formatTable(
	reps []*ga.Report,
	dimFormater func(dim []string) string,
) (dim []string, u [][]string) {
	for _, v := range reps {
		for l := 0; l < len(v.Data.Rows); l++ {
			dim = append(dim, dimFormater(v.Data.Rows[l].Dimensions))

			for m := 0; m < len(v.Data.Rows[l].Metrics); m++ {
				var g []string
				for p := 0; p < len(v.Data.Rows[l].Metrics[m].Values); p++ {
					g = append(g, v.Data.Rows[l].Metrics[m].Values[p])
				}
				u = append(u, g)
			}
		}
	}

	return dim, u
}

func formatNewReturning(
	reps []*ga.Report,
	dimFormater func(dim []string) string,
) (dim []string, new []int, ret []int, err error) {
	for _, v := range reps {
		for l := 0; l < len(v.Data.Rows); l++ {
			if v.Data.Rows[l].Dimensions[0] == newVisitor {
				dim = append(dim, dimFormater(v.Data.Rows[l].Dimensions))
			}

			for m := 0; m < len(v.Data.Rows[l].Metrics); m++ {
				value := v.Data.Rows[l].Metrics[m].Values[0]

				var vu int64
				if vu, err = strconv.ParseInt(value, 0, 0); err != nil {
					return nil, nil, nil, err
				}
				if v.Data.Rows[l].Dimensions[0] == newVisitor {
					new = append(new, int(vu))
				} else {
					ret = append(ret, int(vu))
				}
			}
		}
	}

	return dim, new, ret, nil
}

// The map functions map the properties of the application to the Google Analytics API params.

func mapMetrics(m []string) []*ga.Metric {
	gam := make([]*ga.Metric, len(m))

	for k, v := range m {
		gam[k] = &ga.Metric{Expression: strings.TrimSpace(mapMetric(v))}
	}

	return gam
}

// mapMetric will first try to search an alias for a google analytics metric,
// send the option to GoogleAnalytics otherwise
func mapMetric(metric string) string {
	m, ok := mappingMetrics[metric]
	if !ok {
		return strings.TrimSpace(metric)
	}

	return strings.TrimSpace(m)
}

func mapTimePeriod(m string) []string {
	timePeriods, ok := mappingTimePeriod[m]
	if !ok {
		timePeriods = strings.Split(m, ",")
	}

	return timePeriods
}

func mapDimension(dim string) string {
	d, ok := mappingDimensions[dim]
	if !ok {
		return strings.TrimSpace(dim)
	}

	return strings.TrimSpace(d)
}

func mapOrderBy(o []string) []*ga.OrderBy {
	gam := make([]*ga.OrderBy, len(o))

	for k, v := range o {
		s := strings.Split(v, " ")

		// default
		field := "sessions"
		order := "desc"

		if len(s) == 1 {
			field = s[0]
		}

		if len(s) > 1 {
			field = s[0]
			order = strings.ToLower(s[1])
		}

		gam[k] = &ga.OrderBy{
			FieldName: strings.TrimSpace(mapMetric(field)),
			SortOrder: strings.TrimSpace(mappingOrder[order]),
		}
	}

	return gam
}

func mapHeaders(el string, metrics []string) []string {
	h := make([]string, len(metrics)+1)
	h[0] = el

	for k, v := range metrics {
		head, ok := mappingHeader[v]
		if !ok {
			if strings.Contains(v, "ga:") {
				v = strings.Split(v, "ga:")[1]
			}

			h[k+1] = strings.TrimSpace(v)
			continue
		}
		h[k+1] = strings.TrimSpace(head)
	}

	return h
}
