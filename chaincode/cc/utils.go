package cc

import (
	"context"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func GetAlivePeers() (map[string]string, error) {
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

	res := make(map[string]string, 0)
	for _, container := range containers {
		if strings.HasSuffix(container.Names[0], "seiun.net") {
			res[container.ID] = container.Names[0]
		}
	}

	return res, nil
}

func GetDaysBetween(dateS1 string, dateS2 string, format string) (int, error) {
	date1, err := time.Parse(format, dateS1)
	if err != nil {
		return 0, err
	}

	date2, err := time.Parse(format, dateS2)
	if err != nil {
		return 0, err
	}

	return int(date1.Sub(date2).Hours() / 24), nil
}
