package vipertpl_test

import (
	"bytes"
	"fmt"

	"github.com/molecule-man/vipertpl"
	"github.com/spf13/viper"
)

func ExampleParse() {
	input := []byte(`
foo: 'foo_val'
bar: 'bar_val + {{ ViperGet "foo" }}'
`)

	viper.SetConfigType("yaml")

	err := viper.ReadConfig(bytes.NewBuffer(input))
	if err != nil {
		panic(err)
	}

	err = vipertpl.Parse(viper.GetViper())
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", viper.Get("bar"))
	// Output: "bar_val + foo_val"
}
