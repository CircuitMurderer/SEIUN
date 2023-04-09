package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func main() {
	connector := createConnFrom("netConf.yaml")
	clientConn := connector.createConnect()
	defer clientConn.Close()

	gtw, err := client.Connect(
		connector.createID(),
		client.WithSign(connector.createSign()),
		client.WithClientConnection(clientConn),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)

	if err != nil {
		panic(err)
	}
	defer gtw.Close()

	network := gtw.GetNetwork(connector.Channel)
	contract := network.GetContract(connector.CcType)

	fmt.Println("")
	evalQuery(contract, "GetAllItems")
	evalInvoke(contract, "ChangeStatus", "CERTNUM2", "in!valid")
	evalQuery(contract, "GetAllItems")
}

func formatJSON(data []byte) string {
	fmtJSON := bytes.Buffer{}
	err := json.Indent(&fmtJSON, data, " ", "")
	if err != nil {
		panic(fmt.Errorf("failed to parse JSON: %v", err))
	}
	return fmtJSON.String()
}

func evalQuery(cc *client.Contract, fun string, args ...string) {
	evalRes, err := cc.EvaluateTransaction(fun, args...)
	if err != nil {
		panic(fmt.Errorf("failed to eval transaction: %v", err))
	}

	fmt.Println("---Answer---")
	fmt.Println(formatJSON(evalRes))
}

func evalInvoke(cc *client.Contract, fun string, args ...string) {
	evalRes, err := cc.SubmitTransaction(fun, args...)
	if err != nil {
		panic(fmt.Errorf("failed to eval transaction: %v", err))
	}

	fmt.Println(">>Result<<")
	fmt.Println(string(evalRes))
}
