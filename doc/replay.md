## Replay an Ethereum testnet

See [server setup](server.md) for the generic server setup and [replay
transactions](replay-tx.md) on how to replay failing transactions.

It is assumed that the `evm-bully` is run out of the `evm-bully`
repository directory and that it is located "parallel" to the `nearcore`
and `aurora-engine` repositories:

    .
    ├── aurora-engine
    ├── evm-bully
    └── nearcore

### Install `evm-bully`

Run `make` in the `evm-bully` directory.

TODO
