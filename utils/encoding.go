package utils

import (
	"encoding/json"
	"form-chaincode/dtos"
)

func DecodeValueToPostRequest(value string) (*dtos.PostAssetRequest, error) {
	valueInBytes := []byte(value)
	newDto := &dtos.PostAssetRequest{}
	err := json.Unmarshal(valueInBytes, newDto)
	if err != nil {
		return nil, err
	}
	return newDto, nil
}
