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
NEAR account and install the EVM contract to it (see
[`test_local.sh`](../scripts/test_local.sh)) or let the `evm-bully` do
what with the options `-setup` and `-contract`. We use the latter
approach by employing the wrapper script
[`test_local_setup.sh`](../scripts/test_local_setup.sh):

Example:

    ./scripts/test_local_setup.sh -goerli ../aurora-engine/release.wasm

### Options

-   Use `-autobreak` to automatically repeat with a break point after an
    error. Leads to a [replayable](replay-tx.md) problem `.tar.gz` file.
-   Use `-contract` to set the EVM contract file to deploy. Requires
    option `-setup`.
-   Use `-initial-balance` to set the number of tokens to transfer to
    newly created account. Requires option `-setup`.
-   Use `-keyPath` to set the path to master account key.
-   Use `-release` to run release version of neard (instead of debug
    version).
-   Use `-setup` to setup and run neard before replaying (auto-deploys
    contract). Requires option `-contract`. See [setup
    option](#setup-option) for details.
-   Use `-skip` to skip empty blocks during replay.

#### Testnet options

-   Use `-goerli` to use the Görli testnet.
-   Use `-rinkeby` to use the Rinkeby testnet.
-   Use `-ropsten` to use the Ropsten testnet.

### Setup option

Using the options `-setup` and `-contract` automatically starts `neard`
and deploys the EVM contract before replaying transactions. The
following steps are executed:

-   Switch to the `nearcore` directory, compile the current `HEAD`, and
    start `neard`.
-   Create a NEAR account, using the optional argument supplied to
    `evm-bully replay` as the account name. If no argument has been
    given the account name is automatically generated.
-   Installing the EVM contract supplied to `-contract` under the
    created contract by calling out to
    [`aurora-cli`](https://github.com/aurora-is-near/aurora-cli).

Afterwards the transactions from the supplied testnet are replayed until
an error occurs.
