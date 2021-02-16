## Notes

[ ] 1.  Change chain ID of `neard` to that of Ethereum testnet.
[x] 2.  Sync Ethereum testnet with `geth` (fast sync).
[x] 3.  Read transactions from DB directly.
[ ] 4.  Feed them to NEAR node via Web3 endpoint.

TODO: Genesis block.

-   declaration of genesis block for goerli (and others):
    https://github.com/ethereum/go-ethereum/blob/master/core/genesis_alloc.go#L27
-   Decode RLP like this:
    https://github.com/ethereum/go-ethereum/blob/master/core/genesis.go#L379


### Ethereum testnet sizes

- goerli: ~13 GB
- rinkeby: ~66 GB
- ropsten: ~105 GB
