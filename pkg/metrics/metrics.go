package metrics

import (
	gometrics "github.com/rcrowley/go-metrics"
	"log"
	"os"
	"time"
)

type Metrics struct {
	Registry  gometrics.Registry
	Errors    gometrics.Meter
	Mutations gometrics.Timer
	Selects   gometrics.Timer
	Populate  gometrics.Timer
}

func NewMetrics() *Metrics {
	registry := gometrics.NewRegistry()

	go gometrics.Log(registry, 5 * time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	return &Metrics{
		Registry: registry,
		Errors: gometrics.NewRegisteredMeter("errors", registry),
		Mutations: gometrics.NewRegisteredTimer("mutations", registry),
		Selects: gometrics.NewRegisteredTimer("selects", registry),
		Populate: gometrics.NewRegisteredTimer("populateMutations", registry),
	}
}