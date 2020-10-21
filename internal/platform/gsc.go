// Google Search Console
package platform

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
	sc "google.golang.org/api/webmasters/v3"
)

// SearchConsole connect to the Google Search Console API.
type SearchConsole struct {
	config  *jwt.Config
	service *sc.Service
}

// SearchConsoleResponse returned after requesting the API.
type SearchConsoleResponse struct {
	Dimension   string
	Clicks      float64
	Impressions float64
	Ctr         float64
	Position    float64
}

// NewSearchConsoleClient create a SearchConsole.
func NewSearchConsoleClient(keyfile string) (*SearchConsole, error) {
	data, err := ioutil.ReadFile(keyfile)
	if err != nil {
		var err2 error
		home, _ := homedir.Dir()
		data, err2 = ioutil.ReadFile(home + "/.config/devdash/" + keyfile)
		if err2 != nil {
			return nil, fmt.Errorf("reading keyfile %q failed: %v", keyfile, err)
		}
	}

	// webmaster tools
	web := &SearchConsole{}

	web.config, err = google.JWTConfigFromJSON(data, sc.WebmastersReadonlyScope)
	if err != nil {
		return nil, errors.Errorf("creating JWT config from json keyfile %q failed: %v", keyfile, err)
	}

	web.service, err = sc.NewService(
		context.Background(),
		option.WithHTTPClient(web.config.Client(context.Background())),
	)
	if err != nil {
		return nil, errors.Errorf("can't get webmaster service: %v", err)
	}

	return web, nil
}

// Table of Google Search Console with a dimension and its values.
func (w *SearchConsole) Table(
	startDate string,
	endDate string,
	limit int64,
	address string,
	dimension string,
	filters string,
) ([]SearchConsoleResponse, error) {
	req := &sc.SearchAnalyticsQueryRequest{
		StartDate:  startDate,
		EndDate:    endDate,
		Dimensions: []string{dimension},
		DimensionFilterGroups: []*sc.ApiDimensionFilterGroup{
			{
				Filters: filtersFromString(filters, dimension),
			},
		},
		RowLimit: limit,
	}

	resp, err := w.service.Searchanalytics.Query(address, req).Do()
	if err != nil {
		return nil, err
	}

	return formatSearchTable(resp.Rows), nil
}

// filtersFromString declared in the config.
func filtersFromString(filters string, dimension string) []*sc.ApiDimensionFilter {
	fg := []*sc.ApiDimensionFilter{}
	if filters != "" {
		f := strings.Split(filters, ",")
		for _, v := range f {
			v = strings.TrimSpace(v)
			dim := dimension
			if strings.Contains(v, " ") {
				t := strings.Split(v, " ")
				if strings.Contains(t[0], "*") {
					dim = strings.Trim(t[0], "*")
					v = strings.Join(t[1:], " ")
				}
			}

			operator := "contains"
			if string(v[0]) == "-" {
				operator = "notContains"
				v = v[1:]
			}

			fg = append(fg, &sc.ApiDimensionFilter{
				Dimension:  dim,
				Expression: v,
				Operator:   operator,
			})
		}
	}

	return fg
}

func formatSearchTable(rows []*sc.ApiDataRow) []SearchConsoleResponse {
	results := []SearchConsoleResponse{}
	for _, v := range rows {
		results = append(results, SearchConsoleResponse{
			Dimension:   strings.Join(v.Keys, ","),
			Clicks:      v.Clicks,
			Impressions: v.Impressions,
			Ctr:         v.Ctr,
			Position:    v.Position,
		})
	}

	return results
}
