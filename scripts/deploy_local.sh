#!/bin/sh -e

if [ $# -ne 1 ]
then
  echo "Usage: $0 evm.wasm" >&2
  exit 1
fi

ACCOUNT=test.near

export NEAR_ENV=local

near delete evm.$ACCOUNT $ACCOUNT
near create-account evm.$ACCOUNT --master-account=$ACCOUNT --initial-balance=1000000 --keyPath=$HOME/.near/local/validator_key.json
aurora install --chain 1313161556 --engine evm.$ACCOUNT --signer evm.$ACCOUNT --owner evm.$ACCOUNT $1
