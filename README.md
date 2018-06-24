# theperfectgiftcard-api
[![GoDoc](https://godoc.org/github.com/rayou/theperfectgiftcard-api?status.svg)](https://godoc.org/github.com/rayou/theperfectgiftcard-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/rayou/theperfectgiftcard-api)](https://goreportcard.com/report/github.com/rayou/theperfectgiftcard-api)
[![Coverage Status](https://coveralls.io/repos/github/rayou/theperfectgiftcard-api/badge.svg)](https://coveralls.io/github/rayou/theperfectgiftcard-api)

The Perfect Gift Card API - Fetch your card summary and transaction history via an API. 

Built on top of [go-theperfectgiftcard](https://github.com/rayou/go-theperfectgiftcard) library.

## Example
```sh

# Change card_no and pin to yours.
$ curl --request POST \
  --url https://theperfectgiftcard-api-lhsiiccogb.now.sh/card \
  --data '{
	"card_no": "50211234567890",
	"pin": "0000"
}'

{
  "CardNo": "50211234567890",
  "AccountNo": "000000000",
  "LoadsToDate": "$100.00",
  "PurchasesToDate": "-$54.32",
  "AvailableBalance": "$12.34",
  "PurchasedDate": "1 Jan 2018",
  "ExpiryDate": "1 Jan 2021",
  "Transactions": [
    {
      "Date": "1 Jan 2018 12:04:45 PM",
      "Details": "Store Address",
      "Description": "Refund - Store Address",
      "Amount": "$100.00",
      "Balance": "$100.00"
    },
    {
      "Date": "2 Jan 2018 07:50:53 PM",
      "Details": "Store A",
      "Description": "Purchase - Store A",
      "Amount": "$12.34-",
      "Balance": "$56.78"
    }
  ]
}
```

## How To Use

### Run in docker
```sh
# Clone this repo and switch working folder
$ git clone https://github.com/rayou/theperfectgiftcard-api.git
$ cd theperfectgiftcard-api/

# build docker image
$ docker build -t rayou/theperfectgiftcard-api ./ 

# run in docker
$ docker run --rm -it -p 8000:8000 rayou/theperfectgiftcard-api

# test 
$ curl --request POST \
  --url localhost:8000/card \
  --data '{
	"card_no": "50211234567890",
	"pin": "0000"
}'
```

### Deploy to [Now](https://now.sh)
```sh
# Clone this repo and switch working folder
$ git clone https://github.com/rayou/theperfectgiftcard-api.git
$ cd theperfectgiftcard-api/ 

# install Now cli
$ npm install -g now 

# deploy to Now
$ now

# test
$ curl --request POST \
  --url https://{YOUR_NOW_APP_URL}/card \
  --data '{
	"card_no": "50211234567890",
	"pin": "0000"
}'

```

## Contributing

PRs are welcome.

## Author

Ray Ou - yuhung.ou@live.com

## License

MIT.
