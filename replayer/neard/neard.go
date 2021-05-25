package neard

import (
	"github.com/ethereum/go-ethereum/log"
)

type NEARDaemon struct {
}

func Setup(release bool) (*NEARDaemon, error) {
	log.Info("setup neard")
	return nil, nil
}

func (n *NEARDaemon) Stop() {
	log.Info("stop neard")
}
