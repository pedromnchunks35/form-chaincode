package chaincode

import (
	"encoding/json"
	"fmt"
	"form-chaincode/dtos"
	"form-chaincode/utils"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) GetAssetById(context contractapi.TransactionContextInterface, id string) (*dtos.AssetRequest, error) {
	clearId, err := s.validateGetAssetByIdData(context, id)
	if err != nil {
		return nil, err
	}

	return s.getDataFromLedgerById(context, clearId)
}

func (s *SmartContract) validateGetAssetByIdData(context contractapi.TransactionContextInterface, id string) (string, error) {
	if !utils.IsValidString(id) {
		return "", fmt.Errorf("the id is not valid")
	}

	if !s.exists(context, id) {
		return "", fmt.Errorf("the asset doesn't exist")
	}

	return utils.RemoveStringSpaces(id), nil
}

func (s *SmartContract) getDataFromLedgerById(context contractapi.TransactionContextInterface, clearId string) (*dtos.AssetRequest, error) {
	encodedData, err := context.GetStub().GetState(clearId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving data from ledger")
	}

	data := &dtos.AssetRequest{}
	err = json.Unmarshal(encodedData, data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshling data")
	}
	return data, nil
}
