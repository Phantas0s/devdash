package platform

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	gav3 "google.golang.org/api/analytics/v3"
	ga "google.golang.org/api/analyticsreporting/v4"
	"google.golang.org/api/option"
)

const gaPrefix = "ga:"

// To create fixtures in order to test the data formatting.
// j, _ := json.Marshal(resp)
// fmt.Println(string(j))

const (
	// Decide display as x-axis header for bar metrics
	XHeaderTime uint16 = iota
	XHeaderOtherDim

	earliestDate     = "2005-01-01" // This is the earliest date Google Analytics accepts.
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
	"country":        "ga:country",
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

// Analytics connect to Google Analytics API.
type Analytics struct {
	config          *jwt.Config
	servicev3       *gav3.Service
	realtimeService *gav3.DataRealtimeService
	service         *ga.Service
}

// AnalyticValues which can be possibly send to the Google Analytics API.
// This is a pure value object without behavior.
type AnalyticValues struct {
	ViewID     string
	StartDate  string
	EndDate    string
	TimePeriod string
	Global     bool
	Metrics    []string
	Dimensions []string
	Filters    []string
	Orders     []string
	RowLimit   int64
	XHeaders   uint16
}

// NewAnalyticsClient to connect to Google Analytics APIs.
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

	// analytics reporting v4 service
	an.service, err = ga.NewService(context.Background(), option.WithHTTPClient(an.config.Client(context.Background())))
	if err != nil {
		return nil, fmt.Errorf("creating the analytics reporting service v4 object failed: %v", err)
	}

	// analytics reporting v3 service object.
	an.servicev3, err = gav3.NewService(context.Background(), option.WithHTTPClient(an.config.Client(context.Background())))
	if err != nil {
		return nil, fmt.Errorf("creating the analytics reporting service v3 object failed: %v", err)
	}
	an.realtimeService = gav3.NewDataRealtimeService(an.servicev3)

	return an, nil

}

// SimpleMetric get a value depending on Google Analytics metrics.
func (c *Analytics) SimpleMetric(val AnalyticValues) (string, error) {
	req := &ga.GetReportsRequest{
		ReportRequests: []*ga.ReportRequest{
			{
				ViewId: val.ViewID,
				DateRanges: []*ga.DateRange{
					{StartDate: val.StartDate, EndDate: val.EndDate},
				},
				Metrics:          mapMetrics(val.Metrics),
				IncludeEmptyRows: true,
			},
		},
	}

	resp, err := c.service.Reports.BatchGet(req).Do()
	if err != nil {
		return "", errors.Wrapf(
			err,
			"can't get total metric data from google analytics with start data %s / end_date %s",
			val.StartDate,
			val.EndDate,
		)
	}

	if len(resp.Reports[0].Data.Rows) != 0 {
		return resp.Reports[0].Data.Rows[0].Metrics[0].Values[0], nil
	}

	return "0", nil
}

