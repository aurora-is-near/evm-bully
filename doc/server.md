## Server setup

Instructions for a local Ubuntu 20.04 LTS development environment.

### Some tools

    sudo apt install -y build-essential lld git ncdu tree screen

### Get Rust

    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

Follow instructions that appear after that command, then:

    rustup target add wasm32-unknown-unknown

**Note**: this custom target is not strictly needed for running a nearcore node, but is needed for building smart contracts in Rust.

### Set up NEAR node

Clone the [nearcore repository](https://github.com/near/nearcore) with:

    git clone https://github.com/near/nearcore.git

Navigate to the project root:

    cd nearcore

Build (this will take a while, feel free to move on to future steps while this is happening):

    cargo build -p neard --release --features protocol_feature_evm,nightly_protocol_features

When the build is complete, initialize with:

    ./target/release/neard --home=$HOME/.near/local init

Then run:

    ./target/release/neard --home=$HOME/.near/local run

**Note**: hit Ctrl + C to stop the local node. If you want to pick up where you left off, just use this final "run" command again. If you'd like to start from scratch, remove the folder:

    rm -rf ~/.near/local

and then use the "initialize" and "run" commands.

### Set up Node.js

    curl -sL https://deb.nodesource.com/setup_15.x -o nodesource_setup.sh
    sudo bash nodesource_setup.sh
    sudo apt install -y nodejs
    sudo npm install -g near-cli
    sudo npm install -g aurora-is-near/aurora-cli


### Install Go

See https://golang.org/doc/install

### Install `geth`

See https://github.com/ethereum/go-ethereum

Fast sync testnets:

    geth --ropsten --verbosity=2 --vmodule core=3
    geth --rinkeby --verbosity=2 --vmodule core=3 --port 30304
    geth --goerli  --verbosity=2 --vmodule core=3 --port 30305
