package chaincode

import (
	"encoding/json"
	"fmt"
	"form-chaincode/dtos"
	"form-chaincode/utils"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) PatchAsset(context contractapi.TransactionContextInterface, encodedData []byte, id string) (*dtos.AssetRequest, error) {
	clearId, clearRequest, err := s.validatePatchData(context, encodedData, id)
	if err != nil {
		return nil, err
	}

	return s.patchAsset(context, clearRequest, clearId)
}

func (s *SmartContract) patchAsset(context contractapi.TransactionContextInterface, request *dtos.PutAssetRequest, clearId string) (*dtos.AssetRequest, error) {
	asset, err := s.GetAssetById(context, clearId)
	if err != nil {
		return nil, err
	}

	if utils.IsValidString(request.Hash) {
		asset.Hash = request.Hash
	}

	if utils.IsValidString(request.TypeForm) {
		asset.TypeForm = request.TypeForm
	}

	if utils.IsValidString(request.Timestamp) {
		asset.Timestamp = request.Timestamp
	}

	if utils.IsValidString(request.InsertionType) {
		asset.InsertionType = request.InsertionType
	}

	if utils.IsValidString(request.Description) {
		asset.Description = request.Description
	}

	encodedData, err := json.Marshal(asset)
	if err != nil {
		return nil, fmt.Errorf("error encoding asset after changing values %s", err)
	}

	err = context.GetStub().PutState(clearId, encodedData)
	if err != nil {
		return nil, fmt.Errorf("error updating ledger %s", err)
	}

	return asset, nil
}

func (s *SmartContract) validatePatchData(context contractapi.TransactionContextInterface, encodedData []byte, id string) (string, *dtos.PutAssetRequest, error) {
	clearId := utils.RemoveStringSpaces(id)
	if !utils.IsValidString(clearId) {
		return "", nil, fmt.Errorf("the id is not valid")
	}

	if !s.exists(context, id) {
		return "", nil, fmt.Errorf("it doesn't exit")
	}

	request := &dtos.PutAssetRequest{}
	err := json.Unmarshal(encodedData, request)
	if err != nil {
		return "", nil, fmt.Errorf("decoding the object %s", err)
	}

	if removeSpacesAndCheckIfOnePropertyToChange(request) {
		return "", nil, fmt.Errorf("nothing to change in the request")
	}

	return clearId, request, nil
}

func removeSpacesAndCheckIfOnePropertyToChange(request *dtos.PutAssetRequest) bool {
	request.Timestamp = utils.RemoveStringSpaces(request.Timestamp)
	request.TypeForm = utils.RemoveStringSpaces(request.TypeForm)
	request.Hash = utils.RemoveStringSpaces(request.Hash)
	request.InsertionType = utils.RemoveStringSpaces(request.InsertionType)
	request.Description = utils.RemoveStringSpaces(request.Description)

	return !utils.IsValidString(request.Timestamp) &&
		!utils.IsValidString(request.TypeForm) &&
		!utils.IsValidString(request.Hash) &&
		!utils.IsValidString(request.InsertionType) &&
		!utils.IsValidString(request.Description)
}
