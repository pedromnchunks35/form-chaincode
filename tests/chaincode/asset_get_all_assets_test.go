package chaincode

import (
	"form-chaincode/mocks"
	"github.com/golang/mock/gomock"
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
