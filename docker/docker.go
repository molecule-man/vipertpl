// Package docker contains template functions that use docker client to retrieve
// info about running docker containers
package docker

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// New initializes new Docker struct.
func New() (*Docker, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	cli.NegotiateAPIVersion(ctx)

	cc, err := cli.ContainerList(ctx, container.ListOptions{})
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
				ip := net.ParseIP(port.IP)
				if ip == nil || ip.To4() == nil {
					continue
				}

				name = strings.TrimPrefix(name, "/")
				ports[cntPort{name, port.PrivatePort}] = fmt.Sprintf("%s:%d", port.IP, port.PublicPort)
			}
		}
	}

	return &Docker{ports}
}

// Docker struct holds docker-related template functions.
type Docker struct {
	ports map[cntPort]string
}

// Port is a template function that mimics
// `docker port <container name> <private port>` function.
func (d *Docker) Port(containerName string, port uint16) (string, error) {
	url, ok := d.ports[cntPort{containerName, port}]
	if !ok {
		return "", fmt.Errorf("Docker.Port failed as container %s (port %d) is not found", containerName, port)
	}

	return url, nil
}

type cntPort struct {
	name string
	port uint16
}
