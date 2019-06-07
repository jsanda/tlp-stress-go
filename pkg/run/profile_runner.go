package run

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/jsanda/tlp-stress-go/pkg/generators"
	"github.com/jsanda/tlp-stress-go/pkg/profiles"
	"log"
	"sync"
	"sync/atomic"
)

type profileRunner struct{
	Population            int64
	PartitionKeyGenerator *generators.PartitionKeyGenerator
	Profile               profiles.StressProfile
	StressRunner          profiles.StressRunner
	Concurrency           int64
	Session               *gocql.Session
}

func createRunners(cfg *StressCfg) *profileRunner {
	// TODO tlp-stress has support for multiple profile runners, each being run in a separate thread

	runner := &profileRunner{
		Population: cfg.Population,
		Profile: cfg.Plugin.Instance,
		StressRunner: cfg.Plugin.Instance.GetRunner(cfg.Registry, cfg.Session),
		Concurrency: cfg.Concurrency,
		Session: cfg.Session,
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

func (p *profileRunner) Populate(rows int64, count *int64, done chan struct{}) {
	defer close(done)

	// TODO maxId needs to be configurable
	maxId := uint64(100000)
	ch := p.PartitionKeyGenerator.GenerateKey(rows, maxId)
	mutations := make(chan *profiles.Mutation)
	var wg sync.WaitGroup

	wg.Add(int(p.Concurrency))
	p.applyMutations(&wg, mutations, count)

	for key := range ch {
		op := p.StressRunner.GetNextMutation(key)
		mutations <- op
	}
	close(mutations)
	wg.Wait()
}

func (p *profileRunner) applyMutations(wg *sync.WaitGroup, mutations <-chan *profiles.Mutation, count *int64) {
	for i := int64(0); i < p.Concurrency; i++ {
		go func() {
			for mutation := range mutations {
				if err := mutation.Query.Exec(); err == nil {
					// TODO record execution time metric
				} else {
					log.Printf("An error occurred prepopulating data: %s\n", err)
					// TODO record execution time metric
					// TODO record error metric
				}
				atomic.AddInt64(count, 1)
			}
			wg.Done()
		}()
	}
}
