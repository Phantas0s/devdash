// Google Search Console
package plateform

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	sc "google.golang.org/api/webmasters/v3"
)

type SearchConsole struct {
	config  *jwt.Config
	client  *http.Client
	service *sc.Service
}

type SearchConsoleResponse struct {
	Dimension   string
	Clicks      float64
	Impressions float64
	Ctr         float64
	Position    float64
}

func NewSearchConsoleClient(keyfile string) (*SearchConsole, error) {
	data, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, fmt.Errorf("reading keyfile %q failed: %v", keyfile, err)
	}

	// webmaster tools
	web := &SearchConsole{}

	web.config, err = google.JWTConfigFromJSON(data, sc.WebmastersReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("creating JWT config from json keyfile %q failed: %v", keyfile, err)
	}

	web.client = web.config.Client(context.Background())
	web.service, err = sc.New(web.client)
	if err != nil {
		return nil, fmt.Errorf("can't get webmaster service: %v", err)
	}

	return web, nil
}

func (w *SearchConsole) Table(
	viewID string,
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

// filtersFromString in the config.
// The rules are the following:
// 0. Each filters are separated with a coma
// 1. if there is one word or multiple words, the current dimension will be filtered with Contains
// 2. If there is one word or multiple words with "-" as prefix, the current dimension will be filtered with notContains
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
