package profiles

import (
	"github.com/jsanda/tlp-stress-go/pkg/generators"
)

type StressProfile interface {
	// gocql automatically prepares queries so we do not need to port this
	//Prepare(session *gocql.Session) error

	Schema() []string

	GetRunner() StressRunner

	GetFieldGenerators() map[*generators.Field]generators.FieldGenerator
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

func GetPlugin(name string) (*Plugin, bool) {
	plugin, ok := plugins[name]
	return &plugin, ok
}
