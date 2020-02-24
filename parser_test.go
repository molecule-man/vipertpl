package vipertpl

import (
	"bytes"
	"os"
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/viper"
)

func TestBasicTemplating(t *testing.T) {
	is := is.New(t)

	input := []byte(`foo: "{{\"foo_output\"}}"`)

	viper := viper.New()
	viper.SetConfigType("yaml")
	is.NoErr(viper.ReadConfig(bytes.NewBuffer(input)))

	err := Parse(viper)
	is.NoErr(err)

	is.Equal("foo_output", viper.Get("foo"))
}

func TestChaining(t *testing.T) {
	is := is.New(t)

	input := []byte(`
foo: 'foo_val'
nested:
  bar: 'bar_val + {{ ViperGet "foo" }}'
buz: 'buz_val + {{ ViperGet "nested.bar" }}'
`)

	viper := viper.New()
	viper.SetConfigType("yaml")
	is.NoErr(viper.ReadConfig(bytes.NewBuffer(input)))

	err := Parse(viper)
	is.NoErr(err)

	is.Equal("buz_val + bar_val + foo_val", viper.Get("buz"))
}

func TestCircularDependency(t *testing.T) {
	is := is.New(t)

	input := []byte(`
foo: '{{ ViperGet "bar" }}'
bar: '{{ ViperGet "foo" }}'
`)

	viper := viper.New()
	viper.SetConfigType("yaml")
	is.NoErr(viper.ReadConfig(bytes.NewBuffer(input)))

	err := Parse(viper)
	is.True(err != nil)
}

func TestNonString(t *testing.T) {
	is := is.New(t)

	input := []byte(`
number: 42
boolean: true
bar: 'bar + {{ ViperGet "number" }} + {{ ViperGet "boolean" }}'
`)

	viper := viper.New()
	viper.SetConfigType("yaml")
	is.NoErr(viper.ReadConfig(bytes.NewBuffer(input)))

	err := Parse(viper)
	is.NoErr(err)

	is.Equal("bar + 42 + true", viper.Get("bar"))
	is.Equal(true, viper.Get("boolean"))
	is.Equal(42, viper.Get("number"))
}

func TestTplFuncExec(t *testing.T) {
	if os.Getenv("TEST_INT") != "on" {
		t.Skip("To enable this integration test provide `TEST_INT=on` as an env variable")
	}

	is := is.New(t)

	input := []byte(`foo: '{{ Exec "echo" "foo_val" }}'`)

	viper := viper.New()
	viper.SetConfigType("yaml")
	is.NoErr(viper.ReadConfig(bytes.NewBuffer(input)))

	err := Parse(viper)
	is.NoErr(err)

	is.Equal("foo_val", viper.Get("foo"))
}
