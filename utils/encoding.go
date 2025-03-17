package utils

import (
	"encoding/json"
	"form-chaincode/dtos"
)

func DecodeValueToPostRequest(value []byte) (*dtos.PostAssetRequest, error) {
	newDto := &dtos.PostAssetRequest{}
	err := json.Unmarshal(value, newDto)
	if err != nil {
		return nil, err
	}
	return newDto, nil
}
