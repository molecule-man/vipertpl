package docker

import (
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/matryer/is"
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

	url, found := docker.Port("hello", 80)
	is.True(found)
	is.Equal("docker.host:12345", url)

	url, found = docker.Port("hello", 8080)
	is.True(found)
	is.Equal("docker.host:32345", url)

	url, found = docker.Port("world", 80)
	is.True(found)
	is.Equal("docker.host2:18888", url)

	_, found = docker.Port("non.existent", 80)
	is.True(!found)

	_, found = docker.Port("hello", 3333)
	is.True(!found)
}
