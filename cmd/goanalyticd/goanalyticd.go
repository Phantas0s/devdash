package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/Phantas0s/termetrics/internal/plateform"
)

func main() {
	file := flag.String("config", ".termetrics.yml", "The config file")
	flag.Parse()

	data, _ := ioutil.ReadFile(*file)

	cfg := plateform.MapConfig(data)
	fmt.Println(cfg)
}
