package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

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

func (s *SmartContract) InitLedger(ctx TCI) error {
	items := []CertItem{
		{ID: "CERTNUM1", Title: "Outstanding Student", Owner: "Yuki", Kind: "Honor",
			Family: "Wuhan University", Info: "So great", Status: "valid", Reserve: ""},
		{ID: "CERTNUM2", Title: "Scholarship", Owner: "Aoki", Kind: "Award",
			Family: "Huazhong University", Info: "Great", Status: "valid", Reserve: ""},
	}

	for _, item := range items {
		itemJSON, err := json.Marshal(item)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(item.ID, itemJSON)
		if err != nil {
			return fmt.Errorf("failed to put world state. %v", err)
		}
	}
	return nil
}

func (s *SmartContract) HasItem(ctx TCI, id string) (bool, error) {
	itemJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read world state. %v", err)
	}

	return itemJSON != nil, nil
}

func (s *SmartContract) GetItem(ctx TCI, id string) (*CertItem, error) {
	itemJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get world state. %v", err)
	}
	if itemJSON == nil {
		return nil, fmt.Errorf("no such item. %s", id)
	}

	certItem := CertItem{}
	err = json.Unmarshal(itemJSON, &certItem)
	if err != nil {
		return nil, err
	}

	return &certItem, nil
}

func (s *SmartContract) AddItem(ctx TCI,
	id string, title string, owner string, kind string,
	family string, info string, status string, reserve string) error {

	hasItem, err := s.HasItem(ctx, id)
	if err != nil {
		return err
	}
	if hasItem {
		return fmt.Errorf("item %s already exists", id)
	}

	certItem := CertItem{
		ID:      id,
		Title:   title,
		Owner:   owner,
		Kind:    kind,
		Family:  family,
		Info:    info,
		Status:  status,
		Reserve: reserve,
	}
	itemJSON, err := json.Marshal(certItem)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, itemJSON)
}

func (s *SmartContract) DelItem(ctx TCI, id string) error {
	hasItem, err := s.HasItem(ctx, id)
	if err != nil {
		return err
	}
	if !hasItem {
		return fmt.Errorf("item %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) ChangeStatus(ctx TCI, id string, status string) (string, error) {
	certItem, err := s.GetItem(ctx, id)
	if err != nil {
		return "", err
	}

	oldStatus := certItem.Status
	certItem.Status = status

	itemJSON, err := json.Marshal(certItem)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(id, itemJSON)
	if err != nil {
		return "", err
	}

	return oldStatus, nil
}

func (s *SmartContract) RenewItem(ctx TCI,
	id string, title string, owner string, kind string,
	family string, info string, status string, reserve string) error {

	hasItem, err := s.HasItem(ctx, id)
	if err != nil {
		return err
	}
	if !hasItem {
		return fmt.Errorf("item %s does not exist", id)
	}

	certItem := CertItem{
		ID:      id,
		Title:   title,
		Owner:   owner,
		Kind:    kind,
		Family:  family,
		Info:    info,
		Status:  status,
		Reserve: reserve,
	}
	itemJSON, err := json.Marshal(certItem)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, itemJSON)
}

func (s *SmartContract) GetAllItems(ctx TCI) ([]*CertItem, error) {
	itemsIter, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer itemsIter.Close()

	items := make([]*CertItem, 0)
	for itemsIter.HasNext() {
		item, err := itemsIter.Next()
		if err != nil {
			return nil, err
		}

		certItem := CertItem{}
		err = json.Unmarshal(item.Value, &certItem)
		if err != nil {
			return nil, err
		}
		items = append(items, &certItem)
	}

	return items, nil
}

func main() {
	cc, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Failed to create chaincode: %v", err)
	}
	if err = cc.Start(); err != nil {
		log.Panicf("Failed to start chaincode: %v", err)
	}
}
