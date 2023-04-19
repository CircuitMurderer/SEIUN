package cc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func (s *SmartContract) SubmitReq(ctx TCI, id string, userID string, key string) (string, error) {
	realID := id
	if !strings.HasPrefix(id, "Item-") {
		realID = "Item-" + id
	}

	hasIt, err := s.hasItem(ctx, realID)
	if err != nil {
		return "", err
	}

	if hasIt {
		return "", fmt.Errorf("this ID has already been used")
	}

	cItem := CertItem{
		ID:      realID,
		UserID:  userID,
		Status:  UnauthedCert,
		ReqTime: time.Now().Format("2006-01-02 15:04:05"),
		IsuTime: "",
		RvkTime: "",
		ExpDays: 0,
		Key:     key,
	}

	allPeers, err := s.GetAllPeers(ctx)
	if err != nil {
		return "", err
	}

	alivePeers, err := GetAlivePeers()
	if err != nil {
		return "", err
	}

	if !reflect.DeepEqual(allPeers, alivePeers) {
		err = TBLSVerify(realID, allPeers, alivePeers, 0.7)
		if err != nil {
			return "", err
		}
	}

	itemJSON, err := json.Marshal(cItem)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(realID, itemJSON)
	if err != nil {
		return "", err
	}

	err = s.AddToWaitingList(ctx, realID)
	if err != nil {
		return "", err
	}

	return realID, nil
}

func (s *SmartContract) UserGetCertKey(ctx TCI, id string) (string, error) {
	certItem, err := s.GetCert(ctx, id)
	if err != nil {
		return "", err
	}

	switch certItem.Status {
	case ValidCert:
	case InvalidCert:
		return "", fmt.Errorf("this cert is invalid")
	case OutdatedCert:
		return "", fmt.Errorf("this cert is outdated")
	case UnauthedCert:
		return "", fmt.Errorf("this cert hasn't been authed")
	case RevokedCert:
		return "", fmt.Errorf("this cert has been revoked on %v", certItem.RvkTime)
	default:
		return "", fmt.Errorf("the status of cert is unknown")
	}

	dateFmt := "2006-01-02 15:04:05"
	daysBtw, err := GetDaysBetween(time.Now().Format(dateFmt), certItem.IsuTime, dateFmt)
	if err != nil {
		return "", err
	}

	if daysBtw > certItem.ExpDays {
		err = s.VerifyCert(ctx, certItem.ID, "Outdate", 0)
		if err != nil {
			return "", err
		}

		return "", fmt.Errorf("this cert is outdated")
	}

	allPeers, err := s.GetAllPeers(ctx)
	if err != nil {
		return "", err
	}

	alivePeers, err := GetAlivePeers()
	if err != nil {
		return "", err
	}

	if !reflect.DeepEqual(allPeers, alivePeers) {
		err = SSSVerify(certItem.Key, allPeers, alivePeers, 0.7)
		if err != nil {
			return "", err
		}
	}

	return certItem.Key, nil
}
