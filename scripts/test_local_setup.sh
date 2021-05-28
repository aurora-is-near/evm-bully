#!/bin/sh -e

if [ $# -ne 2 ]
then
  echo "Usage: $0 testnet_flag evm.wasm" >&2
  exit 1
fi

env NEAR_ENV=local evm-bully -v replay \
                   -initial-balance 1000 \
                   -keyPath $HOME/.near/local/validator_key.json \
                   -setup -skip $1 -contract $2
