package chaincode

import (
	"encoding/json"
	"fmt"
	"form-chaincode/dtos"
	"form-chaincode/utils"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strings"
)

func (s *SmartContract) GetAllAssets(
	context contractapi.TransactionContextInterface,
	pageSize string,
	sizeSize string,
	filter string,
) (string, error) {
	page, size, err := validateDataGetAllAssets(pageSize, sizeSize)
	if err != nil {
		return "", err
	}

	query, err := createQuery(filter)
	if err != nil {
		return "", err
	}

	assets, err := s.queryAllSetsWithPagination(context, query, page, size)
	if err != nil {
		return "", err
	}

	encodedAssets, err := json.Marshal(assets)
	if err != nil {
		return "", fmt.Errorf("error encoding the final result %s", err)
	}

	return string(encodedAssets), nil
}

func createQuery(filter string) (string, error) {
	filterByte := []byte(filter)

	filterDecoded := &dtos.Filter{}
	err := json.Unmarshal(filterByte, filterDecoded)
	if err != nil {
		return "", fmt.Errorf("error decoding filter %s", err)
	}

	mainQuery := `{"selector":{`
	initQueryLen := len(mainQuery)
	cleanFilter(filterDecoded)

	if filterDecoded.Hashs != nil {
		encodedArr, err := encodeArray(filterDecoded.Hashs)
		if err != nil {
			return "", err
		}
		queryToAdd := `"hash":{"$in":` + string(encodedArr) + `},`
		mainQuery += queryToAdd
	}

	if filterDecoded.TypeForms != nil {
		encodedArr, err := encodeArray(filterDecoded.TypeForms)
		if err != nil {
			return "", err
		}
		queryToAdd := `"type_form":{"$in":` + string(encodedArr) + `},`
		mainQuery += queryToAdd
	}

	if filterDecoded.InsertionTypes != nil {
		encodedArr, err := encodeArray(filterDecoded.InsertionTypes)
		if err != nil {
			return "", err
		}
		queryToAdd := `"insertion_type":{"$in":` + string(encodedArr) + `},`
		mainQuery += queryToAdd
	}

	if filterDecoded.Ids != nil {
		encodedArr, err := encodeArray(filterDecoded.Ids)
		if err != nil {
			return "", err
		}
		queryToAdd := `"id":{"$in":` + string(encodedArr) + `},`
		mainQuery += queryToAdd
	}

	isValid, err := isTimeFilterValid(&filterDecoded.TimeFilter)
	if err != nil {
		return "", err
	}

	if isValid {
		minimumEncoded, err := json.Marshal(filterDecoded.TimeFilter.Min)
		if err != nil {
			return "", err
		}
		maximumEncoded, err := json.Marshal(filterDecoded.TimeFilter.Max)
		if err != nil {
			return "", err
		}
		queryToAdd := `"timestamp":{"$gte":` + string(minimumEncoded) + `,` + `"$lte":` + string(maximumEncoded) + `},`
		mainQuery += queryToAdd
	}

	if len(mainQuery) != initQueryLen {
		mainQuery = strings.TrimSuffix(mainQuery, ",")
	}

	mainQuery += `}}`
	return mainQuery, nil
}

func isTimeFilterValid(filter *dtos.TimestampFilter) (bool, error) {
	if filter.Min.IsZero() && filter.Max.IsZero() {
		return false, nil
	}

	if filter.Min.IsZero() {
		return false, fmt.Errorf("mininum interval is invalid, while maximum is valid")
	}

	if filter.Max.IsZero() {
		return false, fmt.Errorf("maximum interval is invalid,while minimum is valid")
	}

	if filter.Min.After(filter.Max) {
		return false, fmt.Errorf("minimum interval should not be after the maximum")
	}

	if filter.Min.Equal(filter.Max) {
		return false, fmt.Errorf("intervals should not be equal")
	}

	return true, nil
}

func encodeArray(arr []string) ([]byte, error) {
	encoded, err := json.Marshal(arr)
	if err != nil {
		return nil, fmt.Errorf("error encoding string array %s", err)
	}

	return encoded, nil
}

func cleanFilter(filterDecoded *dtos.Filter) {
	if filterDecoded.Hashs != nil {
		clearAllStringFields(&filterDecoded.Hashs)
	}

	if filterDecoded.Ids != nil {
		clearAllStringFields(&filterDecoded.Ids)
	}

	if filterDecoded.InsertionTypes != nil {
		clearAllStringFields(&filterDecoded.InsertionTypes)
	}

	if filterDecoded.TypeForms != nil {
		clearAllStringFields(&filterDecoded.TypeForms)
	}
}

func clearAllStringFields(value *[]string) {
	for i := 0; i < len(*value); i++ {
		(*value)[i] = utils.RemoveStringSpaces((*value)[i])
	}
}

func (s *SmartContract) queryAllSetsWithPagination(
	context contractapi.TransactionContextInterface,
	query string,
	page int,
	size int,
) ([]*dtos.GetAllAssetsRequest, error) {
	allAssets := []*dtos.GetAllAssetsRequest{}
	bookmark := ""
	for i := 0; i <= page; i++ {
		isInCorrectPage := i == page
		canContinue, newBookMark, err := querySinglePage(context, query, int32(size), bookmark, &allAssets, isInCorrectPage)
		if err != nil {
			return nil, err
		}

		if !isThereANewPage(bookmark, newBookMark) {
			break
		}
		bookmark = newBookMark

		if !canContinue {
			break
		}
	}

	return allAssets, nil
}

func isThereANewPage(currentBookmark string, newBookmark string) bool {
	return currentBookmark != newBookmark
}

func querySinglePage(
	context contractapi.TransactionContextInterface,
	query string,
	size int32,
	bookmark string,
	getAllAssetRequestDto *[]*dtos.GetAllAssetsRequest,
	correctPage bool,
) (canIContinue bool, newBookmark string, err error) {
	iterator, responseMetadata, err := context.GetStub().GetQueryResultWithPagination(
		query,
		size,
		bookmark,
	)
	if err != nil {
		return false, bookmark, fmt.Errorf("error querying the ledger %s", err)
	}
	defer iterator.Close()

	if !correctPage {
		return true, responseMetadata.Bookmark, nil
	}

	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return false, bookmark, fmt.Errorf("error getting an item from the iterator %s", err)
		}

		asset := &dtos.GetAllAssetsRequest{}
		err = json.Unmarshal(queryResponse.Value, asset)
		if err != nil {
			return false, bookmark, fmt.Errorf("error decoding value from the ledger %s", err)
		}
		*getAllAssetRequestDto = append(*getAllAssetRequestDto, asset)
	}

	return responseMetadata.Bookmark != "", responseMetadata.Bookmark, nil
}

func validateDataGetAllAssets(pageString string, sizeString string) (int, int, error) {
	page, size, err := utils.ValidatePageAndSize(pageString, sizeString)
	if err != nil {
		return 0, 0, err
	}
	return page, size, err
}
