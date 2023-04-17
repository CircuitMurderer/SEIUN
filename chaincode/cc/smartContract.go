package cc

import (
	"encoding/json"
	"fmt"
)

func (s *SmartContract) GetCert(ctx TCI, id string) (*CItem, error) {
	itemJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get world state. %v", err)
	}
	if itemJSON == nil {
		return nil, fmt.Errorf("no such item. %s", id)
	}

	certItem := CItem{}
	err = json.Unmarshal(itemJSON, &certItem)
	if err != nil {
		return nil, err
	}

	return &certItem, nil
}

func (s *SmartContract) GetWaitingList(ctx TCI) ([]string, error) {
	waitingJSON, err := ctx.GetStub().GetState("WaitingList")
	if err != nil {
		return nil, fmt.Errorf("failed to get world state. %v", err)
	}

	waitingList := make([]string, 0)
	err = json.Unmarshal(waitingJSON, &waitingList)
	if err != nil {
		return nil, err
	}

	return waitingList, nil
}

func (s *SmartContract) PutWaitingList(ctx TCI, waitingList []string) error {
	waitingJSON, err := json.Marshal(waitingList)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState("WaitingList", waitingJSON)
	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) AddToWaitingList(ctx TCI, id string) error {
	waitingList, err := s.GetWaitingList(ctx)
	if err != nil {
		return err
	}

	waitingList = append(waitingList, id)
	err = s.PutWaitingList(ctx, waitingList)
	if err != nil {
		return err
	}

	return nil
}