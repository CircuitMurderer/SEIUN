package cc

import (
	"encoding/hex"
	"fmt"

	"github.com/hashicorp/vault/shamir"
)

func SecretDistribute(secret string, m, n int) ([]string, error) {
	if len(secret) == 0 {
		return nil, fmt.Errorf("secret empty")
	}
	if m > n {
		return nil, fmt.Errorf("m cannot be larger than n")
	}

	secretByte := []byte(secret) //hex.DecodeString(secret)
	sharesByte, err := shamir.Split(secretByte, n, m)
	if err != nil {
		return nil, err
	}

	shares := make([]string, len(sharesByte))
	for i, shareByte := range(sharesByte) {
		shares[i] = hex.EncodeToString(shareByte)
	}

	return shares, nil
}

func SecretCollect(shares []string) (string, error) {
	sharesByte := make([][]byte, len(shares))
	for i, share := range(shares) {
		tmp, err := hex.DecodeString(share)
		if err != nil {
			return "", err
		}

		sharesByte[i] = tmp
	}

	secretByte, err := shamir.Combine(sharesByte)
	if err != nil {
		return "", err
	}

	return string(secretByte), nil
}
