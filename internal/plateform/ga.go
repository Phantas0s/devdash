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
	"sessions":          "ga:sessions",
	"page_views":        "ga:pageViews",
	"entrances":         "ga:entrances",
	"unique_page_views": "ga:uniquePageviews",
	"users":             "ga:users",
}

var mappingDimensions = map[string]string{
	"page_path":      "ga:pagePath",
	"traffic_source": "ga:source",
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

type Client struct {
	config          *jwt.Config
	client          *http.Client
	servicev3       *gav3.Service
	service         *ga.Service
	realtimeService *gav3.DataRealtimeService
}

// New takes a keyfile for authentication and
// returns a new Google Analytics Reporting Client struct.
func NewGaClient(keyfile string) (*Client, error) {
	data, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, fmt.Errorf("reading keyfile %q failed: %v", keyfile, err)
	}

	client := &Client{}

	client.config, err = google.JWTConfigFromJSON(data, ga.AnalyticsReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("creating JWT config from json keyfile %q failed: %v", keyfile, err)
	}

	client.client = client.config.Client(context.Background())

	// analytics reporting v4 service
	client.service, err = ga.New(client.client)
	if err != nil {
		return nil, fmt.Errorf("creating the analytics reporting service v4 object failed: %v", err)
	}

	// analytics reporting v3 service object.
	client.servicev3, err = gav3.New(client.client)
	if err != nil {
		return nil, fmt.Errorf("creating the analytics reporting service v3 object failed: %v", err)
	}
	client.realtimeService = gav3.NewDataRealtimeService(client.servicev3)

	return client, nil
}

func (c *Client) SimpleMetric(
	viewID string,
	metric, startDate string,
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

	return resp.Reports[0].Data.Rows[0].Metrics[0].Values[0], nil
}

func (c *Client) BarMetric(viewID string, startDate string, endDate string, metric string) ([]string, []int, error) {
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
				Dimensions: []*ga.Dimension{
					{Name: "ga:month"},
					{Name: "ga:day"},
				},
				OrderBys: []*ga.OrderBy{
					{
						FieldName: "ga:month",
						SortOrder: "ASCENDING",
					},
				},
				IncludeEmptyRows: true,
			},
		},
	}

	resp, err := c.service.Reports.BatchGet(req).Do()
	if err != nil {
		return nil, nil, errors.Wrapf(
			err,
			"can't get users data from google analytics with start data %s / end_date %s",
			startDate,
			endDate,
		)
	}

	formater := func(dim []string) string {
		return dim[0] + "-" + dim[1]
	}

	return format(resp.Reports, formater)
}

// RealTimeUsers queries the Analytics Realtime Reporting API V3 using the
// Analytics Reporting API V3 service object.
// It returns the Analytics Realtime Reporting API V3 response
// for how many active users are currently on the site.
func (c *Client) RealTimeUsers(viewID string) (string, error) {
	metric := "rt:activeUsers"

	resp, err := c.realtimeService.Get(gaPrefix+viewID, metric).Do()
	if err != nil {
		return "", err
	}

	return resp.TotalsForAllResults[metric], nil
}

func (c *Client) avgTimeOnPage(viewID string, startDate string, endDate string, global bool) (*ga.GetReportsResponse, error) {
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
				Metrics: []*ga.Metric{
					{Expression: "ga:avgTimeOnPage"},
				},
				Dimensions: []*ga.Dimension{
					{Name: "ga:pagePath"},
				},
				OrderBys: []*ga.OrderBy{
					{
						FieldName: "ga:pageViews",
						SortOrder: "DESCENDING",
					},
				},
				IncludeEmptyRows: true,
			},
		},
	}

	return c.service.Reports.BatchGet(req).Do()
}

// Pages queries the Analytics Reporting API V4 using the
// Analytics Reporting API V4 service object.
func (c *Client) Table(
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
	dim, u, err = formatTable(resp.Reports, formater)
	return
}

// NewVsReturning queries the Analytics Reporting API V4 using the
// Analytics Reporting API V4 service object.
func (c *Client) NewVsReturning(viewID string, startDate string, endDate string) (dim []string, u []int, err error) {
	req := &ga.GetReportsRequest{
		ReportRequests: []*ga.ReportRequest{
			{
				ViewId: viewID,
				DateRanges: []*ga.DateRange{
					{StartDate: startDate, EndDate: endDate},
				},
				Metrics: []*ga.Metric{
					{Expression: "ga:sessions"},
				},
				Dimensions: []*ga.Dimension{
					{Name: "ga:userType"},
					{Name: "ga:month"},
					{Name: "ga:day"},
				},
				OrderBys: []*ga.OrderBy{
					{
						FieldName: "ga:month",
						SortOrder: "ASCENDING",
					},
				},
				IncludeEmptyRows: true,
			},
		},
	}

	resp, err := c.service.Reports.BatchGet(req).Do()
	if err != nil {
		return nil, nil, errors.Wrapf(
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

func format(reps []*ga.Report, dimFormater func(dim []string) string) (dim []string, u []int, err error) {
	for _, v := range reps {
		for l := 0; l < len(v.Data.Rows); l++ {
			dim = append(dim, dimFormater(v.Data.Rows[l].Dimensions))

			for m := 0; m < len(v.Data.Rows[l].Metrics); m++ {
				value := v.Data.Rows[l].Metrics[m].Values[0]

				var vu int64
				if vu, err = strconv.ParseInt(value, 0, 0); err != nil {
					return nil, nil, err
				}
				u = append(u, int(vu))
			}
		}
	}

	return dim, u, nil
}

func formatTable(
	reps []*ga.Report,
	dimFormater func(dim []string) string,
) (dim []string, u [][]string, err error) {
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

	return dim, u, nil
}

func formatNewReturning(
	reps []*ga.Report,
	dimFormater func(dim []string) string,
) (dim []string, u []int, err error) {
	var new []int
	var ret []int
	for _, v := range reps {
		for l := 0; l < len(v.Data.Rows); l++ {
			if v.Data.Rows[l].Dimensions[0] == newVisitor {
				dim = append(dim, dimFormater(v.Data.Rows[l].Dimensions))
			}

			for m := 0; m < len(v.Data.Rows[l].Metrics); m++ {
				value := v.Data.Rows[l].Metrics[m].Values[0]

				var vu int64
				if vu, err = strconv.ParseInt(value, 0, 0); err != nil {
					return nil, nil, err
				}
				if v.Data.Rows[l].Dimensions[0] == newVisitor {
					new = append(new, int(vu))
				} else {
					ret = append(ret, int(vu))
				}
			}
		}
	}

	u = append(u, new...)
	u = append(u, ret...)

	return dim, u, nil
}

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

			// TODO separate camelcase to a cleaner result for header
			h[k+1] = strings.TrimSpace(v)
			continue
		}
		h[k+1] = strings.TrimSpace(head)
	}

	return h
}
