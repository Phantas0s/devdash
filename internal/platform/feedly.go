package platform

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

type Feedly struct {
	Address string
	Client  *http.Client
}

type FeedlyResponse struct {
	Results []map[string]interface{} `json:"results"`
}

const (
	version = "v3"
	domain  = "https://feedly.com"
	method  = "search/feeds"
)

func NewFeedly(address string) *Feedly {
	return &Feedly{
		Address: address,
		Client:  &http.Client{},
	}
}

func (f Feedly) Subscribers() (string, error) {
	url := f.createAPIURL()
	fmt.Println(url)
	resp, err := f.Client.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	c, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fr := &FeedlyResponse{}
	err = json.Unmarshal(c, fr)
	if err != nil {
		return "", err
	}

	if _, ok := fr.Results[0]["subscribers"]; !ok {
		return "", errors.New("empty response from API")
	}

	subs := fr.Results[0]["subscribers"].(float64)

	return strconv.FormatFloat(subs, 'f', 0, 64), nil
}

func (f Feedly) createAPIURL() string {
	params := "n=50&fullTerm=false&organic=true&promoted=true"
	return fmt.Sprintf(
		"%s/%s/%s?q=%s&%s",
		domain,
		version,
		method,
		url.QueryEscape(f.Address),
		params,
	)
}