// BarMetric provides a qualitive dimension linked to a quantitative value, for example a date (dimension) with an int.
func (c *Analytics) BarMetric(val AnalyticValues) ([]string, []int, error) {
	// Add the time dimension to the first two indexes of the slice ga.Dimensions(index 0 and 1)
	tm := mapTimePeriod(val.TimePeriod)
	dim := []*ga.Dimension{}
	for _, v := range tm {
		dim = append(dim, &ga.Dimension{Name: v})
	}

	formater := formatBar
	for _, v := range val.Dimensions {
		// TODO - looks weird. We should maybe pass the formater in this function? Extract "user_returning" case in new function?
		if v == "user_returning" {
			formater = formatBarReturning
		}
		dim = append(dim, &ga.Dimension{Name: mapDimension(v)})
	}

	fi := []*ga.DimensionFilter{}
	// TODO - Why this condition? is filter works only with one dimension only?
	if len(val.Dimensions) == 1 && len(val.Filters) != 0 {
		fi = append(fi, &ga.DimensionFilter{
			CaseSensitive: false,
			DimensionName: mapDimension(val.Dimensions[0]),
			Expressions:   val.Filters,
			Not:           false,
			Operator:      "PARTIAL",
		})
	}

	req := &ga.GetReportsRequest{
		ReportRequests: []*ga.ReportRequest{
			{
				ViewId: val.ViewID,
				DateRanges: []*ga.DateRange{
					{StartDate: val.StartDate, EndDate: val.EndDate},
				},
				Metrics:    mapMetrics(val.Metrics),
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

	if err != nil {
		return nil, nil, errors.Wrapf(
			err,
			"can't get users data from google analytics with start data %s / end_date %s",
			val.StartDate,
			val.EndDate,
		)
	}

	// Decide of the header X-axis
	f := func(dimValues []string) string {
		// TODO need a mapping here instead of using an index (?) - using headers of ga response
		if val.XHeaders == XHeaderOtherDim {
			return dimValues[2]
		}

		if val.XHeaders == XHeaderTime {
			return dimValues[0] + "-" + dimValues[1]
		}

		// By default, display the first dimension for headers X-axis
		return dimValues[0]
	}

	return formater(resp.Reports, f)
}

// RealTimeUsers return the number of visitor currently on the website.
func (c *Analytics) RealTimeUsers(viewID string) (string, error) {
	metric := "rt:activeUsers"

	resp, err := c.realtimeService.Get(gaPrefix+viewID, metric).Do()
	if err != nil {
		return "", err
	}

	return resp.TotalsForAllResults[metric], nil
}

// Table display dimensions and values.
// The headers on the first row are qualitative dimensions, the values can be qualitative or quantitative.
func (c *Analytics) Table(
	an AnalyticValues,
	firstHeader string,
) (headers []string, dim []string, u [][]string, err error) {

	dateRange := []*ga.DateRange{
		{StartDate: an.StartDate, EndDate: an.EndDate},
	}

	if an.Global {
		dateRange = []*ga.DateRange{
			{StartDate: earliestDate, EndDate: an.EndDate},
		}
	}

	req := &ga.GetReportsRequest{
		ReportRequests: []*ga.ReportRequest{
			{
				ViewId:           an.ViewID,
				DateRanges:       dateRange,
				Metrics:          mapMetrics(an.Metrics),
				Dimensions:       mapDimensions(an.Dimensions),
				OrderBys:         mapOrderBy(an.Orders),
				IncludeEmptyRows: true,
				PageSize:         an.RowLimit,
			},
		},
	}

	if len(an.Filters) > 0 {

		// TODO now only one filter is possible for one dimension.
		// If there are more than one dimension, the same set of filter is applied.
		// Possibility to make it multiple filters for multipme dimension?
		filters := []*ga.DimensionFilter{}
		for _, v := range an.Dimensions {
			filters = append(filters, &ga.DimensionFilter{
				CaseSensitive: false,
				DimensionName: v,
				Expressions:   an.Filters,
				Not:           false,
				Operator:      "PARTIAL",
			})
		}

		req.ReportRequests[0].DimensionFilterClauses = []*ga.DimensionFilterClause{
			{
				Filters:  filters,
				Operator: "AND",
			},
		}
	}

	resp, err := c.service.Reports.BatchGet(req).Do()
	if err != nil {
		return nil, nil, nil, errors.Wrapf(
			err,
			"can't get table data from google analytics with start data %s / end_date %s",
			an.StartDate,
			an.EndDate,
		)
	}

	formater := func(dim []string) string {
		return dim[0]
	}

	headers = mapHeaders(firstHeader, an.Metrics)
	dim, u = formatTable(resp.Reports, formater)
	return
}

// StackedBar returns one dimension set linked with multiple values.
func (c *Analytics) StackedBar(an AnalyticValues) (dim []string, new []int, ret []int, err error) {
	d := mapDimensions(an.Dimensions)
	tm := mapTimePeriod(an.TimePeriod)

	// Add the time dimensions to the slice (index 1,2)
	for _, v := range tm {
		d = append(d, &ga.Dimension{Name: v})
	}

	req := &ga.GetReportsRequest{
		ReportRequests: []*ga.ReportRequest{
			{
				ViewId: an.ViewID,
				DateRanges: []*ga.DateRange{
					{StartDate: an.StartDate, EndDate: an.EndDate},
				},
				Metrics:    mapMetrics(an.Metrics),
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
			"can't get stacked bar data from google analytics with start data %s / end_date %s",
			an.StartDate,
			an.EndDate,
		)
	}

	formater := func(dim []string) string {
		return dim[1] + "-" + dim[2]
	}

	return formatNewReturning(resp.Reports, formater)
}

// formatBar to return one slice of dimension which elements are all linked with the elements of another slice with the values.
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

func mapDimensions(dimensions []string) []*ga.Dimension {
	d := []*ga.Dimension{}
	for _, v := range dimensions {
		d = append(d, &ga.Dimension{Name: mapDimension(v)})
	}

	return d
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
