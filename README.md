# Chainlink Scorechain Ethereum Adapter

Chainlink adapter for the [Scorechain Ethereum API](https://ethereum.scorechain.com/api_doc/).

## Environment Variables

* `SC_ETH_TOKEN` - token for the Scorechain Ethereum API
* `SC_ETH_URL` - the URL for the Scorechain Ethereum API

## Usage

The Scorechain Ethereum Adapter exposes the following capabilities:

### Get Status

Get the status of the Scorechain Ethereum API.

#### Params

* `endpoint` - `get-status`

### Get Transaction

Get a transaction from the Scorechain Ethereum API.

#### Params

* `endpoint` - `get-trx`
* `hash` - hash of the desired transaction `string`

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
