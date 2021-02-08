## Server setup

Instructions for a Ubuntu 20.04 LTS development environment.

### Some tools

    sudo apt install -y build-essential lld git ncdu tree screen

### Set up NEAR node

See
https://docs.near.org/docs/develop/evm/evm-local-setup\#set-up-near-node

### Set up Node.js

    curl -sL https://deb.nodesource.com/setup_15.x -o nodesource_setup.sh
    sudo bash nodesource_setup.sh
    sudo apt install -y nodejs

    sudo npm install -g truffle
    sudo npm install -g near-cli

### Install Go

See https://golang.org/doc/install

### Install `geth`

See https://github.com/ethereum/go-ethereum

Fast sync testnets:

    geth --ropsten --verbosity=2 --vmodule core=3
    geth --rinkeby --verbosity=2 --vmodule core=3 --port 30304
    geth --goerli  --verbosity=2 --vmodule core=3 --port 30305
