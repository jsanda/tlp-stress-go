package profiles

import (
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

func (b *basicTimeSeries) GetRunner() StressRunner {
	return nil
}

func (b *basicTimeSeries) GetFieldGenerators() map[*generators.Field]generators.FieldGenerator {
	return map[*generators.Field]generators.FieldGenerator{
		&generators.Field{Table: "sensor_data", Name: "data"}: generators.NewRandom(100, 200),
	}
}
