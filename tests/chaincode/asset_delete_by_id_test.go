package chaincode

import (
	"fmt"
	"form-chaincode/mocks"
	"form-chaincode/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var invalidIdDelete = "   "

func Test_givenInvalidId_when_DeleteAssetById_thenReturnException(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)

	result, err := smartContract.DeleteAssetById(mockedTransaction, invalidIdDelete)
	assert.NotNil(t, result)
	assert.Equal(t, result, false)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "id is not valid")
}

func Test_givenIdForNonExistentAsset_whenDeleteAssetById_thenReturnException(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincodeStub := mocks.NewMockChaincodeStubInterface(controller)

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub)
	mockedChaincodeStub.EXPECT().GetState(utils.RemoveStringSpaces(normalId)).Return(nil, nil)

	result, err := smartContract.DeleteAssetById(mockedTransaction, normalId)
	assert.NotNil(t, result)
	assert.Equal(t, result, false)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "asset does't exist")
}

func Test_givenValidId_whenDeleteAssetById_thenSuccess(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincodeStub := mocks.NewMockChaincodeStubInterface(controller)

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub).Times(2)
	mockedChaincodeStub.EXPECT().GetState(utils.RemoveStringSpaces(normalId)).Return([]byte{1, 0, 1}, nil)

	mockedChaincodeStub.EXPECT().DelState(utils.RemoveStringSpaces(normalId)).Return(nil)
	result, err := smartContract.DeleteAssetById(mockedTransaction, normalId)
	assert.NotNil(t, result)
	assert.Equal(t, result, true)
	assert.Nil(t, err)
}

func Test_givenLedgerError_whenDeleteAssetById_thenError(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincodeStub := mocks.NewMockChaincodeStubInterface(controller)

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub).Times(2)
	mockedChaincodeStub.EXPECT().GetState(utils.RemoveStringSpaces(normalId)).Return([]byte{1, 0, 1}, nil)

	mockedChaincodeStub.EXPECT().DelState(utils.RemoveStringSpaces(normalId)).Return(fmt.Errorf("SOME EXCEPTION"))
	result, err := smartContract.DeleteAssetById(mockedTransaction, normalId)
	assert.NotNil(t, result)
	assert.Equal(t, result, false)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "error deleting state from the ledger")
}
