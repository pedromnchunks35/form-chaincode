# INFO
- This is a chaincode project for inserting different types of forms
and with different types of insertion

# Mockgen
- Generate transaction context from hlf using mock gen
```
mockgen -destination=mocks/mock_transaction_context.go -package=mocks 
github.com/hyperledger/fabric-contract-api-go/contractapi TransactionContextInterface
```
- Generate chaincode stub from hlf using mock gen
```
mockgen -destination=mocks/mock_chaincode_stub.go 
-package=mocks github.com/hyperledger/fabric-chaincode-go/shim ChaincodeStubInterface
```
- Generate chaincode stub from hlf using mock gen
```
mockgen -destination=mocks/mock_iterator.go -package=mocks 
github.com/hyperledger/fabric-chaincode-go/shim StateQueryIteratorInterface
```

# Generate image
- We create a docker file
