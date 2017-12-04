package writebtc

import (
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
	"time"
)

/**
	Btc provides access to the Bitcoin blockchain

	init()     starts everything up.
	shutdown() shuts it all down.
**/
type Btc struct {
	ID       int
	Client   *rpcclient.Client
	sources  []string
	privates []*btcec.PrivateKey
	srcIndex int
	balance  int64
}

func (b *Btc) Shutdown() {
	// For this example gracefully shutdown the Client after 10 seconds.
	// Ordinarily when to shutdown the Client is highly application
	// specific.
	log.Println("Client shutdown in 2 seconds...")
	time.AfterFunc(time.Second*2, func() {
		log.Println("Going down...")
		b.Client.Shutdown()
	})
	defer log.Println("Shutdown done!")
	// Wait until the client either shuts down gracefully (or the user
	// terminates the process with Ctrl+C).
	b.Client.WaitForShutdown()
}

var counter int = 1

func (b *Btc) Init(sources []string, privates []*btcec.PrivateKey) (err error) {
	b.ID = counter
	counter++

	// Only override the handlers for notifications you care about.
	// Also note most of the handlers will only be called if you register
	// for notifications.  See the documentation of the rpcclient
	// NotificationHandlers type for more details about each handler.
	ntfnHandlers := rpcclient.NotificationHandlers{
		OnBlockConnected: func(hash *chainhash.Hash, height int32, t time.Time) {
			go b.newBlock(hash, height)
		},
	}

	// Copy the sources
	b.sources = make([]string, len(sources), len(sources))
	copy(b.sources, sources)

	// Copy the private keys
	b.privates = make([]*btcec.PrivateKey, len(privates), len(privates))
	copy(b.privates, privates)

	// Connect to local btcwallet RPC server using websockets.
	//certHomeDir := btcutil.AppDataDir("btcwallet", false)
	u, _ := user.Current()

	//certHomeDir := btcutil.AppDataDir(".", false)
	certHomeDir := u.HomeDir
	certs, err := ioutil.ReadFile(filepath.Join(certHomeDir, "rpc.cert"))
	if err != nil {
		log.Println(err)
		return
	}

	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:19001",
		Endpoint:     "ws",
		User:         "admin1",
		Pass:         "123",
		Certificates: certs,
		HTTPPostMode: true,
	}

	b.Client, err = rpcclient.New(connCfg, &ntfnHandlers)

	log.Println("in Init, Client: ", b.ID, b.Client)

	if err != nil || b.Client == nil {
		log.Println(err)
		log.Println("xx")
		if b.Client == nil {
			log.Println("client is nil")
		}
		return
	}

	b.srcIndex = 0

	return
}

/**
 * Record a hash into the Bitcoin Block Chain OP_RETURN
**/
func (b *Btc) RecordHash(hash []byte) (txhash *chainhash.Hash, err error) {

	txhash = nil

	if len(hash) != 32 {
		err = fmt.Errorf("len(hash) = %d \n Note that hash[] must be 32 bytes in length", len(hash))
		return
	}

	header := []byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0}
	hash = append(header, hash...)

	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_RETURN)
	builder.AddData(hash)
	opReturn, _ := builder.Script()
	//opReturn, _ := txscript.NullDataScript([]byte("the message"))
	disasm1, err := txscript.DisasmString(opReturn)

	// Create a public key script that pays to the address.
	//addr, err := btcutil.DecodeAddress(b.sources[0], ActiveNet.Params)
	//var changeS []byte
	//changeS, err = txscript.PayToAddrScript(addr)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	//disasm2, err2 := txscript.DisasmString(changeS)
	//if err2 != nil {
	//		log.Println(err2)
	//		return
	//	}
	log.Printf("Op_Return Hex:      %x\n", opReturn)
	log.Println("Op_Return:          ", disasm1)
	//log.Printf("Change Hex:         %x\n", changeS)
	//log.Println("Change Disassembly: ", disasm2)

	tx := wire.NewMsgTx(1)

	txOut := wire.NewTxOut(0, opReturn)
	tx.AddTxOut(txOut)
	//txOut = wire.NewTxOut(int64(1000), changeS)
	//tx.AddTxOut(txOut)

	//client.SendRawTransaction(tx, false)

	return
}

//
// newBlock hashes the current LastHash into the Bitcoin Blockchain 5 minutes after
// the previous block. (If a block is signed quicker than 5 minutes, then the second
// block is ignored.)
//
func (b *Btc) newBlock(hash *chainhash.Hash, height int32) {

	log.Printf("Block connected: %v (%d)", hash, height)

}
