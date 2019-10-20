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

func (f *Feedly) Subscribers() (string, error) {
	url := f.createAPIURL()
	resp, err := f.Client.Get(url)
	if err != nil {
		return "", errors.Wrap(err, "error while fetching feedly API")
	}

	defer resp.Body.Close()
	c, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "error while readling feedly API response")
	}

	subscribers, err := extractSubscribers(c)
	if err != nil {
		return "", err
	}

	return subscribers, nil
}

func (f *Feedly) createAPIURL() string {
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

func extractSubscribers(data []byte) (string, error) {
	fr := &FeedlyResponse{}
	err := json.Unmarshal(data, fr)
	if err != nil {
		return "", errors.Wrap(err, "error while unmarshal feedly API response")
	}

	if _, ok := fr.Results[0]["subscribers"]; !ok {
		return "", errors.New("wrong format for Feedly API")
	}

	s := fr.Results[0]["subscribers"].(float64)
	subs := strconv.FormatFloat(s, 'f', 0, 64)

	return subs, nil
}
