package main

import (
	"cmp"
	"fmt"
	"log"
	"net/http"
	"os"
	"prometheus_F670L/ont"
	internalPrometheus "prometheus_F670L/prometheus"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	session, err := ont.Login(
		cmp.Or(strings.TrimRight(os.Getenv("ENDPOINT"), "/"), "http://192.168.1.1"),
		cmp.Or(os.Getenv("ONT_USERNAME"), "user"),
		cmp.Or(os.Getenv("ONT_PASSWORD"), "user"),
	)

	if err != nil {
		fmt.Println("Login failed:", err)
		return
	}

	log.Println("Login succeeded")
	log.Println("Loading ONT Collector")

	collector := internalPrometheus.NewONTCollector(session)
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	log.Println("Registering metrics")

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		EnableOpenMetrics: false,
	}))

	log.Println("Starting HTTP server on :3000")

	log.Fatal(http.ListenAndServe(":3000", nil))
}
