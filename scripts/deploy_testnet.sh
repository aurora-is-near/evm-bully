#!/bin/sh -e

if [ $# -ne 1 ]
then
  echo "Usage: $0 evm.wasm" >&2
  exit 1
fi

ACCOUNT=evm-bully.testnet

near delete evm.$ACCOUNT $ACCOUNT
near create-account evm.$ACCOUNT --master-account=$ACCOUNT --initial-balance=100
env NEAR_ENV=testnet aurora install --chain 1313161556 --engine evm.$ACCOUNT --signer evm.$ACCOUNT --owner evm.$ACCOUNT $1
