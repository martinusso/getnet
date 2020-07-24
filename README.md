# getnet

[![Build Status](https://travis-ci.org/martinusso/getnet.svg?branch=master)](https://travis-ci.org/martinusso/getnet)
[![Coverage Status](https://coveralls.io/repos/github/martinusso/getnet/badge.svg?branch=master)](https://coveralls.io/github/martinusso/getnet?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/martinusso/getnet)](https://goreportcard.com/report/github.com/martinusso/getnet)

SDK golang para integração com a API Getnet.

Consulte a documentação oficial da API Getnet https://developers.getnet.com.br/api para maiores detalhes sobre os campos. 

## Funcionalidades

- Autenticação
  - Geração do token de acesso

- Tokenização
  - Geração do token do cartão

- Pagamento
  - Verificação de cartão
  - Pagamento com cartão de crédito

## Usando

```
go get "github.com/martinusso/getnet"
```

### Autenticação - Geração do token de acesso


```
var err error

credentials := getnet.ClientCredentials{
	ID:     "ecb847f2-e423-40c0-808c-55d2098a92ab",
	Secret: "1386f27e-0f2e-45f7-9efd-c8fdc1657426"}
credentials.AccessToken, err = credentials.NewAccessToken()
if err != nil {
	log.Fatal(err)
}
```


### Tokenização - Geração do token do cartão

```
var err error

card := getnet.Card{
	CardNumber:      "5155901222280001",
	CardHolderName:  "JOAO DA SILVA",
	SecurityCode:    "123",
	Brand:           "Mastercard",
	ExpirationMonth: "12",
	ExpirationYear:  "20",
}
card.NumberToken, err = card.Token(credentials)
if err != nil {
	log.Fatal(err)
}
```

### Cartão de Crédito

#### Pagamento com cartão de crédito

```
payment := getnet.Payment{
	Amount:   1.00,
	Currency: "BRL",
	Order: getnet.Order{
		OrderID:     "ea3dae62-1125-4eb4-b3ef-dcb720e8899d",
		SalesTax:    0,
		ProductType: Service,
	},
	Customer: getnet.Customer{
		CustomerID:     "customer_id",
		FirstName:      "João",
		LastName:       "da Silva",
		Email:          "customer@email.com.br",
		DocumentType:   "CPF",
		DocumentNumber: "12345678912",
		PhoneNumber:    "27987654321",
		BillingAddress: getnet.BillingAddress{
			Street:     "Av. Brasil",
			Number:     "1000",
			Complement: "Sala 1",
			District:   "São Geraldo",
			City:       "Porto Alegre",
			State:      "RS",
			Country:    "Brasil",
			PostalCode: "90230060",
		},
	},
	Credit: getnet.Credit{
		Delayed:            false,
		Authenticated:      false,
		PreAuthorization:   false,
		SaveCardData:       false,
		TransactionType:    Full,
		NumberInstallments: 1,
		SoftDescriptor:     "Texto exibido na fatura do cartão do comprador",
		DynamicMCC:         1799,
		Card: card, // Tokenização - Geração do token do cartão
	},
}

// credentials obtido em Autenticação - Geração do token de acesso
response, error := payment.Pay(credentials)

```
