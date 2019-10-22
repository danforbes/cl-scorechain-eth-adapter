# Chainlink Scorechain Ethereum Adapter

Chainlink adapter for the [Scorechain Ethereum API](https://ethereum.scorechain.com/api_doc/).

## Environment Variables

* `SC_ETH_TOKEN` - token for the Scorechain Ethereum API
* `SC_ETH_URL` - the URL for the Scorechain Ethereum API

## Usage

The Scorechain Ethereum Adapter exposes the following capabilities:

### [Get Status](https://ethereum.scorechain.com/api_doc/#/blockchain%20data/getSystemStatus)

Get the status of the Scorechain Ethereum API.

#### Params

* `endpoint` - `get-status`

### [Get Transaction](https://ethereum.scorechain.com/api_doc/#/blockchain%20data/getTransaction)

Get a transaction from the Scorechain Ethereum API.

#### Params

* `endpoint` - `get-trx`
* `hash` - hash of the desired transaction `string`

### [Get Scoring](https://ethereum.scorechain.com/api_doc/#/scoring/getAddressScoring)

Get scoring of incoming or outgoing funds of an address from the Scorechain Ethereum API.

#### Params

* `endpoint` - `get-scoring`
* `address` - address to be scored `string`
* `direction` - direction of the scoring `{incoming|outgoing}`
* `depth` - depth of the analysis `[1-10]` (default: `3`)

## Test

```
go test
```

## Build

```
go build -o ./build/cl-sc-eth
```

## Lambda

Zip the binary

```
zip cl-sc-eth.zip ./build/cl-sc-eth
```

Then upload to AWS Lambda and use `cl-sc-eth` as the handler.
