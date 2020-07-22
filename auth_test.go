package getnet

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const (
	token     = "7cdc8d2f-98e3-49b2-9129-fdf0f389c11c"
	tokenType = "Bearer"
	expiresIn = 3600
	scope     = "oob"
)

func TestNewAccessToken(t *testing.T) {
	server := serverTestAuth()
	defer server.Close()

	urlStaging = server.URL

	var err error
	c := fixtureCredentials()
	c.AccessToken, err = c.NewAccessToken()
	if err != nil {
		t.Errorf("There should not be an error, error: %s", err)
	}
	if c.AccessToken.Token != token {
		t.Errorf("Expected '%s', got '%s'", token, c.AccessToken.Token)
	}
	if c.AccessToken.ExpiresIn != expiresIn {
		t.Errorf("Expected '%d', got '%d'", expiresIn, c.AccessToken.ExpiresIn)
	}
}

func TestNewAccessTokenUnauthorized(t *testing.T) {
	errorMessage := "Não autorizado."
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		err := ErrorResponseSchemaV2{Description: errorMessage}
		rw.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(rw).Encode(err)
	}))
	defer server.Close()

	urlStaging = server.URL

	var err error
	c := fixtureCredentials()
	_, err = c.NewAccessToken()
	if err.Error() != errorMessage {
		t.Errorf("Expected '%s', got '%s'", errorMessage, err.Error())
	}
}

func TestExpired(t *testing.T) {
	at := AccessToken{
		ExpiresIn: 11,
		createdAt: time.Now(),
	}
	if at.Expired() {
		t.Errorf("Expected not expired token")
	}
	time.Sleep(1001 * time.Millisecond)

	if !at.Expired() {
		t.Errorf("Expected expired token")
	}
}

func serverTestAuth() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		if req.Header.Get("Content-type") != "application/x-www-form-urlencoded" {
			err := ErrorResponseSchemaV2{
				Error:       "Mensagem de erro",
				Description: "Descrição do erro"}
			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode(err)
			return
		}

		auth := req.Header.Get("Authorization")
		if !strings.Contains(auth, "Basic ") {
			err := ErrorResponseSchemaV2{
				Error:       "Mensagem de erro",
				Description: "Descrição do erro"}
			rw.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(rw).Encode(err)
			return
		}

		at := AccessToken{
			Token:     token,
			TokenType: tokenType,
			ExpiresIn: expiresIn,
			Scope:     scope,
		}
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(at)
	}))
}

func fixtureCredentials() ClientCredentials {
	return ClientCredentials{
		ClientID:     "client-credentials-id-1",
		ClientSecret: "client-credentials-secret-A",
		Sandbox:      true,
		AccessToken: AccessToken{
			Token: "5d6a5e20-01ed-4672-8e20-690ce727deb8",
		}}
}
