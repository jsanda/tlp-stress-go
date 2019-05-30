package generators

import (
	"math/rand"
	"time"
)

type PartitionKeyGenerator struct {
	GenFunc func(uint64) uint64
	Prefix string
}

func Sequence(prefix string) PartitionKeyGenerator {
	current := uint64(0)
	genFunc := func(max uint64) uint64 {
		if current > max {
			current = 0
		}
		current++
		return current
	}
	return PartitionKeyGenerator{GenFunc:genFunc, Prefix:prefix}
}

func Random(prefix string) PartitionKeyGenerator {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	genFunc := func(max uint64) uint64 {
		return uint64(seededRand.Int63n(int64(max)))
	}
	return PartitionKeyGenerator{GenFunc:genFunc, Prefix:prefix}
}

// TODO func Normal()
