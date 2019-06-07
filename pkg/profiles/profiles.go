package profiles

import (
	"github.com/gocql/gocql"
	"github.com/jsanda/tlp-stress-go/pkg/generators"
)

type PopulationOption interface {}

type Standard struct{
	populationOption PopulationOption
}

type Custom struct{
	populationOption PopulationOption
	Rows uint64
}

type Operation struct {
	Query *gocql.Query
}

type Mutation struct {
	Operation
}

type Select struct {
	Operation
}

type StressProfile interface {
	// gocql automatically prepares queries so we do not need to port this
	//Prepare(session *gocql.Session) error

	Schema() []string

	GetRunner(registry *generators.Registry, session *gocql.Session) StressRunner

	GetFieldGenerators() map[generators.Field]generators.FieldGenerator

	GetPopulationOption() PopulationOption
}

type StressRunner interface {
	GetNextMutation(key *generators.PartitionKey) *Mutation

	GetNextSelect(key *generators.PartitionKey) *Select
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
