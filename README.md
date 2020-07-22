# getnet

SDK golang para integração com a API Getnet.

Consulte a documentação oficial da API Getnet https://developers.getnet.com.br/api para maiores detalhes sobre os campos. 


## Cartão de Crédito

### Pagamento com cartão de crédito

```
var err error
credentials := ClientCredentials{
	ID:     "1",
	Secret: "A"}
credentials.AccessToken, err = credentials.NewAccessToken()
if err != nil {
	// error
}

payment := Payment{
	Amount:   100, // Valor da compra em centavos (ex: R$ 1.00 = 100)
	Currency: "BRL",
	Order: Order{
		OrderID:     "ea3dae62-1125-4eb4-b3ef-dcb720e8899d",
		SalesTax:    0,
		ProductType: Service,
	},
	Customer: Customer{
		CustomerID:     "customer_id",
		FirstName:      "João",
		LastName:       "da Silva",
		Email:          "customer@email.com.br",
		DocumentType:   "CPF",
		DocumentNumber: "12345678912",
		PhoneNumber:    "27987654321", // Telefone do comprador. (sem máscara)
		BillingAddress: BillingAddress{
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
	Credit: Credit{
		Delayed:            false,
		Authenticated:      false,
		Pre_authorization:  false,
		SaveCardData:       false,
		TransactionType:    Full,
		NumberInstallments: 1,
		SoftDescriptor:     "Texto exibido na fatura do cartão do comprador",
		DynamicMCC:         1799,
		Card: Card{
			CardNumber:      "5155901222280001",
			CardHolderName:  "JOAO DA SILVA",
			SecurityCode:    "123",
			Brand:           "Mastercard",
			ExpirationMonth: "12",
			ExpirationYear:  "20",
		},
	},
}

response, error := payment.Pay(credentials)
```
