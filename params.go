package squid

import "github.com/flosch/pongo2"

type Params map[string]interface{}

func (any Params) Update() pongo2.Context {
	result := pongo2.Context{}
	for k, v := range any {
		result[k] = v
	}
	return result
}