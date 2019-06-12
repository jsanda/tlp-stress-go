package run

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/jsanda/tlp-stress-go/pkg/generators"
	"github.com/jsanda/tlp-stress-go/pkg/metrics"
	"github.com/jsanda/tlp-stress-go/pkg/profiles"
	gometrics "github.com/rcrowley/go-metrics"
	"log"
	"math/rand"
	"time"

	//"math/rand"
	"sync"
	//"time"
)

type profileRunner struct{
	Population            uint64
	PartitionKeyGenerator *generators.PartitionKeyGenerator
	Profile               profiles.StressProfile
	StressRunner          profiles.StressRunner
	Concurrency           uint64
	Session               *gocql.Session
	Metrics               *metrics.Metrics
	Duration              uint64
	Iterations            uint64
	Partitions            uint64
}

func createRunners(cfg *StressCfg) *profileRunner {
	// TODO tlp-stress has support for multiple profile runners, each being run in a separate thread

	runner := &profileRunner{
		Population: cfg.Population,
		Profile: cfg.Plugin.Instance,
		StressRunner: cfg.Plugin.Instance.GetRunner(cfg.Registry, cfg.Session),
		Concurrency: cfg.Concurrency,
		Session: cfg.Session,
		Metrics: cfg.Metrics,
		Duration: cfg.Duration,
		Iterations: cfg.Iterations,
		Partitions: cfg.Partitions,
	}

	thread := 0
	prefix := fmt.Sprintf("%s.%d.", cfg.Id, thread)

	log.Printf("Creating generator %s", cfg.PartitionKeyGenerator)

	switch cfg.PartitionKeyGenerator {
	case "random":
		runner.PartitionKeyGenerator = generators.Random(prefix)
	case "sequence":
		runner.PartitionKeyGenerator = generators.Sequence(prefix)
	default:
		log.Fatalf("%s is not a valid generator\n", cfg.PartitionKeyGenerator)
	}

	return runner
}

func (p *profileRunner) Populate(rows uint64, done chan struct{}) {
	defer close(done)

	// TODO maxId needs to be configurable
	maxId := uint64(100000)
	ch := p.PartitionKeyGenerator.GenerateKey(rows, maxId)
	ops := make(chan *profiles.Operation)
	var wg sync.WaitGroup

	wg.Add(int(p.Concurrency))
	p.execPopulate(&wg, ops)

	for key := range ch {
		op := p.StressRunner.GetNextMutation(key)
		ops <- op
	}
	close(ops)
	wg.Wait()
}

func (p *profileRunner) execPopulate(wg *sync.WaitGroup, ops <-chan *profiles.Operation) {
	for i := uint64(0); i < p.Concurrency; i++ {
		go func() {
			var err error
			for mutation := range ops {
				p.Metrics.Populate.Time(func() {
					err = mutation.Query.Exec()
				})
				if err != nil {
					log.Printf("An error occurred prepopulating data: %s\n", err)
					p.Metrics.Errors.Mark(1)
				}
				p.Metrics.PopulateCount.Inc(1)
			}
			wg.Done()
		}()
	}
}

func (p *profileRunner) Run() {
	if p.Duration == 0 {
		log.Printf("Running the profile for %d iterations...\n", p.Iterations)
	} else {
		log.Printf("Running the profile for %dmin\n", p.Duration)
	}

	totalValues := p.Iterations
	ch := p.PartitionKeyGenerator.GenerateKey(totalValues, p.Partitions)
	readRate := float64(0.1)
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	ops := make(chan *profiles.Operation)
	var wg sync.WaitGroup

	wg.Add(int(p.Concurrency))
	p.execOperations(&wg, ops)

	for key := range ch {
		if int32(readRate * 100) > rand.Int31n(100) {
			ops <- p.StressRunner.GetNextSelect(key)
		} else {
			ops <- p.StressRunner.GetNextMutation(key)
		}

	}
}

func (p *profileRunner) execOperations(wg *sync.WaitGroup, ops <-chan *profiles.Operation) {
	for i := uint64(0); i < p.Concurrency; i++ {
		go func() {
			var err error
			var timer gometrics.Timer
			var counter gometrics.Counter

			for op := range ops {
				if op.OperationType == profiles.Mutation {
					timer = p.Metrics.Mutations
					counter = p.Metrics.MutationsCount
				} else {
					timer = p.Metrics.Selects
					counter = p.Metrics.SelectsCount
				}
				timer.Time(func() {
					err = op.Query.Exec()
				})
				if err != nil {
					p.Metrics.Errors.Mark(1)
				}
				counter.Inc(1)
			}
			wg.Done()
		}()
	}
}