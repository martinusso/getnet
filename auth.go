package getnet

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	authTokenURL = "/auth/oauth/v2/token"
)

var (
	urlStaging    = "https://api-sandbox.getnet.com.br"
	urlProduction = "https://api.getnet.com.br"
)

type ClientCredentials struct {
	ClientID     string
	ClientSecret string
	SellerID     string
	Sandbox      bool
	AccessToken  AccessToken
}

func (cc ClientCredentials) Basic() string {
	basic := fmt.Sprintf("%s:%s", cc.ClientID, cc.ClientSecret)
	token := base64.StdEncoding.EncodeToString([]byte(basic))
	return fmt.Sprintf("Basic %s", token)
}

func (cc ClientCredentials) Bearer() string {
	return fmt.Sprintf("Bearer %s", cc.AccessToken.Token)
}

func (cc ClientCredentials) HasSeller() bool {
	return strings.TrimSpace(cc.SellerID) != ""
}

func (cc ClientCredentials) NewAccessToken() (AccessToken, error) {
	formData := url.Values{}
	formData.Add("scope", "oob")
	formData.Add("grant_type", "client_credentials")
	res, err := NewRestClient(cc).AuthBasic().FormData(authTokenURL, formData)
	if err != nil {
		return AccessToken{}, err
	}

	var at AccessToken
	err = json.Unmarshal(res.Body, &at)
	at.createdAt = time.Now()
	return at, err
}

func (cc ClientCredentials) URL() string {
	if cc.Sandbox {
		return urlStaging
	}
	return urlProduction
}

type AccessToken struct {
	Token     string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
	Scope     string `json:"scope"`
	createdAt time.Time
}

func (at AccessToken) Expired() bool {
	return time.Since(at.createdAt) > time.Duration(at.ExpiresIn-10)*time.Second
}
