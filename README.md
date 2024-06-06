# Currency Exchange Service

## Clone repository
```
git clone https://github.com/BenjuKoirala/currency-exchange-service.git
```
## Installing dependencies
* Clone repo
```
$ cd CurrencyExchangeService
$ go get
```


## Running exchange server
```
$ cd CurrencyExchangeService
```
* Update config.json with API key from ExchangeRate-API (https://app.exchangerate-api.com/keys)
```
$ go run server/server.go
```

## Running tests
```
$ cd CurrencyExchangeService/server
$ go test
```

## Building exchange cli
```
$ cd CurrencyExchangeService/exchange-cli
$ go build -o target/exchange-cli
$ ./target/exchange-cli getrate -base USD -target NPR
```


