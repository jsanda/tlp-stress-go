package profiles

import "github.com/gocql/gocql"

type keyValue struct {}

func NewKeyValue() StressProfile {
	return &keyValue{}
}

func (k keyValue) Prepare(session *gocql.Session) error {
	return nil
}

func (k keyValue) Schema() []string {
	return make([]string, 1)
}

func (k keyValue) GetRunner() StressRunner {
	return nil
}
