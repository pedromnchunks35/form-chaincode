package main

import (
	"fmt"
	"form-chaincode/chaincode"
	"form-chaincode/utils"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
)

type serverConfig struct {
	CCID    string
	Address string
}

func main() {
	config := serverConfig{
		CCID:    utils.GetEnvOrDefault("CHAINCODE_ID", "123456"),
		Address: utils.GetEnvOrDefault("CHAINCODE_SERVER_ADDRESS", "localhost:8080"),
	}
	assetChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating basic chaincode: %v", err)
	}

	server := &shim.ChaincodeServer{
		CCID:     config.CCID,
		Address:  config.Address,
		CC:       assetChaincode,
		TLSProps: utils.GetTLSProperties(),
	}

	fmt.Printf("Starting server for chaincode %v in address: %v \n", config.CCID, config.Address)
	if err := server.Start(); err != nil {
		log.Fatalf("Something went wrong starting the chaincode: %v", err)
	}
}
