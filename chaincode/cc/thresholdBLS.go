package cc

import (
	"math"

	"go.dedis.ch/kyber/v3/pairing/bn256"
	"go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/sign/bdn"
	"go.dedis.ch/kyber/v3/sign/tbls"
)

type KeyPair struct {
	PriKey *share.PriShare
	PubKey *share.PubPoly
	Sign   []byte
}

func GenKeys(allPeers map[string]string, suite *bn256.Suite, data string, t float64) (map[string]*KeyPair, error) {
	n := len(allPeers)
	m := int(math.Ceil(float64(n) * t))

	msg := []byte(data)
	keyPairs := make(map[string]*KeyPair, 0)
	secret := suite.G1().Scalar().Pick(suite.RandomStream())

	priPoly := share.NewPriPoly(suite.G2(), m, secret, suite.RandomStream())
	pubPoly := priPoly.Commit(suite.G2().Point().Base())

	i := 0
	priShares := priPoly.Shares(n)
	for peer := range allPeers {
		keyPairs[peer] = new(KeyPair)

		keyPairs[peer].PriKey = priShares[i]
		keyPairs[peer].PubKey = pubPoly

		sign, err := tbls.Sign(suite, keyPairs[peer].PriKey, msg)
		if err != nil {
			return nil, err
		}

		keyPairs[peer].Sign = sign
		i += 1
	}

	return keyPairs, nil
}

func SignVerify(alivePeers map[string]string, suite *bn256.Suite, keyPairs map[string]*KeyPair, data string, t float64) error {
	n := len(keyPairs)
	m := int(math.Ceil(float64(n) * t))
	
	msg := []byte(data)
	signs := make([][]byte, 0)
	for peer, keyPair := range keyPairs {
		_, exists := alivePeers[peer]
		if !exists {
			continue
		}

		signs = append(signs, keyPair.Sign)
	}

	for peer, keyPair := range keyPairs {
		_, exists := alivePeers[peer]
		if !exists {
			continue
		}

		sign, err := tbls.Recover(suite, keyPair.PubKey, msg, signs, m, n)
		if err != nil {
			return err
		}

		err = bdn.Verify(suite, keyPair.PubKey.Commit(), msg, sign)
		if err != nil {
			return err
		}
	}

	return nil
}

func TBLSVerify(data string, allPeers map[string]string, alivePeers map[string]string, threshold float64) error {
	suite := bn256.NewSuite()
	keyPairs, err := GenKeys(allPeers, suite, data, threshold)
	if err != nil {
		return err
	}

	err = SignVerify(alivePeers, suite, keyPairs, data, threshold)
	if err != nil {
		return err
	}

	return nil
}
