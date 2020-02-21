// Package docker contains template functions that use docker client to retrieve
// info about running docker containers
package docker

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// New initializes new Docker struct
func New() (*Docker, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	cc, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	return newDocker(cc), nil
}

func newDocker(cc []types.Container) *Docker {
	ports := make(map[cntPort]string, len(cc))

	for _, c := range cc {
		for _, port := range c.Ports {
			if port.PublicPort == 0 {
				continue
			}

			for _, name := range c.Names {
				name = strings.TrimPrefix(name, "/")
				ports[cntPort{name, port.PrivatePort}] = fmt.Sprintf("%s:%d", port.IP, port.PublicPort)
			}
		}
	}

	return &Docker{ports}
}

// Docker struct holds docker-related template functions
type Docker struct {
	ports map[cntPort]string
}

// Port is a template function that mimics
// `docker port <container name> <private port>` function
func (d *Docker) Port(containerName string, port uint16) (string, bool) {
	url, ok := d.ports[cntPort{containerName, port}]
	return url, ok
}

type cntPort struct {
	name string
	port uint16
}
