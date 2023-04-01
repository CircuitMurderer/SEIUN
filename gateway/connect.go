package main

import (
	"crypto/x509"
	"fmt"
	"os"
	"path"

	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/yaml.v3"
)

type Connector struct {
	MspID       string `yaml:"mspID"`
	CryptoPath  string `yaml:"cryptoPath"`
	CertPath    string `yaml:"certPath"`
	KeyPath     string `yaml:"keyPath"`
	TlsCertPath string `yaml:"tlsCertPath"`
	PeerPoint   string `yaml:"peerPoint"`
	GatewayPeer string `yaml:"gatewayPeer"`
	Channel     string `yaml:"channel"`
	CcType      string `yaml:"ccType"`
}

func createConnFrom(filePath string) Connector {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	c := Connector{}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		panic(err)
	}

	c.CertPath = c.CryptoPath + c.CertPath
	c.KeyPath = c.CryptoPath + c.KeyPath
	c.TlsCertPath = c.CryptoPath + c.TlsCertPath

	return c
}

func loadCert(filePath string) (*x509.Certificate, error) {
	certPEM, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cert: %v", err)
	}
	return identity.CertificateFromPEM(certPEM)
}

func (c *Connector) createConnect() *grpc.ClientConn {
	cert, err := loadCert(c.TlsCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(cert)
	tpCred := credentials.NewClientTLSFromCert(certPool, c.GatewayPeer)

	conn, err := grpc.Dial(c.PeerPoint, grpc.WithTransportCredentials(tpCred))
	if err != nil {
		panic(fmt.Errorf("failed to create connection: %v", err))
	}

	return conn
}

func (c *Connector) createID() *identity.X509Identity {
	cert, err := loadCert(c.CertPath)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(c.MspID, cert)
	if err != nil {
		panic(err)
	}

	return id
}

func (c *Connector) createSign() identity.Sign {
	files, err := os.ReadDir(c.KeyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read key dir: %v", err))
	}

	pKeyPEM, err := os.ReadFile(path.Join(c.KeyPath, files[0].Name()))
	if err != nil {
		panic(fmt.Errorf("failed to read key file: %v", err))
	}

	pKey, err := identity.PrivateKeyFromPEM(pKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(pKey)
	if err != nil {
		panic(err)
	}

	return sign
}
