package run

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/jsanda/tlp-stress-go/pkg/generators"
	"github.com/jsanda/tlp-stress-go/pkg/metrics"
	"github.com/jsanda/tlp-stress-go/pkg/profiles"
	//gometrics "github.com/rcrowley/go-metrics"
	"log"
	//"math/rand"
	"sync"
	"sync/atomic"
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

func (p *profileRunner) Populate(rows uint64, count *int64, done chan struct{}) {
	defer close(done)

	// TODO maxId needs to be configurable
	maxId := uint64(100000)
	ch := p.PartitionKeyGenerator.GenerateKey(rows, maxId)
	ops := make(chan *profiles.Operation)
	var wg sync.WaitGroup

	wg.Add(int(p.Concurrency))
	p.applyMutations(&wg, ops, count)

	for key := range ch {
		op := p.StressRunner.GetNextMutation(key)
		ops <- op
	}
	close(ops)
	wg.Wait()
}

func (p *profileRunner) applyMutations(wg *sync.WaitGroup, ops <-chan *profiles.Operation, count *int64) {
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
				atomic.AddInt64(count, 1)
			}
			wg.Done()
		}()
	}
}

func (p *profileRunner) Run() {
//	if p.Duration == 0 {
//		log.Printf("Running the profile for %d iterations...\n", p.Iterations)
//	} else {
//		log.Printf("Running the profile for %dmin\n", p.Duration)
//	}
//
//	totalValues := p.Iterations
//	ch := p.PartitionKeyGenerator.GenerateKey(totalValues, p.Partitions)
//	readRate := float64(0.1)
//	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
//	operations := make(chan *profiles.Operation)
//	var wg sync.WaitGroup
//
//	wg.Add(int(p.Concurrency))
//	p.execOperations(&wg, )
//
//	for key := range ch {
//		if int32(readRate * 100) > rand.Int31n(100) {
//			read := p.StressRunner.GetNextSelect(key)
//			op =
//		} else {
//			op = p.StressRunner.GetNextMutation(key)
//		}
//	}
}
//
//func (p *profileRunner) execOperations(wg *sync.WaitGroup, ops <- chan *profiles.Operation, timer *gometrics.Timer) {
//
//}

