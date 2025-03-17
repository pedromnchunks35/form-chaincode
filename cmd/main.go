package main

import (
	"fmt"
	"form-chaincode/utils"
	chaincode "github.com/hyperledger/fabric-chaincode-go"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
	"os"
)

type serverConfig struct {
	CCID    string
	Address string
}

func main() {
	config := serverConfig{
		CCID:    os.Getenv("CHAINCODE_ID"),
		Address: os.Getenv("CHAINCODE_SERVER_ADDRESS"),
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
