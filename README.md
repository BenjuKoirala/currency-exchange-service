# Currency Exchange Service

CurrencyExchangeService is a gRPC service that provides exchange rates between different currencies. This service includes a command-line interface (CLI) client for querying exchange rates.

# Features
* Get exchange rates between currencies
* Simple gRPC service implementation using ProtoBufs
* Command-line client for easy access to the service

# Prerequisites
* Go 1.16+
* Protocol Buffers compiler (protoc)
* gRPC and Protocol Buffers Go plugins

## Clone repository
```
git clone https://github.com/BenjuKoirala/currency-exchange-service.git
```
```
cd currency-exchange-service
```

## (Optional) Generate gRPC and Protocol Buffers code:
```
protoc --go_out=. --go-grpc_out=. proto/currency.proto
```

## Building the CLI client
```
go build -o target/exchange-cli ./exchange-cli
```

# Usage
## Running exchange server
To start the gRPC server:
* Update config.json with API key from ExchangeRate-API (https://app.exchangerate-api.com/keys)
```
go run server/server.go
```

## Running tests
```
cd ./server
```
```
go test
```

## Using the CLI client
```
./target/exchange-cli getrate -b USD -t NPR
```
* `-b` specifies the base currency (e.g., USD).
* `-t` specifies the target currency (e.g., NPR).


