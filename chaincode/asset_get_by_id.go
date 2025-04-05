package chaincode

import (
	"encoding/json"
	"fmt"
	"form-chaincode/dtos"
	"form-chaincode/utils"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) GetAssetById(context contractapi.TransactionContextInterface, id string) (string, error) {
	clearId, err := s.validateGetAssetByIdData(context, id)
	if err != nil {
		return "", err
	}

	asset, err := s.getDataFromLedgerById(context, clearId)
	if err != nil {
		return "", err
	}

	assetEncoded, err := json.Marshal(asset)
	if err != nil {
		return "", fmt.Errorf("error encoding the object %s", err.Error())
	}

	return string(assetEncoded), nil
}

func (s *SmartContract) validateGetAssetByIdData(context contractapi.TransactionContextInterface, id string) (string, error) {
	cleanId := utils.RemoveStringSpaces(id)

	if !utils.IsValidString(cleanId) {
		return "", fmt.Errorf("the id is not valid")
	}

	if !s.exists(context, cleanId) {
		return "", fmt.Errorf("the asset doesn't exist")
	}

	return cleanId, nil
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
