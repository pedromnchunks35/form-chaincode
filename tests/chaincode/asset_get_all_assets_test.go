package chaincode

import (
	"encoding/json"
	"fmt"
	"form-chaincode/dtos"
	"form-chaincode/mocks"
	"form-chaincode/utils"
	"github.com/golang/mock/gomock"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GivenInvalidPageChar_whenGetAllAssets_thenException(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	assets, err := smartContract.GetAllAssets(mockedTransaction, "l", "10", []byte{1, 0})
	assert.Nil(t, assets)
	assert.NotNil(t, err)
}

func Test_GivenInvalidSizeChar_whenGetAllAssets_thenException(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	assets, err := smartContract.GetAllAssets(mockedTransaction, "0", "l", []byte{1, 0})
	assert.Nil(t, assets)
	assert.NotNil(t, err)
}

func Test_GivenInvalidPage_whenGetAllAssets_thenException(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	assets, err := smartContract.GetAllAssets(mockedTransaction, "-1", "10", []byte{1, 0})
	assert.Nil(t, assets)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "page and size are not consistent")
}

func Test_GivenInvalidSize_whenGetAllAssets_thenException(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	assets, err := smartContract.GetAllAssets(mockedTransaction, "0", "-1", []byte{1, 0})
	assert.Nil(t, assets)
	assert.NotNil(t, err)
}

func Test_GivenInvalidSize2_whenGetAllAssets_thenException(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	assets, err := smartContract.GetAllAssets(mockedTransaction, "0", "0", []byte{1, 0})
	assert.Nil(t, assets)
	assert.NotNil(t, err)
}

func Test_GivenNilFilter_whenGetAllAssets_thenException(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	assets, err := smartContract.GetAllAssets(mockedTransaction, "0", "10", nil)
	assert.Nil(t, assets)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "error decoding filter")
}

func Test_GivenEmptyFilterAndErrorIterating_whenGetAllAssets_thenReturnException(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincodeStub := mocks.NewMockChaincodeStubInterface(controller)
	filter := &dtos.Filter{}
	encodedFilter, err := json.Marshal(filter)
	assert.Nil(t, err)

	expectedQuery := `{"selector":{}}`

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub).Times(1)

	mockedChaincodeStub.EXPECT().GetQueryResultWithPagination(expectedQuery, int32(1), "").Return(nil, nil, fmt.Errorf("lol"))

	assets, err := smartContract.GetAllAssets(mockedTransaction, "0", "1", encodedFilter)
	assert.Nil(t, assets)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "error querying the ledger")
}

func Test_GivenEmptyFilterAndOneSizePage_whenGetAllAssets_thenReturnOneItem(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincodeStub := mocks.NewMockChaincodeStubInterface(controller)
	mockedIterator := mocks.NewMockStateQueryIteratorInterface(controller)
	filter := &dtos.Filter{}
	encodedFilter, err := json.Marshal(filter)
	assert.Nil(t, err)

	expectedQuery := `{"selector":{}}`

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub).Times(1)

	metadata := &peer.QueryResponseMetadata{
		FetchedRecordsCount: 1,
		Bookmark:            "",
	}
	mockedChaincodeStub.EXPECT().GetQueryResultWithPagination(expectedQuery, int32(1), "").Return(mockedIterator, metadata, nil)
	mockedIterator.EXPECT().HasNext().Return(true).Times(1)
	mockedIterator.EXPECT().HasNext().Return(false).Times(1)
	mockedIterator.EXPECT().Close().Return(nil)

	asset := &dtos.GetAllAssetsRequest{
		Id:            normalId,
		TypeForm:      normalTypeForm,
		Description:   normalDescription,
		Timestamp:     normalTimestamp,
		InsertionType: normalInsertionType,
		Hash:          normalHash,
	}
	encodedAsset, err := json.Marshal(asset)
	assert.Nil(t, err)
	queryResult := &queryresult.KV{
		Value: encodedAsset,
	}
	mockedIterator.EXPECT().Next().Return(queryResult, nil)

	assets, err := smartContract.GetAllAssets(mockedTransaction, "0", "1", encodedFilter)
	assert.Equal(t, len(assets), 1)

	singleAsset := assets[0]
	assert.Equal(t, asset.Id, singleAsset.Id)
	assert.Equal(t, asset.TypeForm, singleAsset.TypeForm)
	assert.Equal(t, asset.Description, singleAsset.Description)
	assert.Equal(t, asset.Timestamp, singleAsset.Timestamp)
	assert.Equal(t, asset.InsertionType, singleAsset.InsertionType)
	assert.Equal(t, asset.Hash, singleAsset.Hash)
}

