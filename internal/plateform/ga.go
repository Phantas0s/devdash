package plateform

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	gav3 "google.golang.org/api/analytics/v3"
	ga "google.golang.org/api/analyticsreporting/v4"
)

const gaPrefix = "ga:"

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
	// TODO: remove v3 once v4 supports the realtime reporting API.
	client.servicev3, err = gav3.New(client.client)
	if err != nil {
		return nil, fmt.Errorf("creating the analytics reporting service v3 object failed: %v", err)
	}
	client.realtimeService = gav3.NewDataRealtimeService(client.servicev3)

	return client, nil
}

// GetReport queries the Analytics Reporting API V4 using the
// Analytics Reporting API V4 service object.
func (c *Client) GetReport(viewID string, startDate string, endDate string) (*ga.GetReportsResponse, error) {
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
						SortOrder: "DESCENDING",
					},
				},
				IncludeEmptyRows: true,
			},
		},
	}

	return c.service.Reports.BatchGet(req).Do()
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
