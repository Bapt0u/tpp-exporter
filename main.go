package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"sync"
	"venafiTools/tpp"
)

var config tpp.Config

func main() {
	fmt.Println("##################################")
	fmt.Println("#                                #")
	fmt.Println("# Venafi-exporter is starting... #")
	fmt.Println("#                                #")
	fmt.Println("##################################")
	fmt.Print("")

	log.Print("Loading conf file at ./conf/venafi.yml")

	config.GetConf("./conf/venafi.yml")

	var wg sync.WaitGroup
	var v = new(tpp.VenafiTpp)
	v.Connect(&config)

	// Create a non global registry
	reg := prometheus.NewRegistry()

	// Create new metrics and register them using the custom registry
	m := tpp.NewMetrics(reg)

	wg.Add(1)
	go m.Updater(v, config.VenafiTpp.ScrapeTime)

	// Expose metrics and custom registry via an http server
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Print("system; msg=Awake and ready to serve on localhost:2112")
	err := http.ListenAndServe(":2112", nil)
	tpp.CheckErr(err)
}
