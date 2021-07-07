## Replay an Ethereum testnet

See [server setup](server.md) for the generic server setup and [replay
transactions](replay-tx.md) on how to replay failing transactions.

It is assumed that the `evm-bully` is run out of the `evm-bully`
repository directory and that it is located "parallel" to the
[`nearcore`](https://github.com/near/nearcore/) and
[`aurora-engine`](https://github.com/aurora-is-near/aurora-engine)
repositories:

    .
    ├── aurora-engine
    ├── evm-bully
    └── nearcore

### Install `evm-bully`

Run `make` in the `evm-bully` directory.

### Compile Aurora Engine

Run `make evm-bully=yes` in `aurora-engine` directory.

### Replay transactions

Replaying transactions requires that the corresponding Ethereum testnet
is synched, see [synching testnets](server.md#synching-testnets).

In order to run `evm-bully replay` we either need to manually create a
NEAR account and install the EVM contract to it or let the `evm-bully`
do what with the options `-setup` and `-contract`. We use the latter
approach by employing the wrapper script
[`test_local_setup.sh`](../scripts/test_local_setup.sh):

Example:

    scripts/test_local_setup -goerli ../aurora-engine/release.wasm

### Options

TODO
