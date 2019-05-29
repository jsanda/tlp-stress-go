package profiles

import "github.com/jsanda/tlp-stress-go/pkg/generators"

type keyValue struct {}

func NewKeyValue() StressProfile {
	return &keyValue{}
}

func (k *keyValue) Schema() []string {
	return make([]string, 1)
}

func (k *keyValue) GetRunner() StressRunner {
	return nil
}

func (k *keyValue) GetFieldGenerators() map[*generators.Field]generators.FieldGenerator {
	return nil
}
