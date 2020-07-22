package main

import (
	"fmt"
	"log"

	"github.com/martinusso/getnet"
)

const (
	clientID     = ""
	clientSecret = ""
	sellerID     = ""
)

func main() {
	credentials := credentials()

	card := creditCard(credentials)

	verifyCard(card, credentials)

	paymentWithCreditCard(card, credentials)
}

func credentials() getnet.ClientCredentials {
	credentials := getnet.ClientCredentials{
		SellerID:     sellerID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Sandbox:      true}

	fmt.Println("Geração do token de acesso")

	var err error
	credentials.AccessToken, err = credentials.NewAccessToken()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("access_token: ")
	fmt.Println(credentials.AccessToken.Token)
	fmt.Print("token_type: ")
	fmt.Println(credentials.AccessToken.TokenType)
	fmt.Print("expires_in: ")
	fmt.Println(credentials.AccessToken.ExpiresIn)
	fmt.Print("scope: ")
	fmt.Println(credentials.AccessToken.Scope)

	return credentials
}

func creditCard(credentials getnet.ClientCredentials) getnet.Card {
	fmt.Println("\nGeração do token do cartão")

	card := getnet.Card{
		CardNumber:      "5155901222280001",
		Brand:           getnet.Mastercard,
		CardHolderName:  "Emilio Botín",
		ExpirationYear:  "23",
		ExpirationMonth: "12",
		SecurityCode:    "123",
	}
	var err error
	card.NumberToken, err = card.Token(credentials)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("number_token: ")
	fmt.Println(card.NumberToken)

	return card
}

func paymentWithCreditCard(card getnet.Card, credentials getnet.ClientCredentials) {
	fmt.Println("\nPagamento com cartão de crédito")

	payment := getnet.Payment{
		Amount: 12.34,
		Order: getnet.Order{
			OrderID:     "beed376a-9774-4b8d-80be-a02536e8771f",
			SalesTax:    0,
			ProductType: getnet.Service,
		},
		Customer: getnet.Customer{
			CustomerID:     "ea05ba48-d193-4eb8-a4e9-c9cae1e3e2aa",
			FirstName:      "João",
			LastName:       "da Silva",
			Name:           "João da Silva",
			Email:          "customer@email.com.br",
			DocumentType:   "CPF",
			DocumentNumber: "12345678912",
			PhoneNumber:    "5551999887766",
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
		Device: getnet.Device{
			DeviceID:  "ad38ae20-223c-4875-8797-0749abeb7e08",
			IPAddress: "127.0.0.1",
		},
		Shippings: []getnet.Shipping{
			getnet.Shipping{
				FirstName:   "João",
				Name:        "João da Silva",
				Email:       "customer@email.com.br",
				PhoneNumber: "5551999887766",
				Address: getnet.Address{
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
		},
		Credit: getnet.Credit{
			Card: card,
		},
	}
	pr, err := payment.Pay(credentials)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("payment_id: ")
	fmt.Println(pr.PaymentID)
	fmt.Print("seller_id: ")
	fmt.Println(pr.SellerID)
	fmt.Print("amount: ")
	fmt.Println(pr.Amount)
	fmt.Print("currency: ")
	fmt.Println(pr.Currency)
	fmt.Print("order_id: ")
	fmt.Println(pr.OrderID)
	fmt.Print("status: ")
	fmt.Println(pr.Status)
	fmt.Print("received_at: ")
	fmt.Println(pr.ReceivedAt)
}

func verifyCard(card getnet.Card, credentials getnet.ClientCredentials) {
	fmt.Println("\nVerificação de cartão")

	verification, err := card.Verify(credentials)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("status: ")
	fmt.Println(verification.Status)
	fmt.Print("verification_id: ")
	fmt.Println(verification.VerificationID)
	fmt.Print("authorization_code: ")
	fmt.Println(verification.AuthorizationCode)
}
