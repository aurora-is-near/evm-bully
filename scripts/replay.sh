#!/bin/sh -e

if [ $# -ne 0 ]
then
  echo "Usage: $0" >&2
  exit 1
fi

evm-bully -v replay -accountId evm-bully.testnet -goerli evm.evm-bully.testnet
