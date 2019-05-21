package profiles

import "github.com/gocql/gocql"

type StressProfile interface {
	Prepare(session *gocql.Session) error

	Schema() []string

	GetRunner() StressRunner
}

type StressRunner interface {

}
