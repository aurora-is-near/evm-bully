package replayer

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// ChainID queries the given endpoint for the chain ID.
func ChainID(ctx context.Context, endpoint string) (uint64, error) {
	c, err := rpc.DialContext(ctx, endpoint)
	if err != nil {
		return 0, err
	}
	ec := ethclient.NewClient(c)
	defer ec.Close()
	chainID, err := ec.ChainID(ctx)
	if err != nil {
		return 0, err
	}
	cID, err := bigIntToUint64(chainID)
	if err != nil {
		return 0, err
	}
	return cID, nil
}
