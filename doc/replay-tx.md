## Replay failing transactions

It is assumed that the `evm-bully` is run out of the `evm-bully` repository
directory and that it is located "parallel" to the `nearcore` and
`aurora-engine` repositories:

```
.
├── aurora-engine
├── evm-bully
└── nearcore
```

### Install `evm-bully`

Run `make` in the `evm-bully` directory.

### Extract problem `.tar.gz` file

Example:

    tar xvzf rinkeby-block-55-tx-0.tar.gz

### Reproduce problem

Example:

    evm-bully replay-tx rinkeby-block-55-tx-0

This automatically starts the debug version of `neard` in `../nearcore`.

### Options

- Use `-release` to run release version of `neard`.
- Use `-build` to build "nearcore" version of `neard` given in `breakpoint.json` before running it.
- Use `-contract` to update the EVM contract in `neard` state before replaying transaction (for debugging).

`-contract` requires that given the contract in `../aurora-engine/` has been
build with `make evm-bully=yes` and that
[aurora-cli](https://github.com/aurora-is-near/aurora-cli) is installed in
`$PATH`.
