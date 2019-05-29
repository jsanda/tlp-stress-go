package generators

type Field struct {
	Table string
	Name string
}

type FieldGenerator interface {
	GetInt() int64

	//GetFloat() float32

	GetText() string

	GetDescription() string

	SetParameters([]string) error
}

type Registry struct {
	defaults map[*Field]*FieldGenerator

	overrides map[*Field]*FieldGenerator
}

func NewRegistry() *Registry {
	return &Registry{
		defaults: make(map[*Field]*FieldGenerator),
		overrides: make(map[*Field]*FieldGenerator),
	}
}

func (r *Registry) SetDefault(field *Field, generator *FieldGenerator) {
	r.defaults[field] = generator
}
