## Notes

1.  Change chain ID of `neard` to that of Ethereum testnet.
2.  Sync Ethereum testnet with `geth` (fast sync).
3.  Read transactions from DB directly (with go-leveldb?).
4.  Feed them to NEAR node via Web3 endpoint.

TODO: Genesis block.

-   header of block bodies:
    https://github.com/ethereum/go-ethereum/blob/master/core/rawdb/schema.go\#L81
    (use block hash to access it)
-   declaration of genesis block for goerli (and others):
    https://github.com/ethereum/go-ethereum/blob/master/core/genesis_alloc.go\#L27
-   Decode RLP like this:
    https://github.com/ethereum/go-ethereum/blob/master/core/genesis.go\#L379
