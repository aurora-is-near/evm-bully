#!/bin/sh -e

if [ $# -ne 1 ]
then
  echo "Usage: $0 evm.wasm" >&2
  exit 1
fi

near delete evm.evm-bully.testnet evm-bully.testnet
near create-account evm.evm-bully.testnet --master-account=evm-bully.testnet --initial-balance=100
near deploy --account-id=evm.evm-bully.testnet --wasm-file=$1
