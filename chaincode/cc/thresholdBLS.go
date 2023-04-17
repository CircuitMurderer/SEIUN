package cc

import (
	"go.dedis.ch/kyber/v3/pairing/bn256"
	"go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/sign/bdn"
	"go.dedis.ch/kyber/v3/sign/tbls"
)

type KeyPairs struct {
	PriKey *share.PriShare
	PubKey *share.PubPoly
}

func GenKeys(m, n int, allPeers map[string]string, suite *bn256.Suite, data string) (map[string]*KeyPairs, [][]byte, error) {
	msg := []byte(data)
	keyPairs := make(map[string]*KeyPairs, 0)
	secret := suite.G1().Scalar().Pick(suite.RandomStream())
	
	priPoly := share.NewPriPoly(suite.G2(), m, secret, suite.RandomStream())
	pubPoly := priPoly.Commit(suite.G2().Point().Base())
	signs := make([][]byte, 0)

	i := 0
	priShares := priPoly.Shares(n)
	for peer := range allPeers {
		keyPairs[peer].PriKey = priShares[i]
		keyPairs[peer].PubKey = pubPoly

		sign, err := tbls.Sign(suite, keyPairs[peer].PriKey, msg)
		if err != nil {
			return nil, nil, err
		}

		signs = append(signs, sign)
		i += 1
	}

	return keyPairs, signs, nil
}

func SignVerify(m, n int, alivePeers map[string]string, suite *bn256.Suite, keyPairs map[string]*KeyPairs, signs [][]byte, data string) error {
	msg := []byte(data)

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



func SignID() error {
	/*
	suite := bn256.NewSuite()
	data := "THISISATEST"
	n := 10
	m := 7

	keyPairs, signs, err := genKeys(m, n, suite, data)
	if err != nil {
		return err
	}
	
	err = signVerify(m, n, suite, keyPairs, signs[1:8], data)
	if err != nil {
		return err
	}
	*/

	return nil
}