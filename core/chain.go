package core

import (
	"fmt"

	"github.com/ChainSafe/ChainBridgeV2/chains"
	"github.com/ChainSafe/ChainBridgeV2/keystore"
	msg "github.com/ChainSafe/ChainBridgeV2/message"
	log "github.com/ChainSafe/log15"
)

type Chain struct {
	cfg      *ChainConfig
	conn     chains.Connection
	listener chains.Listener
	writer   chains.Writer
}

type ChainConfig struct {
	Id            msg.ChainId // ChainID
	Endpoint      string      // url for rpc endpoint
	Receiver      string      // bridge address to call
	Emitter       string      // bridge address where events occur
	From          string
	Keystore      *keystore.Keystore
	Subscriptions []string `toml:"subscriptions"`
}

func NewChain(cfg *ChainConfig) *Chain {
	return &Chain{
		cfg: cfg,
	}
}

func (c *Chain) SetConnection(conn chains.Connection) {
	c.conn = conn
}

func (c *Chain) SetListener(listener chains.Listener) {
	c.listener = listener
}

func (c *Chain) SetWriter(Writer chains.Writer) {
	c.writer = Writer
}

func (c *Chain) GetWriter() chains.Writer {
	return c.writer
}

func (c *Chain) Start() error {
	if c.conn == nil {
		return fmt.Errorf("no connection specified")
	}
	if c.listener == nil {
		return fmt.Errorf("no listener specified")
	}
	if c.writer == nil {
		return fmt.Errorf("no Writer specified")
	}

	err := c.conn.Connect()
	if err != nil {
		return err
	}

	err = c.listener.Start()
	if err != nil {
		return err
	}

	err = c.writer.Start()
	if err != nil {
		return err
	}

	log.Debug("Successfully started chain")
	return nil
}

func (c *Chain) Id() msg.ChainId {
	return c.cfg.Id
}

func (c *Chain) Connection() chains.Connection {
	return c.conn
}

func (c *Chain) Stop() error {
	err := c.listener.Stop()
	if err != nil {
		return err
	}

	err = c.writer.Stop()
	if err != nil {
		return err
	}

	c.conn.Close()

	return nil
}
