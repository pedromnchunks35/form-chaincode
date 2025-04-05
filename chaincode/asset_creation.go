package chaincode

import (
	"encoding/json"
	"fmt"
	"form-chaincode/dtos"
	"form-chaincode/utils"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"time"
)

func (s *SmartContract) CreateAsset(context contractapi.TransactionContextInterface, encodedValue string) (string, error) {
	newDto, err := s.validateAsset(context, encodedValue)
	if err != nil {
		return "", err
	}

	asset, err := s.postAsset(context, newDto)
	if err != nil {
		return "", err
	}

	assetEncoded, err := json.Marshal(asset)
	if err != nil {
		return "", fmt.Errorf("error encoding the asset %s", err.Error())
	}

	return string(assetEncoded), nil
}

func (s *SmartContract) postAsset(context contractapi.TransactionContextInterface, cleanDto *dtos.PostAssetRequest) (*dtos.PostAssetRequest, error) {
	encodedCleanDto, err := json.Marshal(cleanDto)
	if err != nil {
		return nil, fmt.Errorf("encoding cleaned object %s", err)
	}

	err = context.GetStub().PutState(cleanDto.Id, encodedCleanDto)
	if err != nil {
		return nil, fmt.Errorf("inserting cleaned object %s", err)
	}

	return cleanDto, nil
}

func (s *SmartContract) validateAsset(context contractapi.TransactionContextInterface, value string) (*dtos.PostAssetRequest, error) {
	newDto, err := utils.DecodeValueToPostRequest(value)
	if err != nil {
		return nil, fmt.Errorf("decoding the given value results in: %s", err)
	}

	if !removeSpacesAndArePostRequestFieldsValid(newDto) {
		return nil, fmt.Errorf("some fields are not valid")
	}

	if s.exists(context, newDto.Id) {
		return nil, fmt.Errorf("already exists")
	}
	return newDto, nil
}

func removeSpacesAndArePostRequestFieldsValid(request *dtos.PostAssetRequest) bool {
	request.Id = utils.RemoveStringSpaces(request.Id)
	request.TypeForm = utils.RemoveStringSpaces(request.TypeForm)
	request.InsertionType = utils.RemoveStringSpaces(request.InsertionType)
	request.Hash = utils.RemoveStringSpaces(request.Hash)

	return areAllPostFieldsValid(
		request.Id,
		request.TypeForm,
		request.Description,
		request.Timestamp,
		request.InsertionType,
		request.Hash,
	)
}

func areAllPostFieldsValid(
	id string,
	typeForm string,
	description string,
	timestamp time.Time,
	insertionsType string,
	hash string,
) bool {
	return utils.IsValidString(id) &&
		utils.IsValidString(typeForm) &&
		utils.IsValidString(description) &&
		!timestamp.IsZero() &&
		utils.IsValidString(insertionsType) &&
		utils.IsValidString(hash)
}
