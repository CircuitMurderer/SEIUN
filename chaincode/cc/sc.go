package cc

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

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

func (s *SmartContract) GetAlivePeers() ([]string, error) {
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

	res := make([]string, 0)
	for _, container := range containers {
		if strings.HasSuffix(container.Names[0], "seiun.net") {
			res = append(res, container.ID)
		}
	}

	return res, nil
}

func (s *SmartContract) GetHost() (string, error) {
	out, err := exec.Command("hostname").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}