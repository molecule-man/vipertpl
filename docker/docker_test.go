package docker

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/docker/docker/api/types"
	"github.com/matryer/is"
	"github.com/molecule-man/vipertpl"
	"github.com/spf13/viper"
)

func TestPort(t *testing.T) {
	is := is.New(t)
	docker := newDocker([]types.Container{
		{
			Names: []string{"hello"},
			Ports: []types.Port{
				{IP: "docker.host", PrivatePort: 80, PublicPort: 12345},
				{IP: "docker.host", PrivatePort: 8080, PublicPort: 32345},
			},
		},
		{
			Names: []string{"world"},
			Ports: []types.Port{
				{IP: "docker.host2", PrivatePort: 80, PublicPort: 18888},
			},
		},
	})

	url, err := docker.Port("hello", 80)
	is.NoErr(err)
	is.Equal("docker.host:12345", url)

	url, err = docker.Port("hello", 8080)
	is.NoErr(err)
	is.Equal("docker.host:32345", url)

	url, err = docker.Port("world", 80)
	is.NoErr(err)
	is.Equal("docker.host2:18888", url)

	_, err = docker.Port("non.existent", 80)
	is.True(err != nil)

	_, err = docker.Port("hello", 3333)
	is.True(err != nil)
}

func TestTemplatingWithDockerPort(t *testing.T) {
	is := is.New(t)

	input := []byte(`foo: '{{ DockerPort "testcnt" 80 }}'`)

	viper := viper.New()
	viper.SetConfigType("yaml")
	is.NoErr(viper.ReadConfig(bytes.NewBuffer(input)))

	dockerFuncs := newDocker([]types.Container{
		{
			Names: []string{"testcnt"},
			Ports: []types.Port{
				{IP: "docker.host", PrivatePort: 80, PublicPort: 12345},
			},
		},
	})

	parser := vipertpl.New(template.FuncMap{
		"DockerPort": dockerFuncs.Port,
	})
	err := parser.Parse(viper)
	is.NoErr(err)

	is.Equal("docker.host:12345", viper.Get("foo"))
}
