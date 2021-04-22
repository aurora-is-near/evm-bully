#!/bin/sh -e

if [ $# -ne 0 ]
then
  echo "Usage: $0" >&2
  exit 1
fi

ACCOUNT=evm-bully.testnet

evm-bully -v replay -accountId evm.$ACCOUNT -goerli evm.$ACCOUNT
