package chaincode

import (
	"encoding/json"
	"form-chaincode/dtos"
	"form-chaincode/mocks"
	"form-chaincode/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_givenInvalidId_whenPatchAsset_thenException(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)

	asset, err := smartContract.PatchAsset(mockedTransaction, "", emptyString)
	assert.Nil(t, asset)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "the id is not valid")
}

func Test_givenValidIdButAssetDoesNotExist_whenPatchAsset_thenException(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincode := mocks.NewMockChaincodeStubInterface(controller)

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincode)
	mockedChaincode.EXPECT().GetState(utils.RemoveStringSpaces(normalId)).Return(nil, nil)
	asset, err := smartContract.PatchAsset(mockedTransaction, "", normalId)
	assert.Nil(t, asset)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "it doesn't exist")
}

func Test_givenNilStructure_whenPatchAsset_thenException(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincode := mocks.NewMockChaincodeStubInterface(controller)

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincode)
	mockedChaincode.EXPECT().GetState(utils.RemoveStringSpaces(normalId)).Return([]byte{1, 0}, nil)

	asset, err := smartContract.PatchAsset(mockedTransaction, "", normalId)
	assert.Nil(t, asset)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "decoding the object")
}

func Test_givenNothingToPut_whenPatchAsset_thenException(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincode := mocks.NewMockChaincodeStubInterface(controller)

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincode)
	mockedChaincode.EXPECT().GetState(utils.RemoveStringSpaces(normalId)).Return([]byte{1, 0}, nil)

	assetToPut := &dtos.PutAssetRequest{}
	encoded, err := json.Marshal(assetToPut)
	assert.Nil(t, err)

	asset, err := smartContract.PatchAsset(mockedTransaction, string(encoded), normalId)
	assert.Nil(t, asset)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "nothing to change in the request")
}

func Test_givenSomethingToPut_whenPatchAsset_thenReturnAsset(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincode := mocks.NewMockChaincodeStubInterface(controller)

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincode).Times(4)

	assetToPut := &dtos.PutAssetRequest{
		Hash: utils.RemoveStringSpaces("something"),
	}
	encoded, err := json.Marshal(assetToPut)
	assert.Nil(t, err)

	givenAsset := &dtos.AssetRequest{
		TypeForm: "something2",
	}
	encodedAssetFromDb, err := json.Marshal(givenAsset)
	assert.Nil(t, err)

	mockedChaincode.EXPECT().GetState(utils.RemoveStringSpaces(normalId)).Return(encodedAssetFromDb, nil).Times(3)

	mockedChaincode.EXPECT().PutState(utils.RemoveStringSpaces(normalId), gomock.Any()).Times(1)

	asset, err := smartContract.PatchAsset(mockedTransaction, string(encoded), normalId)
	assert.Nil(t, err)
	assert.NotNil(t, asset)
	assert.Equal(t, asset.Hash, assetToPut.Hash)
	assert.Equal(t, asset.TypeForm, givenAsset.TypeForm)
}
