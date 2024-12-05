# StatsApi

Wrap prometheus api and make it easy to use

## Example 
```go
package main

import (
	"context"
	"fmt"
	"net/http"
	
)

var client stats.Client

func index(w http.ResponseWriter, req *http.Request) {
	// #4. Set another count metric
	err := client.Count(req.Context(),
		"svc",                                         // contextTag
		"http_requests_total",                          // key
		1,                                              // amount
		stats.Tag{Key: "url", Value: req.URL.Path},  // tag
		stats.Tag{Key: "method", Value: req.Method}, // another tag
	)
	if err != nil {
		log.Warnf("cannot count cause: %s", err.Error())
	}

	fmt.Fprint(w, "hello\n")
}

func main() {

	// #1. Initializing the stats api client
	client = stats.NewClient(nil)

	// #2. Set a gauge goutils
	err := client.Gauge(context.TODO(),
		"svc",                                       // context Tag
		"app_info",                                   // key
		1,                                            // value
		stats.Tag{Key: "version", Value: "1.0.1"}, // tag
	)
	if err != nil {
		log.Warnf("cannot set gauge cause: %s", err.Error())
	}

	http.HandleFunc("/", index)

	// #3. Registry the goutils endpoint, where prometheus will scrape goutils from
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```
