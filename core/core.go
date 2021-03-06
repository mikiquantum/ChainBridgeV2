package core

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	msg "github.com/ChainSafe/ChainBridgeV2/message"
	"github.com/ChainSafe/ChainBridgeV2/router"

	"github.com/ChainSafe/ChainBridgeV2/keystore"
	"github.com/ChainSafe/log15"
)

type Core struct {
	registry map[msg.ChainId]*Chain
	route    *router.Router
	ks       *keystore.Keystore
}

func NewCore(ks *keystore.Keystore) *Core {
	return &Core{
		registry: make(map[msg.ChainId]*Chain),
		route:    router.NewRouter(),
		ks:       ks,
	}
}

func (c *Core) AddChain(chain *Chain) {
	err := c.route.Listen(chain.Id(), chain.GetWriter())
	if err != nil {
		log15.Error("Failed to add chain, will not be started", "id", chain.Id(), "err", err)
		return
	}
	c.registry[chain.Id()] = chain
	chain.listener.SetRouter(c.route)
}

// Start will call all registered chains' Start methods and block forever (or until signal is received)
func (c *Core) Start() {
	for _, ch := range c.registry {
		err := ch.Start()
		if err != nil {
			log15.Error(
				"failed to start chain",
				"chain", ch.Id(),
				"err", err,
			)
			return
		} else {
			log15.Info(fmt.Sprintf("Started %s chain", ch.Id()))
		}
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigc)

	// Block here and wait for a signal
	<-sigc
	log15.Info("Interrupt received, shutting down now.")
	for _, ch := range c.registry {
		err := ch.Stop()
		if err != nil {
			log15.Error(
				"failed to shutdown chain",
				"chain", ch.Id(),
				"err", err,
			)
		}
	}
}
