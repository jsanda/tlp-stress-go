package generators

import (
	"math/rand"
	"strconv"
	"time"
)

type random struct {
	min int64
	max int64
	seededRand *rand.Rand
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewRandom(min int64, max int64) FieldGenerator {
	return &random{
		min: min,
		max: max,
		seededRand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *random) SetParameters(params []string) error {
	min, err := strconv.ParseInt(params[0], 10, 64)
	if err != nil {
		return err
	}
	max, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return err
	}
	r.min = min
	r.max = max

	return nil
}

func (r *random) GetInt() int64 {
	return r.randomInt()
}

func (r *random) GetText() string {
	// This code is shamelessly taken from https://www.calhoun.io/creating-random-strings-in-go/
	// John Calhoun has some great tutorials and exercises!
	length := r.randomInt()
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (r *random) GetDescription() string {
	return "Completely random data with even distribution"
}

func (r *random) randomInt() int64 {
	return r.seededRand.Int63n(r.max - r.min) + r.min
}
