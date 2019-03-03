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

var mappingGscHeader = map[string]string{
	"page":  "Page",
	"query": "Query",
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
) ([][]string, error) {
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

	return formatSearchTable(resp.Rows, dimension), nil
}

func filtersFromString(filters string, dimension string) []*sc.ApiDimensionFilter {
	fg := []*sc.ApiDimensionFilter{}
	if filters != "" {
		f := strings.Split(filters, ",")
		for _, v := range f {
			v = strings.TrimSpace(v)
			dim := dimension
			if strings.Contains(v, " ") {
				t := strings.Split(v, " ")
				dim = t[0]
				v = t[1]
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

func formatSearchTable(rows []*sc.ApiDataRow, dimension string) [][]string {
	result := [][]string{{mappingGscHeader[dimension], "Clicks", "Impressions", "CTR", "Position"}}

	for _, v := range rows {
		result = append(result, []string{
			strings.Join(v.Keys, ","),
			fmt.Sprintf("%g", v.Clicks),
			fmt.Sprintf("%g", v.Impressions),
			fmt.Sprintf("%.5f", v.Ctr),
			fmt.Sprintf("%.2f", v.Position),
		})
	}

	return result
}
