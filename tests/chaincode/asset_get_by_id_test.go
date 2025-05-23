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
	"time"
)

var emptyString = " "
var normalId = "some _id"
var normalTypeForm = "some_type_form"
var normalDescription = "some_description"
var normalTimestamp = time.Now()
var normalInsertionType = "some_insertion_type"
var normalHash = "some_hash"

func Test_given_invalid_id_string_when_GetAssetById_thenReturnException(t *testing.T) {
	controller := gomock.NewController(t)
	mockedStub := mocks.NewMockTransactionContextInterface(controller)

	asset, err := smartContract.GetAssetById(mockedStub, emptyString)
	assert.Equal(t, "", asset)
	assert.NotNil(t, err)
	assert.Equal(t, err, fmt.Errorf("the id is not valid"))
}

func Test_given_invalid_id_whenGetAssetById_thenReturnException(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincodeStub := mocks.NewMockChaincodeStubInterface(controller)

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub)
	mockedChaincodeStub.EXPECT().GetState(utils.RemoveStringSpaces(normalId)).Return(nil, nil)

	asset, err := smartContract.GetAssetById(mockedTransaction, normalId)
	assert.Equal(t, "", asset)
	assert.NotNil(t, err)
	assert.Equal(t, err, fmt.Errorf("the asset doesn't exist"))
}

func Test_given_valid_id_whenGetAssetById_thenReturnTrueObject(t *testing.T) {
	controller := gomock.NewController(t)
	mockedTransaction := mocks.NewMockTransactionContextInterface(controller)
	mockedChaincodeStub := mocks.NewMockChaincodeStubInterface(controller)

	mockedTransaction.EXPECT().GetStub().Return(mockedChaincodeStub).Times(2)

	asset := &dtos.AssetRequest{
		Id:            utils.RemoveStringSpaces(normalId),
		TypeForm:      normalTypeForm,
		Description:   normalDescription,
		Timestamp:     normalTimestamp,
		InsertionType: normalInsertionType,
		Hash:          normalHash,
	}
	encodedAsset, err := json.Marshal(asset)
	assert.Nil(t, err)

	mockedChaincodeStub.EXPECT().GetState(utils.RemoveStringSpaces(normalId)).Return(encodedAsset, nil).Times(2)

	resultString, err := smartContract.GetAssetById(mockedTransaction, normalId)
	assert.Nil(t, err)

	result := &dtos.AssetRequest{}
	err = json.Unmarshal([]byte(resultString), result)
	assert.Nil(t, err)

	assert.NotNil(t, asset)
	assert.Equal(t, asset.Id, result.Id)
	assert.Equal(t, asset.Hash, result.Hash)
	assert.Equal(t, asset.Timestamp.Equal(result.Timestamp), true)
	assert.Equal(t, asset.Description, result.Description)
	assert.Equal(t, asset.InsertionType, result.InsertionType)
	assert.Equal(t, asset.TypeForm, result.TypeForm)
}
