package chaincode

import (
	"encoding/json"
	"fmt"
	"form-chaincode/dtos"
	"form-chaincode/mocks"
	"form-chaincode/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var normalIdCreation = "some _id"
var normalTypeFormCreation = "some_ type_form"
var normalDescriptionCreation = "some_ description"
var normalTimestampCreation = "some _tim estamp"
var normalInsertionTypeCreation = "s o me_insertion_type"
var normalHashCreation = "som e _has h"

func Test_givenNilAsset_whenCreateAsset_thenReturnError(t *testing.T) {
	controller := gomock.NewController(t)
	mockedStub := mocks.NewMockTransactionContextInterface(controller)

	result, err := smartContract.CreateAsset(mockedStub, "")
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "decoding the given value results in")
}

func Test_givenCompleteObjectWithEmtpyStrings_whenCreateAsset_thenReturnError(t *testing.T) {
	controller := gomock.NewController(t)
	mockedStub := mocks.NewMockTransactionContextInterface(controller)

	request := &dtos.PostAssetRequest{
		Id:            normalIdCreation,
		TypeForm:      emptyString,
		Description:   normalDescriptionCreation,
		Timestamp:     normalTimestampCreation,
		InsertionType: normalInsertionTypeCreation,
		Hash:          normalHashCreation,
	}
	encodedData, err := json.Marshal(request)
	assert.Nil(t, err)

	result, err := smartContract.CreateAsset(mockedStub, string(encodedData))
	assert.NotNil(t, err.Error(), "some fields are not valid")
	assert.Nil(t, result)
}

func Test_givenAlreadyExistentObject_whenCreateAsset_thenReturnError(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincodeStub := mocks.NewMockChaincodeStubInterface(controller)

	request := &dtos.PostAssetRequest{
		Id:            normalIdCreation,
		TypeForm:      normalTypeFormCreation,
		Description:   normalDescriptionCreation,
		Timestamp:     normalTimestampCreation,
		InsertionType: normalInsertionTypeCreation,
		Hash:          normalHashCreation,
	}
	encodedData, err := json.Marshal(request)
	assert.Nil(t, err)

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub)
	mockedChaincodeStub.EXPECT().GetState(utils.RemoveStringSpaces(normalIdCreation)).Return([]byte{0, 1, 0}, nil)

	result, err := smartContract.CreateAsset(mockedTransaction, string(encodedData))
	assert.NotNil(t, err.Error(), "already exists")
	assert.Nil(t, result)
}

func Test_givenCompleteValidObject_whenCreateAsset_thenReturnSameObject(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincodeStub := mocks.NewMockChaincodeStubInterface(controller)

	request := &dtos.PostAssetRequest{
		Id:            normalIdCreation,
		TypeForm:      normalTypeFormCreation,
		Description:   normalDescriptionCreation,
		Timestamp:     normalTimestampCreation,
		InsertionType: normalInsertionTypeCreation,
		Hash:          normalHashCreation,
	}
	encodedData, err := json.Marshal(request)
	assert.Nil(t, err)

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub).Times(2)
	mockedChaincodeStub.EXPECT().GetState(utils.RemoveStringSpaces(normalIdCreation)).Return(nil, nil)

	cleanRequest := &dtos.PostAssetRequest{
		Id:            utils.RemoveStringSpaces(normalIdCreation),
		TypeForm:      utils.RemoveStringSpaces(normalTypeFormCreation),
		Description:   normalDescriptionCreation,
		Timestamp:     utils.RemoveStringSpaces(normalTimestampCreation),
		InsertionType: utils.RemoveStringSpaces(normalInsertionTypeCreation),
		Hash:          utils.RemoveStringSpaces(normalHashCreation),
	}
	cleanEncodedData, err := json.Marshal(cleanRequest)
	assert.Nil(t, err)

	mockedChaincodeStub.EXPECT().PutState(utils.RemoveStringSpaces(normalIdCreation), cleanEncodedData).Return(nil)
	result, err := smartContract.CreateAsset(mockedTransaction, string(encodedData))
	assert.Nil(t, err)

	assert.Equal(t, result.Id, cleanRequest.Id)
	assert.Equal(t, result.TypeForm, cleanRequest.TypeForm)
	assert.Equal(t, result.Description, cleanRequest.Description)
	assert.Equal(t, result.Timestamp, cleanRequest.Timestamp)
	assert.Equal(t, result.InsertionType, cleanRequest.InsertionType)
	assert.Equal(t, result.Hash, cleanRequest.Hash)
}

func Test_givenExceptionOnPut_whenCreateAsset_thenReturnException(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincodeStub := mocks.NewMockChaincodeStubInterface(controller)

	request := &dtos.PostAssetRequest{
		Id:            normalIdCreation,
		TypeForm:      normalTypeFormCreation,
		Description:   normalDescriptionCreation,
		Timestamp:     normalTimestampCreation,
		InsertionType: normalInsertionTypeCreation,
		Hash:          normalHashCreation,
	}
	encodedData, err := json.Marshal(request)
	assert.Nil(t, err)

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub).Times(2)
	mockedChaincodeStub.EXPECT().GetState(utils.RemoveStringSpaces(normalIdCreation)).Return(nil, nil)

	cleanRequest := &dtos.PostAssetRequest{
		Id:            utils.RemoveStringSpaces(normalIdCreation),
		TypeForm:      utils.RemoveStringSpaces(normalTypeFormCreation),
		Description:   normalDescriptionCreation,
		Timestamp:     utils.RemoveStringSpaces(normalTimestampCreation),
		InsertionType: utils.RemoveStringSpaces(normalInsertionTypeCreation),
		Hash:          utils.RemoveStringSpaces(normalHashCreation),
	}
	cleanEncodedData, err := json.Marshal(cleanRequest)
	assert.Nil(t, err)

	mockedChaincodeStub.EXPECT().PutState(utils.RemoveStringSpaces(normalIdCreation), cleanEncodedData).Return(
		fmt.Errorf("some exception"),
	)
	result, err := smartContract.CreateAsset(mockedTransaction, string(encodedData))
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "inserting cleaned object")
}
