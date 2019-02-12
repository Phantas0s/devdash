package plateform

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	gav3 "google.golang.org/api/analytics/v3"
	ga "google.golang.org/api/analyticsreporting/v4"
)

const gaPrefix = "ga:"
const earliestDate = "2005-01-01"

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

// GetReport queries the Analytics Reporting API V4 using the
// Analytics Reporting API V4 service object.
func (c *Client) Users(viewID string, startDate string, endDate string) ([]string, []int, error) {
	req := &ga.GetReportsRequest{
		ReportRequests: []*ga.ReportRequest{
			{
				ViewId: viewID,
				DateRanges: []*ga.DateRange{
					{StartDate: startDate, EndDate: endDate},
				},
				Metrics: []*ga.Metric{
					{Expression: "ga:users"},
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
func (c *Client) Pages(viewID string, startDate string, endDate string, global bool) (dim []string, u []int, err error) {

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
					{Expression: "ga:pageViews"},
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
		return dim[0]
	}

	return format(resp.Reports, formater)
}

// NewVsReturning queries the Analytics Reporting API V4 using the
// Analytics Reporting API V4 service object.
func (c *Client) ReturningVsNew(viewID string, startDate string, endDate string) (dim []string, u []int, err error) {
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

	return formatReturningNew(resp.Reports, formater)
}

func (c *Client) TrafficSource(viewID string, startDate string, endDate string) (dim []string, u []int, err error) {
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
					{Name: "ga:source"},
				},
				OrderBys: []*ga.OrderBy{
					{
						FieldName: "ga:sessions",
						SortOrder: "DESCENDING",
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
			"can't get traffic source data from google analytics with start data %s / end_date %s",
			startDate,
			endDate,
		)
	}

	formater := func(dim []string) string {
		return dim[0]
	}

	return format(resp.Reports, formater)
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

func formatReturningNew(reps []*ga.Report, dimFormater func(dim []string) string) (dim []string, u []int, err error) {
	var new []int
	var ret []int
	for _, v := range reps {
		for l := 0; l < len(v.Data.Rows); l++ {
			if v.Data.Rows[l].Dimensions[0] == "New Visitor" {
				dim = append(dim, dimFormater(v.Data.Rows[l].Dimensions))
			}

			for m := 0; m < len(v.Data.Rows[l].Metrics); m++ {
				value := v.Data.Rows[l].Metrics[m].Values[0]

				var vu int64
				if vu, err = strconv.ParseInt(value, 0, 0); err != nil {
					return nil, nil, err
				}
				if v.Data.Rows[l].Dimensions[0] == "New Visitor" {
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

func debug(rep *ga.Report) {
	test, _ := json.Marshal(rep)
	fmt.Println(string(test))
}
