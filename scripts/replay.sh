#!/bin/sh -e

if [ $# -ne 0 ]
then
  echo "Usage: $0" >&2
  exit 1
fi

evm-bully replay -accountId evm.evm-bully.testnet -goerli evm.evm-bully.testnet
