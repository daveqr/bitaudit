package server

import (
	"github.com/btcsuite/btcd/rpcclient"
	"io/ioutil"
)

func InitConfig() rpcclient.ConnConfig {
	certs, _ := ioutil.ReadFile("rpc.cert")

	connCfg := rpcclient.ConnConfig{
		Host:         "localhost:19001",
		Endpoint:     "ws",
		User:         "admin1",
		Pass:         "123",
		Certificates: certs,
		HTTPPostMode: true,
	}

	return connCfg
}
