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
		certItem.ExpDays = -1
	case "Outdate":
		certItem.Status = OutdatedCert
		certItem.ExpDays = 0
	case "Revoke":
		certItem.Status = RevokedCert
		certItem.ExpDays = -1
		certItem.RvkTime = time.Now().Format("2006-01-02 15:04:05")
	default:
		return fmt.Errorf("unknown status")
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
