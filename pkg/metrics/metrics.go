package metrics

import (
	gometrics "github.com/rcrowley/go-metrics"
	"log"
	"os"
	"time"
)

type Metrics struct {
	Registry       gometrics.Registry
	Errors         gometrics.Meter
	Mutations      gometrics.Timer
	MutationsCount gometrics.Counter
	Selects        gometrics.Timer
	SelectsCount   gometrics.Counter
	Populate       gometrics.Timer
	PopulateCount  gometrics.Counter
}

func NewMetrics() *Metrics {
	registry := gometrics.NewRegistry()

	go gometrics.Log(registry, 5 * time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	return &Metrics{
		Registry: registry,
		Errors: gometrics.NewRegisteredMeter("errors", registry),
		Mutations: gometrics.NewRegisteredTimer("mutations", registry),
		MutationsCount: gometrics.NewRegisteredCounter("mutationsCount", registry),
		Selects: gometrics.NewRegisteredTimer("selects", registry),
		SelectsCount: gometrics.NewRegisteredCounter("selectsCount", registry),
		Populate: gometrics.NewRegisteredTimer("populate", registry),
		PopulateCount: gometrics.NewRegisteredCounter("populateCount", registry),
	}
}