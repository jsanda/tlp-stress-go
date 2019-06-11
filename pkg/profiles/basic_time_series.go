package profiles

import (
	"github.com/gocql/gocql"
	"github.com/jsanda/tlp-stress-go/pkg/generators"
)

type basicTimeSeries struct{}

func NewBasicTimeSeries() StressProfile {
	return &basicTimeSeries{}
}

func (b *basicTimeSeries) Schema() []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS sensor_data (
  sensor_id text,
  timestamp timeuuid,
  data text,
  PRIMARY KEY (sensor_id, timestamp))
  WITH CLUSTERING ORDER BY (timestamp DESC)`}
}

func (b *basicTimeSeries) GetRunner(registry *generators.Registry, session *gocql.Session) StressRunner {
	return &stressRunner{
		dataField: registry.GetGenerator("sensor_data", "data"),
		session: session,
	}
}

func (b *basicTimeSeries) GetFieldGenerators() map[generators.Field]generators.FieldGenerator {
	return map[generators.Field]generators.FieldGenerator{
		generators.Field{Table: "sensor_data", Name: "data"}: generators.NewRandom(100, 200),
	}
}

func (b *basicTimeSeries) GetPopulationOption() PopulationOption {
	return Standard{}
}

type stressRunner struct {
	dataField generators.FieldGenerator
	session   *gocql.Session
}

func (s *stressRunner) GetNextMutation(key *generators.PartitionKey) *Operation {
	data := s.dataField.GetText()
	timestamp := gocql.TimeUUID()
	query := s.session.Query("INSERT INTO sensor_data (sensor_id, timestamp, data) VALUES (?, ?, ?)",
		key.GetText(), timestamp, data)

	return &Operation{query, Mutation}
}

func (s *stressRunner) GetNextSelect(key *generators.PartitionKey) *Operation {
	return nil
}
