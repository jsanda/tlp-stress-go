package run

import (
	"github.com/jsanda/tlp-stress-go/pkg/generators"
	"log"
	"time"
)

type profileRunnerCfg struct{
	PartitionKeyGenerator *generators.PartitionKeyGenerator
}

type profileRunner struct{
	Population int64
	PartitionKeyGenerator *generators.PartitionKeyGenerator
}

func createRunners() *profileRunner {
	return &profileRunner{}
}

func (p *profileRunner) Populate(rows uint64, done chan struct{}) {
	defer close(done)
	log.Printf("Populating Cassandra with %d rows\n", rows)

	// TODO maxId needs to be configurable
	maxId := uint64(100000)
	ch := p.PartitionKeyGenerator.GenerateKey(rows, maxId)

	for key := range ch {
		// Get the next mutation for the key
		// Execute the query for the mutation
		key.Prefix = "foo"
		p.Population++
		time.Sleep(50 * time.Nanosecond)
	}
}
