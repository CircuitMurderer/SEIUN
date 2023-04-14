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

func genKeys(m, n int, suite *bn256.Suite, data string) ([]KeyPairs, [][]byte, error) {
	msg := []byte(data)
	keyPairs := make([]KeyPairs, n)
	secret := suite.G1().Scalar().Pick(suite.RandomStream())
	
	priPoly := share.NewPriPoly(suite.G2(), m, secret, suite.RandomStream())
	pubPoly := priPoly.Commit(suite.G2().Point().Base())
	signs := make([][]byte, 0)

	for i, pShare := range priPoly.Shares(n) {
		keyPairs[i].PriKey = pShare
		keyPairs[i].PubKey = pubPoly

		sign, err := tbls.Sign(suite, keyPairs[i].PriKey, msg)
		if err != nil {
			return nil, nil, err
		}

		signs = append(signs, sign)
	}

	return keyPairs, signs, nil
}

func signVerify(m, n int, suite *bn256.Suite, keyPairs []KeyPairs, signs [][]byte, data string) error {
	msg := []byte(data)

	for _, keyPair := range keyPairs {
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

	return nil
}