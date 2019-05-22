package profiles

import "github.com/gocql/gocql"

type StressProfile interface {
	Prepare(session *gocql.Session) error

	Schema() []string

	GetRunner() StressRunner
}

type StressRunner interface {

}

type Plugin struct {
	Name string
	Instance StressProfile
}

var plugins = map[string]Plugin{
	"BasicTimeSeries": {
		Name: "BasicTimeSeries",
		Instance: NewBasicTimeSeries(),
	},
	"KeyValue": {
		Name: "KeyValue",
		Instance: NewKeyValue(),
	},
}

func GetPlugins() map[string]Plugin {
	return plugins
}
