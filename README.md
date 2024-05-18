# Ethereum Subscriber

This repository was created during a hiring process and presents the solution to a homework task.

## Usage

The service can be started with `go run main.go`. You can then subscribe to an ethereum address using

```shell
curl "http://localhost:8000/subscribe?address=0xdAC17F958D2ee523a2206206994597C13D831ec7"
```

The timer is set to 12 seconds so you need to wait a bit for something to happen.

You can then check the transactions for an address using

```shell
curl "http://localhost:8000/transactions?address=0xdAC17F958D2ee523a2206206994597C13D831ec7" | jq
```

## File Structure

`main.go` starts a very REST(ish) API. It creates a store, gateway and passes those to the block poller.
The block poller runs in a go routine on a timer every 12s.

The store is an in memory store that uses maps internally to save the subscriptions and transactions.
It has an interface that the block poller uses so the actual implementation can be changed to a persistent store later if needed.
It uses a mutex to avoid race conditions.

The gateway is responsible for calling out to the Ethereum gateway and deserializing the results.

The block poller runs periodically.
It gets the current block and checks if any of the transactions contained in it have active subscriptions, if so it adds them to the store.