func Test_GivenEmptyFilterAndFiveSizePage_whenGetAllAssets_thenReturnFiveItems(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincodeStub := mocks.NewMockChaincodeStubInterface(controller)
	mockedIterator := mocks.NewMockStateQueryIteratorInterface(controller)
	filter := &dtos.Filter{}
	encodedFilter, err := json.Marshal(filter)
	assert.Nil(t, err)

	expectedQuery := `{"selector":{}}`

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub).Times(1)

	metadata := &peer.QueryResponseMetadata{
		FetchedRecordsCount: 5,
		Bookmark:            "",
	}
	mockedChaincodeStub.EXPECT().GetQueryResultWithPagination(expectedQuery, int32(5), "").Return(mockedIterator, metadata, nil)
	mockedIterator.EXPECT().HasNext().Return(true).Times(5)
	mockedIterator.EXPECT().HasNext().Return(false).Times(1)
	mockedIterator.EXPECT().Close().Return(nil)

	asset := &dtos.GetAllAssetsRequest{
		Id:            normalId,
		TypeForm:      normalTypeForm,
		Description:   normalDescription,
		Timestamp:     normalTimestamp,
		InsertionType: normalInsertionType,
		Hash:          normalHash,
	}
	encodedAsset, err := json.Marshal(asset)
	assert.Nil(t, err)
	queryResult := &queryresult.KV{
		Value: encodedAsset,
	}
	mockedIterator.EXPECT().Next().Return(queryResult, nil).Times(5)

	assets, err := smartContract.GetAllAssets(mockedTransaction, "0", "5", encodedFilter)
	assert.Equal(t, len(assets), 5)
}

func Test_GivenEmptyFilterAndNextpageSizeOne_whenGetAllAssets_thenReturnOneItem(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincodeStub := mocks.NewMockChaincodeStubInterface(controller)
	mockedIterator := mocks.NewMockStateQueryIteratorInterface(controller)
	filter := &dtos.Filter{}
	encodedFilter, err := json.Marshal(filter)
	assert.Nil(t, err)

	expectedQuery := `{"selector":{}}`

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub).Times(2)

	metadata := &peer.QueryResponseMetadata{
		FetchedRecordsCount: 1,
		Bookmark:            "firstBookmark",
	}
	mockedChaincodeStub.EXPECT().GetQueryResultWithPagination(expectedQuery, int32(1), "").Return(mockedIterator, metadata, nil).Times(1)

	mockedChaincodeStub.EXPECT().GetQueryResultWithPagination(expectedQuery, int32(1), metadata.Bookmark).Return(mockedIterator, metadata, nil).Times(1)

	mockedIterator.EXPECT().HasNext().Return(true).Times(1)
	mockedIterator.EXPECT().HasNext().Return(false).Times(1)

	asset := &dtos.GetAllAssetsRequest{
		Id:            normalId,
		TypeForm:      normalTypeForm,
		Description:   normalDescription,
		Timestamp:     normalTimestamp,
		InsertionType: normalInsertionType,
		Hash:          normalHash,
	}
	encodedAsset, err := json.Marshal(asset)
	assert.Nil(t, err)
	queryResult := &queryresult.KV{
		Value: encodedAsset,
	}
	mockedIterator.EXPECT().Next().Return(queryResult, nil).Times(1)

	mockedIterator.EXPECT().Close().Return(nil).Times(2)
	assets, err := smartContract.GetAllAssets(mockedTransaction, "1", "1", encodedFilter)
	assert.Nil(t, err)
	assert.NotNil(t, assets)
	assert.Equal(t, len(assets), 1)
}

func Test_GivenEmptyFilterAndHasNextFalseSizeOne_whenGetAllAssets_thenReturnException(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincodeStub := mocks.NewMockChaincodeStubInterface(controller)
	mockedIterator := mocks.NewMockStateQueryIteratorInterface(controller)
	filter := &dtos.Filter{}
	encodedFilter, err := json.Marshal(filter)
	assert.Nil(t, err)

	expectedQuery := `{"selector":{}}`

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub).Times(1)

	metadata := &peer.QueryResponseMetadata{
		FetchedRecordsCount: 1,
		Bookmark:            "",
	}
	mockedChaincodeStub.EXPECT().GetQueryResultWithPagination(expectedQuery, int32(1), "").Return(mockedIterator, metadata, nil).Times(1)
	mockedIterator.EXPECT().Close().Return(nil).Times(1)

	assets, err := smartContract.GetAllAssets(mockedTransaction, "1", "1", encodedFilter)
	assert.Nil(t, err)
	assert.NotNil(t, assets)
	assert.Equal(t, len(assets), 0)
}

func Test_GivenHashFilterAndNoItems_whenGetAllAssets_thenArgsOk(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincodeStub := mocks.NewMockChaincodeStubInterface(controller)
	mockedIterator := mocks.NewMockStateQueryIteratorInterface(controller)
	filter := &dtos.Filter{
		Hash: normalHash,
	}
	encodedFilter, err := json.Marshal(filter)
	assert.Nil(t, err)

	expectedQuery := `{"selector":{"hash":"` + normalHash + `"}}`

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub).Times(1)

	metadata := &peer.QueryResponseMetadata{
		FetchedRecordsCount: 1,
		Bookmark:            "",
	}
	mockedChaincodeStub.EXPECT().GetQueryResultWithPagination(expectedQuery, int32(1), "").Return(mockedIterator, metadata, nil).Times(1)
	mockedIterator.EXPECT().Close().Return(nil).Times(1)

	assets, err := smartContract.GetAllAssets(mockedTransaction, "1", "1", encodedFilter)
	assert.Nil(t, err)
	assert.NotNil(t, assets)
	assert.Equal(t, len(assets), 0)
}

