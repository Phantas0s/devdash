// Google Search Console
package plateform

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
	sc "google.golang.org/api/webmasters/v3"
)

type SearchConsole struct {
	config  *jwt.Config
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

	web.service, err = sc.NewService(context.Background(), option.WithHTTPClient(web.config.Client(context.Background())))
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
