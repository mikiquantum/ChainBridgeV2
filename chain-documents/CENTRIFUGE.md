# Running a Centrifuge -> Ethereum NFT environemnt

### Requirements
1. Rust nightly
2. Rust Wasm toolchain 
3. Node >10.16.0

### Prerequisites
Ensure balance in from accounts defined in config.toml under defined keystore (`./keys` by default)

### Start
1. `make build`
2. `make start_cent`    # Note this will take some time, you can run step 3 in parallel
3. `make start_eth`     # Starts a ganache instance
4. `make deploy_eth`    # Deploys smart contracts
5. `make run`           # Run bridge
6. `make cent_asset_tx` # Begins the transfer process
