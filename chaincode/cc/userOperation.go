package cc

import "encoding/json"

func (s *SmartContract) SubmitReq(ctx TCI, id string, userID string) (string, error) {
	/*
	realID := ""
	for {
		realID = "Item-" + GenRandomString(123)
		
		hasItem, err := s.HasItem(ctx, realID)
		if err != nil {
			return "", err
		}
		if !hasItem {
			break
		}
	}
	*/

	realID := "Item-" + id
	key := "TestKey"
	cItem := CItem{
		ID:      realID,
		UserID:  userID,
		Status:  UnauthedCert,
		IsuTime: "",
		RvkTime: "",
		Key:     key,
		Shares:  make(map[string]string, 0),
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

	return key, nil
}
