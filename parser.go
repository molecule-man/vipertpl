// Package vipertpl extends viper's functionality with ability to use golang
// templates in string variables
package vipertpl

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	"github.com/spf13/viper"
)

// Parse goes though all the keys stored in viper config and parses them with
// golang's internal templating engine
func Parse(v *viper.Viper) error {
	for _, key := range v.AllKeys() {
		parser := parser{v, map[string]struct{}{}}
		val := v.Get(key)

		_, err := parser.parseTpl(key, val)
		if err != nil {
			return err
		}
	}

	return nil
}

type parser struct {
	v           *viper.Viper
	visitedKeys map[string]struct{}
}

func (p *parser) viperGet(key string) (interface{}, error) {
	if _, isVisited := p.visitedKeys[key]; isVisited {
		return "", fmt.Errorf("not able to parse tpl for key %s: %w", key, ErrCircularDependency)
	}

	val := p.v.Get(key)
	p.visitedKeys[key] = struct{}{}

	return p.parseTpl(key, val)
}

func (p *parser) parseTpl(key string, rawVal interface{}) (interface{}, error) {
	val, ok := rawVal.(string)
	if !ok {
		return rawVal, nil
	}

	tpl, err := template.New(val).
		Funcs(template.FuncMap{
			"ViperGet": p.viperGet,
		}).
		Parse(val)
	if err != nil {
		return val, err
	}

	var buf bytes.Buffer

	if err := tpl.Execute(&buf, nil); err != nil {
		return val, err
	}

	p.v.Set(key, buf.String())

	return buf.String(), nil
}

// ErrCircularDependency is returned when there is a circular dependency caused
// by using tpl "ViperGet" function
var ErrCircularDependency = errors.New("circular dependency")
