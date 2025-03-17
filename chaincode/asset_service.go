package chaincode

import (
	"form-chaincode/utils"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) exists(context contractapi.TransactionContextInterface, id string) bool {
	value, err := context.GetStub().GetState(id)
	if utils.ValueExists(err, value) {
		return true
	}
	return false
}
