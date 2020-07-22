package getnet

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	numberToken = "dfe05208b105578c070f806c80abd3af09e246827d29b866cf4ce16c205849977c9496cbf0d0234f42339937f327747075f68763537b90b31389e01231d4d13c"
)

func TestCardToken(t *testing.T) {
	server := serverTestTokenCard()
	defer server.Close()

	urlStaging = server.URL

	var err error
	credentials := fixtureCredentials()

	card := Card{
		CardNumber: "5155901222280001"}

	card.NumberToken, err = card.Token(credentials)
	if err != nil {
		t.Errorf("There should not be an error, error: %s", err)
	}
	if card.NumberToken != numberToken {
		t.Errorf("Expected '%s', got '%s'", numberToken, card.NumberToken)
	}
}

func TestCardTokenBadRequest(t *testing.T) {
	errorMessage := "Mensagem detalhada do erro."
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		err := ErrorResponseSchemaV1{Details: []Detail{
			{DescriptionDetail: errorMessage},
		},
			Message: errorMessage}
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
	}))
	defer server.Close()

	urlStaging = server.URL

	var err error
	credentials := fixtureCredentials()

	card := Card{}
	_, err = card.Token(credentials)
	if err.Error() != errorMessage {
		t.Errorf("Expected '%s', got '%s'", errorMessage, err.Error())
	}
}

func TestCardVerify(t *testing.T) {
	server := serverTestVerifyCard()
	defer server.Close()

	urlStaging = server.URL

	var err error
	credentials := fixtureCredentials()

	card := Card{
		CardNumber:  "5155901222280001",
		NumberToken: numberToken,
		Brand:       Mastercard}

	ver, err := card.Verify(credentials)
	if err != nil {
		t.Errorf("There should not be an error, error: %s", err)
	}
	if !ver.Verified() {
		t.Errorf("Expected verified card, got '%s'", ver.Status)
	}
}

func TestCardBrandNotVerified(t *testing.T) {
	server := serverTestVerifyCard()
	defer server.Close()

	urlStaging = server.URL

	var err error
	credentials := fixtureCredentials()

	card := Card{
		CardNumber:  "5155901222280001",
		NumberToken: numberToken}

	ver, err := card.Verify(credentials)
	if err != nil {
		t.Errorf("There should not be an error, error: %s", err)
	}
	if !ver.NotVerified() {
		t.Errorf("Expected a not verified card, got '%s'", ver.Status)
	}
}

func TestCardVerifyNumberTokenRequired(t *testing.T) {
	server := serverTestVerifyCard()
	defer server.Close()

	urlStaging = server.URL

	var err error
	credentials := fixtureCredentials()

	card := Card{
		CardNumber:  "5155901222280001",
		NumberToken: ""}

	ver, err := card.Verify(credentials)
	if err.Error() != errNumberToken.Error() {
		t.Errorf("Expected '%s', got '%s'", errNumberToken, err)
	}
	if ver.Status != "" {
		t.Errorf("Expected an empty status, got '%s'", ver.Status)
	}
}

func serverTestTokenCard() *httptest.Server {
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
		if !strings.Contains(auth, "Bearer ") {
			err := ErrorResponseSchemaV1{
				Message: "Invalid Authorization",
				Name:    "auth/bearer"}
			rw.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(rw).Encode(err)
			return
		}

		token := Token{
			NumberToken: numberToken,
		}
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(token)
	}))
}

func serverTestVerifyCard() *httptest.Server {
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
		if !strings.Contains(auth, "Bearer ") {
			err := ErrorResponseSchemaV1{
				Message: "Invalid Authorization",
				Name:    "auth/bearer"}
			rw.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(rw).Encode(err)
			return
		}

		ver := Verification{
			Status:            Verified,
			VerificationID:    "51d2f214-ffd8-4da3-96ae-97aa16c068d3",
			AuthorizationCode: "6964722471672911",
		}
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(ver)
	}))
}
