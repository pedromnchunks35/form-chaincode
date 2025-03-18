package chaincode

import (
	"fmt"
	"form-chaincode/utils"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) DeleteAssetById(context contractapi.TransactionContextInterface, id string) (bool, error) {
	clearId, err := s.validateDataDeleteById(context, id)
	if err != nil {
		return false, err
	}

	return s.deleteDataFromLedgerById(context, clearId)
}

func (s *SmartContract) validateDataDeleteById(context contractapi.TransactionContextInterface, id string) (string, error) {
	clearId := utils.RemoveStringSpaces(id)
	if !utils.IsValidString(clearId) {
		return "", fmt.Errorf("id is not valid")
	}

	if !s.exists(context, clearId) {
		return "", fmt.Errorf("asset does't exist")
	}

	return clearId, nil
}

func (s *SmartContract) deleteDataFromLedgerById(context contractapi.TransactionContextInterface, clearId string) (bool, error) {
	err := context.GetStub().DelState(clearId)
	if err != nil {
		return false, fmt.Errorf("error deleting state from the ledger %s", err)
	}

	return true, nil
}
