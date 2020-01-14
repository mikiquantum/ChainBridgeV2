package ethereum

import (
	"encoding/binary"
	"math/big"

	"github.com/ChainSafe/ChainBridgeV2/chains"
	"github.com/ChainSafe/ChainBridgeV2/common"
	msg "github.com/ChainSafe/ChainBridgeV2/message"
	"github.com/ChainSafe/log15"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

var _ chains.Writer = &Writer{}

type Writer struct {
	cfg   Config
	conn  *Connection
	nonce uint64
}

func NewWriter(conn *Connection, cfg *Config) *Writer {
	return &Writer{
		cfg:  *cfg,
		conn: conn,
	}
}

func (w *Writer) Start() error {
	log15.Debug("Starting ethereum writer...")
	log15.Warn("Writer.Start() not fully implemented")
	return nil
}

// ResolveMessage handles any given message based on type
// Note: We are currently panicking here, we should develop a better method for handling failures (possibly a queue?)
func (w *Writer) ResolveMessage(m msg.Message) {
	log15.Trace("Attempting to resolve message", "type", m.Type, "src", m.Source, "dst", m.Destination)
	var tx *ethtypes.Transaction
	var calldata []byte

	if m.Type == msg.DepositAssetType {
		log15.Info("Handling Deposit Asset message", "to", w.conn.cfg.receiver, "msgdata", m.Data)
		id := common.FunctionId("store(bytes32)")
		calldata = append(id, m.Data...)
	} else if m.Type == msg.GenericDepositType {
		log15.Info("Handling generic deposit message...", "to", w.conn.cfg.receiver, "msgdata", m.Data)
		currBlock, err := w.conn.LatestBlock()
		if err != nil {
			panic(err)
		}

		address := ethcommon.HexToAddress(w.conn.kp.Public().Address())

		nonce, err := w.conn.NonceAt(address, currBlock.Number())
		if err != nil {
			panic(err)
		}

		id := common.FunctionId("createDepositProposal(bytes32,uint256,uint256)")
		calldata := append(id, ethcrypto.Keccak256(m.Data)...)

		// add nonce to calldata and increment
		nonceBytes := make([]byte, 32)
		binary.BigEndian.PutUint64(nonceBytes[24:], w.nonce)
		calldata = append(calldata, nonceBytes...)
		w.nonce += 1

		// add origin chain to calldata
		chainIdBytes := make([]byte, 32)
		chainIdBytes[31] = uint8(m.Source)
		calldata = append(calldata, chainIdBytes...)

		tx = ethtypes.NewTransaction(
			nonce,
			w.conn.cfg.receiver,
			big.NewInt(0),  // TODO: value?
			1000000,        // TODO: gasLimit
			big.NewInt(10), // TODO: gasPrice
			calldata,
		)

		data, err := tx.MarshalJSON()
		if err != nil {
			panic(err)
		}

		err = w.conn.SubmitTx(data)
		if err != nil {
			panic(err)
		}

		nonce, err = w.conn.PendingNonceAt(address)
		if err != nil {
			panic(err)
		}

		// executeDeposit(uint _originChainId, uint _depositId, address _to, bytes memory _data)
		id = common.FunctionId("executeDeposit(uint256,uint256,address,bytes)")
		calldata = make([]byte, 32)
		calldata[31] = uint8(m.Source)

		// add nonce (depositId) to calldata
		calldata = append(calldata, nonceBytes...)

		// add destination chain to calldata
		// toBytes := make([]byte, 20)
		// calldata = append(calldata, toBytes...)

		// add hash data
		calldata = append(calldata, m.Data...)

		tx = ethtypes.NewTransaction(
			nonce,
			w.conn.cfg.receiver,
			big.NewInt(0),  // TODO: value?
			1000000,        // TODO: gasLimit
			big.NewInt(10), // TODO: gasPrice
			calldata,
		)

		data, err = tx.MarshalJSON()
		if err != nil {
			panic(err)
		}

		err = w.conn.SubmitTx(data)
		if err != nil {
			panic(err)
		}

		return
	} else {
		panic("not implemented")
	}

	currBlock, err := w.conn.LatestBlock()
	if err != nil {
		panic(err)
	}
	address := ethcommon.HexToAddress(w.conn.kp.Public().Address())

	nonce, err := w.conn.NonceAt(address, currBlock.Number())
	if err != nil {
		panic(err)
	}

	tx = ethtypes.NewTransaction(
		nonce,
		w.conn.cfg.receiver,
		// TODO: Make these configurable
		big.NewInt(0),
		6721975,
		big.NewInt(20000000000),
		calldata,
	)

	data, err := tx.MarshalJSON()
	if err != nil {
		panic(err)
	}

	err = w.conn.SubmitTx(data)
	if err != nil {
		panic(err)
	}
}

func (w *Writer) Stop() error {
	return nil
}
