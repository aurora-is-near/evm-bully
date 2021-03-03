package replayer

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// BlockNumber queries the given endpoint for the current block number.
func BlockNumber(ctx context.Context, endpoint string) (uint64, error) {
	c, err := rpc.DialContext(ctx, endpoint)
	if err != nil {
		return 0, err
	}
	ec := ethclient.NewClient(c)
	defer ec.Close()
	return ec.BlockNumber(ctx)
}
