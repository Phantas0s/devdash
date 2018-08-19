package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
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

// New takes a keyfile for auththentication and
// returns a new Google Analytics Reporting Client struct.
// Your credentials should be obtained from the Google
// Developer Console (https://console.developers.google.com).
// Navigate to your project, then see the "Credentials" page
// under "APIs & Auth".
// To create a service account client, click "Create new Client ID",
// select "Service Account", and click "Create Client ID". A JSON
// key file will then be downloaded to your computer.
func newClient(keyfile string, debug bool) (*Client, error) {
	// Read the keyfile.
	data, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, fmt.Errorf("reading keyfile %q failed: %v", keyfile, err)
	}

	// Create the initial client.
	client := &Client{}

	// Create a JWT config from the keyfile.
	client.config, err = google.JWTConfigFromJSON(data, ga.AnalyticsReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("creating JWT config from json keyfile %q failed: %v", keyfile, err)
	}

	// The following GET request will be authorized and authenticated
	// on the behalf of your service account.
	if debug {
		ctx := context.WithValue(
			context.Background(),
			oauth2.HTTPClient,
			&http.Client{},
			// &http.Client{Transport: &logTransport{http.DefaultTransport}},
		)
		client.client = client.config.Client(ctx)
	} else {
		client.client = client.config.Client(context.Background())
	}

	// Construct the analytics reporting v4 service object.
	client.service, err = ga.New(client.client)
	if err != nil {
		return nil, fmt.Errorf("creating the analytics reporting service v4 object failed: %v", err)
	}

	// Construct the analytics reporting v3 service object.
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
// It returns the Analytics Reporting API V4 response
func (c *Client) GetReport(viewID string) (*ga.GetReportsResponse, error) {
	req := &ga.GetReportsRequest{
		ReportRequests: []*ga.ReportRequest{
			{
				ViewId: viewID,
				DateRanges: []*ga.DateRange{
					{StartDate: "7daysAgo", EndDate: "today"},
				},
				Metrics: []*ga.Metric{
					{Expression: "ga:pageviews"},
					{Expression: "ga:uniquePageviews"},
					{Expression: "ga:users"},
				},
				Dimensions: []*ga.Dimension{
					{Name: "ga:pagePath"},
				},
				OrderBys: []*ga.OrderBy{
					{FieldName: "ga:pageviews", SortOrder: "DESCENDING"},
				},
			},
		},
	}

	// Call the BatchGet method and return the response.
	return c.service.Reports.BatchGet(req).Do()
}

// GetRealtimeActiveUsers queries the Analytics Realtime Reporting API V3 using the
// Analytics Reporting API V3 service object.
// It returns the Analytics Realtime Reporting API V3 response
// for how many active users are currently on the site.
func (c *Client) GetRealtimeActiveUsers(viewID string) (string, error) {
	metric := "rt:activeUsers"

	// Call the realtime get method.
	resp, err := c.realtimeService.Get(gaPrefix+viewID, metric).Do()
	if err != nil {
		return "", err
	}

	return resp.TotalsForAllResults[metric], nil
}
