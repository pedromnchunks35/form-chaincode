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
- Generate state query iterator mock
```
mockgen -destination=mocks/mock_iterator.go -package=mocks 
github.com/hyperledger/fabric-chaincode-go/shim StateQueryIteratorInterface
```
- Generate history iterator mock
```
mockgen -destination=mocks/mock_history_iterator.go -package=mocks 
github.com/hyperledger/fabric-chaincode-go/shim HistoryQueryIteratorInterface
```

# Generate image
- We create a docker file
```
docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t pedrosilvamnchunks/form-chaincode:latest --push .
```

# Timestamp format
- Timestamp format will be the one from  `ISO 8601` which is the same as RFC3339
- E.g: "2025-04-05T12:30:45Z"
