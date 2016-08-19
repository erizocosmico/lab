package lab

// Params will hold all the parameters to call a strategy.
type Params interface {
	// Get returns a parameter by its name.
	Get(string) (interface{}, bool)
	// Set sets a parameter.
	Set(string, interface{})
}

type params map[string]interface{}

func (c params) Get(name string) (interface{}, bool) {
	v, ok := c[name]
	return v, ok
}

func (c params) Set(name string, value interface{}) {
	c[name] = value
}

func newParams() Params {
	return make(params)
}
