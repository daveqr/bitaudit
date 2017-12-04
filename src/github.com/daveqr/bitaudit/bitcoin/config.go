package bitcoin

import (
	"github.com/btcsuite/btcd/rpcclient"
	"io/ioutil"
)

var config = initConfig()

func initConfig() rpcclient.ConnConfig {
	certs, _ := ioutil.ReadFile("/tmp/rpc.cert")

	config := rpcclient.ConnConfig{
		Host:         "localhost:19001",
		//Endpoint:     "ws",
		User:         "admin1",
		Pass:         "123",
		Certificates: certs,
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	return config
}
