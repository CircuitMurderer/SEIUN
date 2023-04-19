package cc

import (
	"encoding/json"
	"fmt"
	"time"
)

func (s *SmartContract) VerifyCert(ctx TCI, id string, status string, expDays int) error {
	certItem, err := s.GetCert(ctx, id)
	if err != nil {
		return err
	}

	switch status {
	case "Valid":
		certItem.Status = ValidCert
		certItem.ExpDays = expDays
		certItem.IsuTime = time.Now().Format("2006-01-02 15:04:05")
	case "Invalid":
		certItem.Status = InvalidCert
	case "Outdate":
		certItem.Status = OutdatedCert
	case "Revoke":
		certItem.Status = RevokedCert
		certItem.RvkTime = time.Now().Format("2006-01-02 15:04:05")
	default:
		return fmt.Errorf("no such status")
	}

	itemJSON, err := json.Marshal(certItem)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, itemJSON)
	if err != nil {
		return err
	}

	return nil
}
