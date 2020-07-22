package getnet

import (
	"encoding/json"
	"errors"
)

var errNumberToken = errors.New("Obrigatório informar o número do cartão tokenizado. Gerado previamente por meio do endpoint /v1/tokens/card.")

type Brand string

const (
	endpointTokenCard        = "/v1/tokens/card"
	endpointCardVerification = "/v1/cards/verification"

	Mastercard Brand = "Mastercard"
	Visa       Brand = "Visa"
	Amex       Brand = "Amex"
	Elo        Brand = "Elo"
	Hipercard  Brand = "Hipercard"

	Verified    = "VERIFIED"
	NotVerified = "NOT VERIFIED"
	Denied      = "DENIED"
	Error       = "ERROR"
)

type Card struct {
	CardNumber      string `json:"-"`
	NumberToken     string `json:"number_token"`
	Brand           Brand  `json:"brand"`
	CardHolderName  string `json:"cardholder_name"`
	ExpirationMonth string `json:"expiration_month"`
	ExpirationYear  string `json:"expiration_year"`
	SecurityCode    string `json:"security_code"`
	CustomerID      string `json:"-"`
}

func (c Card) Token(cc ClientCredentials) (string, error) {
	payload := struct {
		CardNumber string `json:"card_number"`
		CustomerID string `json:"customer_id,omitempty"`
	}{
		CardNumber: c.CardNumber,
		CustomerID: c.CustomerID,
	}

	res, err := NewRestClient(cc).Post(endpointTokenCard, payload)
	if err != nil {
		return "", err
	}

	var token Token
	err = json.Unmarshal(res.Body, &token)
	return token.NumberToken, err
}

func (c Card) Verify(cc ClientCredentials) (Verification, error) {
	if c.NumberToken == "" {
		return Verification{}, errNumberToken
	}

	if c.Brand != Mastercard && c.Brand != Visa {
		return Verification{Status: NotVerified}, nil
	}

	res, err := NewRestClient(cc).Post(endpointCardVerification, c)
	if err != nil {
		return Verification{}, err
	}

	var ver Verification
	err = json.Unmarshal(res.Body, &ver)
	return ver, err
}

type Token struct {
	NumberToken string `json:"number_token"`
}

type Verification struct {
	Status            string `json:"status"`
	VerificationID    string `json:"verification_id"`
	AuthorizationCode string `json:"authorization_code"`
}

func (v Verification) Verified() bool {
	return v.Status == Verified
}

func (v Verification) NotVerified() bool {
	return v.Status == NotVerified
}

func (v Verification) Denied() bool {
	return v.Status == Denied
}
