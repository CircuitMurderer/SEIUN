package cc

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func GetAlivePeers() ([]string, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.WithHost("tcp://172.21.0.1:2375"), client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	// containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	res := make([]string, 0)
	for _, container := range containers {
		if strings.HasSuffix(container.Names[0], "seiun.net") {
			res = append(res, container.ID)
		}
	}

	return res, nil
}
