package cc

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

type TCI contractapi.TransactionContextInterface

type SmartContract struct {
	contractapi.Contract
}

type CertItem struct {
	ID      string `json:"ID"`
	Title   string `json:"Title"`
	Owner   string `json:"Owner"`
	Kind    string `json:"Kind"`
	Family  string `json:"Family"`
	Info    string `json:"Info"`
	Status  string `json:"Status"`
	Reserve string `json:"Reserve"`
}

type CItem struct {
	ID      string `json:"ID"`
	UserID  string `json:"UsrID"`
	Status  int    `json:"Status"`
	IsuTime string `json:"IsuTime"`
	RvkTime string `json:"RvkTime"`

	Key    string            `json:"Key"`
	Shares map[string]string `json:"Shares"`
}

type AlivePeers struct {
	PeerInfo map[string]string `json:"PeerInfo"`
}

type WorldState struct {

}

const (
	ValidCert = iota
	InvalidCert
	UnauthedCert
	OutdatedCert
	RevokedCert
	OtherStatus
)
