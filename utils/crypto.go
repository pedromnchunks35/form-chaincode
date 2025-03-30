package utils

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"log"
	"os"
	"strconv"
)

func GetTLSProperties() shim.TLSProperties {
	tlsDisabledStr := GetEnvOrDefault("CHAINCODE_TLS_DISABLED", "true")
	key := GetEnvOrDefault("CHAINCODE_TLS_KEY", "")
	cert := GetEnvOrDefault("CHAINCODE_TLS_CERT", "")
	clientCACert := GetEnvOrDefault("CHAINCODE_CLIENT_CA_CERT", "")

	tlsDisabled := getBoolOrDefault(tlsDisabledStr, false)
	var keyBytes, certBytes, clientCACertBytes []byte
	var err error

	if !tlsDisabled {
		keyBytes, err = os.ReadFile(key)
		if err != nil {
			log.Panicf("error while reading the crypto file: %s", err)
		}
		certBytes, err = os.ReadFile(cert)
		if err != nil {
			log.Panicf("error while reading the crypto file: %s", err)
		}
	}

	if clientCACert != "" {
		clientCACertBytes, err = os.ReadFile(clientCACert)
		if err != nil {
			log.Panicf("error while reading the crypto file: %s", err)
		}
	}

	return shim.TLSProperties{
		Disabled:      tlsDisabled,
		Key:           keyBytes,
		Cert:          certBytes,
		ClientCACerts: clientCACertBytes,
	}
}

func GetEnvOrDefault(env, defaultVal string) string {
	value, ok := os.LookupEnv(env)
	if !ok || value == "" {
		value = defaultVal
	}
	return value
}

func getBoolOrDefault(value string, defaultVal bool) bool {
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultVal
	}
	return parsed
}
