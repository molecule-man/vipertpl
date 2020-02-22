# vipertpl

[![GoDoc](https://godoc.org/github.com/molecule-man/vipertpl?status.svg)](https://godoc.org/github.com/molecule-man/vipertpl)
[![CircleCI](https://circleci.com/gh/molecule-man/vipertpl.svg?style=svg)](https://circleci.com/gh/molecule-man/vipertpl)

Package `vipertpl` extends [viper's](https://github.com/spf13/viper)
functionality with ability to use golang
[templates](https://golang.org/pkg/text/template/) in string variables.

## Usage example

```go
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
```

## Template funcs

### ViperGet

`ViperGet` is a built-in function, a wrapper over `viper.Get` function which
applies template parsing to the output of `viper.Get` (can be used recursively.
See the following example).

```yaml
foo: "foo value"
bar: 'bar value + {{ ViperGet "foo" }}'
buz: 'buz value + {{ ViperGet "bar" }}'
```

```go
// ... read config with viper ...
vipertpl.Parse(viper.GetViper())
fmt.Printf("%#v", viper.Get("buz"))
// Output: "buz value + bar val + foo val"
```

### Docker port

It is possible to invoke command [docker port CONTAINER
PRIVATE_PORT](https://docs.docker.com/engine/reference/commandline/port/) inside
the template:

```yaml
database:
  dsn: 'user:password@tcp({{ DockerPort "myservice.mysql" 3306 }})
```

```go
import (
	"text/template"

	"github.com/molecule-man/vipertpl"
	"github.com/molecule-man/vipertpl/docker"
	"github.com/spf13/viper"
)

// ... read config with viper ...
// ...

dockerFuncs, err := docker.New()
if err != nil {
	panic(err)
}

// as the `docker port` function is not built-in we must initialize parser with
// this function added to a list of available template functions:
parser := vipertpl.New(template.FuncMap{
	"DockerPort": dockerFuncs.Port,
})

parser.Parse(viper.GetViper())
fmt.Printf("%#v", viper.Get("database.dsn"))
// Given that there is docker container running under name myservice.mysql with
// private port published at port 32173 the output will be
// "user:password@tcp(0.0.0.0:32173)"
```
