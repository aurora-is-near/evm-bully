#!/bin/sh -e

if [ $# -ne 2 ]
then
  echo "Usage: $0 testnet_flag evm.wasm" >&2
  exit 1
fi

ACCOUNT=test.near
EVM=`openssl rand -hex 16`

export NEAR_ENV=local

near create-account $EVM.$ACCOUNT --master-account=$ACCOUNT --initial-balance=1000 --keyPath=$HOME/.near/local/validator_key.json
aurora install --chain 1313161556 --engine $EVM.$ACCOUNT --signer $EVM.$ACCOUNT --owner $EVM.$ACCOUNT $2

evm-bully -v replay -accountId $EVM.$ACCOUNT $1 -skip $EVM.$ACCOUNT