func Test_GivenHashAndIdsFilterAndNoItems_whenGetAllAssets_thenArgsOk(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincodeStub := mocks.NewMockChaincodeStubInterface(controller)
	mockedIterator := mocks.NewMockStateQueryIteratorInterface(controller)

	ids := make([]string, 0)
	ids = append(ids, normalId)

	filter := &dtos.Filter{
		Hash: normalHash,
		Ids:  ids,
	}
	encodedFilter, err := json.Marshal(filter)
	assert.Nil(t, err)

	expectedQuery := `{"selector":{"hash":"` + normalHash + `"` + `,"id":{"$in":["` + utils.RemoveStringSpaces(normalId) + `"]` + `}}}`

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub).Times(1)

	metadata := &peer.QueryResponseMetadata{
		FetchedRecordsCount: 1,
		Bookmark:            "",
	}
	mockedChaincodeStub.EXPECT().GetQueryResultWithPagination(expectedQuery, int32(1), "").Return(mockedIterator, metadata, nil).Times(1)
	mockedIterator.EXPECT().Close().Return(nil).Times(1)

	assets, err := smartContract.GetAllAssets(mockedTransaction, "1", "1", encodedFilter)
	assert.Nil(t, err)
	assert.NotNil(t, assets)
	assert.Equal(t, len(assets), 0)
}

func Test_GivenHashAndIdsAndTypeFormsFilterAndNoItems_whenGetAllAssets_thenArgsOk(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincodeStub := mocks.NewMockChaincodeStubInterface(controller)
	mockedIterator := mocks.NewMockStateQueryIteratorInterface(controller)

	ids := make([]string, 0)
	ids = append(ids, normalId)

	typeForms := make([]string, 0)
	typeForms = append(typeForms, normalTypeForm)

	filter := &dtos.Filter{
		Hash:      normalHash,
		Ids:       ids,
		TypeForms: typeForms,
	}
	encodedFilter, err := json.Marshal(filter)
	assert.Nil(t, err)

	expectedQuery := `{"selector":{"hash":"` + normalHash + `"` + `,"type_form":{"$in":["` + utils.RemoveStringSpaces(normalTypeForm) + `"]}` + `,"id":{"$in":["` + utils.RemoveStringSpaces(normalId) + `"]` + `}}}`

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub).Times(1)

	metadata := &peer.QueryResponseMetadata{
		FetchedRecordsCount: 1,
		Bookmark:            "",
	}
	mockedChaincodeStub.EXPECT().GetQueryResultWithPagination(expectedQuery, int32(1), "").Return(mockedIterator, metadata, nil).Times(1)
	mockedIterator.EXPECT().Close().Return(nil).Times(1)

	assets, err := smartContract.GetAllAssets(mockedTransaction, "1", "1", encodedFilter)
	assert.Nil(t, err)
	assert.NotNil(t, assets)
	assert.Equal(t, len(assets), 0)
}

func Test_GivenHashAndIdsAndTypeFormsAndInsertionTypesFilterAndNoItems_whenGetAllAssets_thenArgsOk(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincodeStub := mocks.NewMockChaincodeStubInterface(controller)
	mockedIterator := mocks.NewMockStateQueryIteratorInterface(controller)

	ids := make([]string, 0)
	ids = append(ids, normalId)

	typeForms := make([]string, 0)
	typeForms = append(typeForms, normalTypeForm)

	insertionTypes := make([]string, 0)
	insertionTypes = append(insertionTypes, normalInsertionType)

	filter := &dtos.Filter{
		Hash:           normalHash,
		Ids:            ids,
		TypeForms:      typeForms,
		InsertionTypes: insertionTypes,
	}
	encodedFilter, err := json.Marshal(filter)
	assert.Nil(t, err)

	expectedQuery := `{"selector":{"hash":"` + normalHash + `"` + `,"type_form":{"$in":["` + utils.RemoveStringSpaces(normalTypeForm) + `"]}` + `,"insertion_type":{"$in":["` + utils.RemoveStringSpaces(normalInsertionType) + `"]}` + `,"id":{"$in":["` + utils.RemoveStringSpaces(normalId) + `"]` + `}}}`

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub).Times(1)

	metadata := &peer.QueryResponseMetadata{
		FetchedRecordsCount: 1,
		Bookmark:            "",
	}
	mockedChaincodeStub.EXPECT().GetQueryResultWithPagination(expectedQuery, int32(1), "").Return(mockedIterator, metadata, nil).Times(1)
	mockedIterator.EXPECT().Close().Return(nil).Times(1)

	assets, err := smartContract.GetAllAssets(mockedTransaction, "1", "1", encodedFilter)
	assert.Nil(t, err)
	assert.NotNil(t, assets)
	assert.Equal(t, len(assets), 0)
}
