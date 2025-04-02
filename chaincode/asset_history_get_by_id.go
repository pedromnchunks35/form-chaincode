package chaincode

import (
	"encoding/json"
	"fmt"
	"form-chaincode/utils"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
)

func (s *SmartContract) GetHistoryAssetById(context contractapi.TransactionContextInterface, id string) (string, error) {
	cleanId := utils.RemoveStringSpaces(id)
	if !s.exists(context, cleanId) {
		return "", fmt.Errorf("the asset doesn't exist")
	}

	assetHistory, err := GetHistoryFromCleanKey(context, cleanId)
	if err != nil {
		return "", err
	}

	return MarshalHistoryAndReturnStringValue(assetHistory)
}

func GetHistoryFromCleanKey(context contractapi.TransactionContextInterface, cleanId string) ([]*queryresult.KeyModification, error) {
	iterator, err := context.GetStub().GetHistoryForKey(cleanId)
	if err != nil {
		return nil, fmt.Errorf("something went wrong getting the item history: %s", err.Error())
	}
	defer iterator.Close()

	assetHistory := []*queryresult.KeyModification{}
	for iterator.HasNext() {
		asset, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("something went wring retriving the next item from the history: %s", err.Error())
		}
		assetHistory = append(assetHistory, asset)
	}

	return assetHistory, nil
}

func MarshalHistoryAndReturnStringValue(assetHistory []*queryresult.KeyModification) (string, error) {
	assetHistoryEncoded, err := json.Marshal(assetHistory)
	if err != nil {
		return "", fmt.Errorf("something went wrong encoding the final result %s", err.Error())
	}

	return string(assetHistoryEncoded), nil
}
