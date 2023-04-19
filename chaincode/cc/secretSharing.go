package cc

import (
	"fmt"
	"math"

	"github.com/hashicorp/vault/shamir"
)

func secretDistribute(secret string, allPeers map[string]string, t float64) (map[string]string, error) {
	if len(secret) == 0 {
		return nil, fmt.Errorf("secret empty")
	}

	n := len(allPeers)
	m := int(math.Ceil(float64(n) * t))

	secretByte := []byte(secret) //hex.DecodeString(secret)
	sharesByte, err := shamir.Split(secretByte, n, m)
	if err != nil {
		return nil, err
	}

	i := 0
	shares := make(map[string]string, 0)
	for k := range allPeers {
		shares[k] = string(sharesByte[i])
		i += 1
	}

	return shares, nil
}

func secretCollect(shares map[string]string, alivePeers map[string]string) (string, error) {
	sharesByte := make([][]byte, 0)
	for peer, share := range shares {
		_, exists := alivePeers[peer]
		if !exists {
			continue
		}

		tmp := []byte(share)
		sharesByte = append(sharesByte, tmp)
	}

	secretByte, err := shamir.Combine(sharesByte)
	if err != nil {
		return "", err
	}

	return string(secretByte), nil
}

func SSSVerify(data string, allPeers map[string]string, alivePeers map[string]string, threshold float64) error {
	shares, err := secretDistribute(data, allPeers, threshold)
	if err != nil {
		return err
	}

	finData, err := secretCollect(shares, alivePeers)
	if err != nil {
		return err
	}

	if finData != data {
		return fmt.Errorf("combined key not equals the origin key")
	}

	return nil
}
