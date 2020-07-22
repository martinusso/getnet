package getnet

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPaymentCredit(t *testing.T) {
	server := serverTestPaymentCredit()
	defer server.Close()

	urlStaging = server.URL

	credentials := fixtureCredentials()

	p := Payment{}
	pr, err := p.Pay(credentials)
	if err != nil {
		t.Errorf("There should not be an error, error: %s", err)
	}
	paymentID := "06f256c8-1bbf-42bf-93b4-ce2041bfb87e"
	if pr.PaymentID != paymentID {
		t.Errorf("Expected '%s', got '%s'", paymentID, pr.PaymentID)
	}

	if pr.Amount != 1.23 {
		t.Errorf("Expected '%f', got '%f'", 1.23, pr.Amount)
	}

	expected := "2017-03-19 16:30:30.764 +0000 UTC"
	if pr.ReceivedAt.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, pr.ReceivedAt.String())
	}

	expected = "2017-03-19 16:30:30 +0000 UTC"
	if pr.Credit.AuthorizedAt.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, pr.Credit.AuthorizedAt.String())
	}
}

func serverTestPaymentCredit() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		if req.Header.Get("Content-type") != "application/json; charset=utf-8" {
			err := ErrorResponseSchemaV1{
				Message: "Mensagem detalhada do erro",
				Name:    "Nome do modulo ou sistema em que o erro ocorreu"}
			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode(err)
			return
		}

		auth := req.Header.Get("Authorization")
		if !strings.Contains(auth, "Bearer") {
			err := ErrorResponseSchemaV1{
				Message: "Invalid Authorization",
				Name:    "auth/bearer"}
			rw.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(rw).Encode(err)
			return
		}

		rw.WriteHeader(http.StatusCreated)
		payload := `{
"payment_id": "06f256c8-1bbf-42bf-93b4-ce2041bfb87e",
"seller_id": "6eb2412c-165a-41cd-b1d9-76c575d70a28",
"amount": 123,
"currency": "BRL",
"order_id": "6d2e4380-d8a3-4ccb-9138-c289182818a3",
"status": "APPROVED",
"received_at": "2017-03-19T16:30:30.764Z",
"credit": {
  "delayed": false,
  "authorization_code": "000000099999",
  "authorized_at": "2017-03-19T16:30:30Z",
  "reason_code": "0",
  "reason_message": "transaction approved",
  "acquirer": "GETNET",
  "soft_descriptor": "Descrição para fatura",
  "brand": "Mastercard",
  "terminal_nsu": "0099999",
  "acquirer_transaction_id": "10000024",
  "transaction_id": "1002217281190421"
}
}`
		rw.Write([]byte(payload))
	}))
}
