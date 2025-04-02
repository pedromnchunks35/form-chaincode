package chaincode

import (
	"encoding/json"
	"form-chaincode/dtos"
	"form-chaincode/mocks"
	"form-chaincode/utils"
	"github.com/golang/mock/gomock"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func Test_given_invalidId_thenException(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincodeStub := mocks.NewMockChaincodeStubInterface(controller)

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub)
	mockedChaincodeStub.EXPECT().GetState(utils.RemoveStringSpaces(normalId)).Return([]byte{}, nil)

	result, err := smartContract.GetHistoryAssetById(mockedTransaction, normalId)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "the asset doesn't exist")
	assert.Equal(t, "", result)
}

func Test_given_validIdAndOneItemHistory_thenReturnArrayLength1(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincodeStub := mocks.NewMockChaincodeStubInterface(controller)
	mockedHistoryIteratorMock := mocks.NewMockHistoryQueryIteratorInterface(controller)

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub).Times(1)
	mockedChaincodeStub.EXPECT().GetState(utils.RemoveStringSpaces(normalId)).Return([]byte{123}, nil)

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub).Times(1)
	mockedChaincodeStub.EXPECT().GetHistoryForKey(utils.RemoveStringSpaces(normalId)).Return(mockedHistoryIteratorMock, nil)

	mockedHistoryIteratorMock.EXPECT().HasNext().Return(true)

	value := dtos.GetAllAssetsRequest{
		Id:            utils.RemoveStringSpaces(normalId),
		TypeForm:      utils.RemoveStringSpaces(normalTypeForm),
		Description:   utils.RemoveStringSpaces(normalDescription),
		Timestamp:     utils.RemoveStringSpaces(normalTimestamp),
		InsertionType: utils.RemoveStringSpaces(normalInsertionType),
		Hash:          utils.RemoveStringSpaces(normalHash),
	}

	valueEncoded, err := json.Marshal(value)
	assert.Nil(t, err)

	item := &queryresult.KeyModification{
		TxId:      utils.RemoveStringSpaces(normalId),
		Timestamp: &timestamppb.Timestamp{},
		Value:     valueEncoded,
		IsDelete:  false,
	}
	mockedHistoryIteratorMock.EXPECT().Next().Return(item, nil).Times(1)
	mockedHistoryIteratorMock.EXPECT().HasNext().Return(false).Times(1)
	mockedHistoryIteratorMock.EXPECT().Close().Return(nil).Times(1)

	result, err := smartContract.GetHistoryAssetById(mockedTransaction, normalId)
	assert.Nil(t, err)
	assert.NotEqual(t, "", result)

	decodedResult := &[]*queryresult.KeyModification{}
	err = json.Unmarshal([]byte(result), decodedResult)
	assert.Nil(t, err)

	assert.Equal(t, len((*decodedResult)), 1)

	singleItem := (*decodedResult)[0]
	assert.Equal(t, singleItem, item)
}

func Test_given_validIdAndTwoItemHistory_thenReturnArrayLength2(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincodeStub := mocks.NewMockChaincodeStubInterface(controller)
	mockedHistoryIteratorMock := mocks.NewMockHistoryQueryIteratorInterface(controller)

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub).Times(1)
	mockedChaincodeStub.EXPECT().GetState(utils.RemoveStringSpaces(normalId)).Return([]byte{123}, nil)

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub).Times(1)
	mockedChaincodeStub.EXPECT().GetHistoryForKey(utils.RemoveStringSpaces(normalId)).Return(mockedHistoryIteratorMock, nil)

	mockedHistoryIteratorMock.EXPECT().HasNext().Return(true).Times(2)
	value := dtos.GetAllAssetsRequest{
		Id:            utils.RemoveStringSpaces(normalId),
		TypeForm:      utils.RemoveStringSpaces(normalTypeForm),
		Description:   utils.RemoveStringSpaces(normalDescription),
		Timestamp:     utils.RemoveStringSpaces(normalTimestamp),
		InsertionType: utils.RemoveStringSpaces(normalInsertionType),
		Hash:          utils.RemoveStringSpaces(normalHash),
	}

	valueEncoded, err := json.Marshal(value)
	assert.Nil(t, err)

	item := &queryresult.KeyModification{
		TxId:      utils.RemoveStringSpaces(normalId),
		Timestamp: &timestamppb.Timestamp{},
		Value:     valueEncoded,
		IsDelete:  false,
	}
	mockedHistoryIteratorMock.EXPECT().Next().Return(item, nil).Times(2)
	mockedHistoryIteratorMock.EXPECT().HasNext().Return(false).Times(1)
	mockedHistoryIteratorMock.EXPECT().Close().Return(nil).Times(1)

	result, err := smartContract.GetHistoryAssetById(mockedTransaction, normalId)
	assert.Nil(t, err)
	assert.NotEqual(t, "", result)

	decodedResult := &[]*queryresult.KeyModification{}
	err = json.Unmarshal([]byte(result), decodedResult)
	assert.Nil(t, err)

	assert.Equal(t, len((*decodedResult)), 2)

	singleItem := (*decodedResult)[0]
	assert.Equal(t, singleItem, item)

	otherItem := (*decodedResult)[1]
	assert.Equal(t, otherItem, item)
}
