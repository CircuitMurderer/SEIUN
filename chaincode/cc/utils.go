package cc

import (
	"context"
	"crypto/sha256"
	"math/rand"
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

func GetHashedPeers(peers map[string]string) (map[string]string, error) {
	hashedPeers := make(map[string]string, 0)
	hasher := sha256.New()
	for peer, name := range peers {
		_, err := hasher.Write([]byte(peer))
		if err != nil {
			return nil, err
		}

		hashed := hasher.Sum(nil)
		hashedPeers[string(hashed)] = name
		hasher.Reset()
	}

	return hashedPeers, nil
}

func GenRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	letters := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	
	strByte := make([]byte, length)
	for i := range strByte {
		strByte[i] = letters[rand.Intn(len(letters))]
	}

	return string(strByte)
}
