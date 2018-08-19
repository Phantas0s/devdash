package main

import (
	"flag"
	"fmt"
	"io/ioutil"
)

func main() {
	file := flag.String("config", ".termetrics.yml", "The config file")
	flag.Parse()

	data, _ := ioutil.ReadFile(*file)
	cfg := mapConfig(data)

	client, err := newClient(cfg.GoogleAnalytics.Keyfile, false)
	if err != nil {
		fmt.Println(err)
	}
	res, err := client.GetReport(cfg.GoogleAnalytics.ViewID)
	if err != nil {
		fmt.Println(err)
	}

	json, _ := res.MarshalJSON()
	fmt.Println(string(json))
}
