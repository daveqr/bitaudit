package writebtc

import (
	"github.com/btcsuite/btcd/chaincfg"
)

var ActiveNet = testNet3Params

type params struct {
	*chaincfg.Params
	connect string
	port    string
}

var mainNetParams = params{
	Params:  &chaincfg.MainNetParams,
	connect: "localhost:19001",
	port:    "19001",
}

var testNet3Params = params{
	Params:  &chaincfg.TestNet3Params,
	connect: "localhost:19001",
	port:    "19001",
}

var simNetParams = params{
	Params:  &chaincfg.SimNetParams,
	connect: "localhost:19001",
	port:    "19001",
}
