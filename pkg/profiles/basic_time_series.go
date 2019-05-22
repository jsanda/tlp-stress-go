package profiles

import "github.com/gocql/gocql"

type basicTimeSeries struct{}

func NewBasicTimeSeries() StressProfile {
	return &basicTimeSeries{}
}

func (b basicTimeSeries) Prepare(session *gocql.Session) error {
	return nil
}

func (b basicTimeSeries) Schema() []string {
	return make([]string, 1)
}

func (b basicTimeSeries) GetRunner() StressRunner {
	return nil
}
