package getnet

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type RestClient struct {
	credentials ClientCredentials
	authBasic   bool
}

func NewRestClient(c ClientCredentials) RestClient {
	return RestClient{
		credentials: c}
}

func (r RestClient) AuthBasic() RestClient {
	r.authBasic = true
	return r
}

func (r RestClient) FormData(endpoint string, form url.Values) (Response, error) {
	contentType := "application/x-www-form-urlencoded"
	return r.send(http.MethodPost, endpoint, contentType, strings.NewReader(form.Encode()))
}

func (r RestClient) Get(endpoint string) (Response, error) {
	contentType := "application/json; charset=utf-8"
	return r.send(http.MethodGet, endpoint, contentType, nil)
}

func (r RestClient) Post(endpoint string, value interface{}) (Response, error) {
	contentType := "application/json; charset=utf-8"
	body, err := json.Marshal(value)
	if err != nil {
		return Response{}, err
	}
	return r.send(http.MethodPost, endpoint, contentType, bytes.NewReader(body))
}

func (r RestClient) send(method, endpoint, contentType string, body io.Reader) (Response, error) {
	url := r.credentials.URL() + endpoint
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return Response{}, err
	}

	req.Header.Add("Content-type", contentType)
	if r.authBasic {
		req.Header.Add("Authorization", r.credentials.Basic())
	} else {
		req.Header.Add("Authorization", r.credentials.Bearer())
	}
	if r.credentials.HasSeller() {
		req.Header.Add("seller_id", r.credentials.SellerID)
	}

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Response{}, err
	}

	err = r.getError(res.StatusCode, content, endpoint)
	return Response{
		Body: content,
		Code: res.StatusCode,
	}, err
}

func (r RestClient) getError(statusCode int, payload []byte, endpoint string) error {
	if strings.HasPrefix(strconv.Itoa(statusCode), "2") {
		return nil
	}

	if strings.Contains(endpoint, "/v1/") {
		var errResponse ErrorResponseSchemaV1
		if err := json.Unmarshal(payload, &errResponse); err != nil {
			return err
		}
		return errors.New(errResponse.String())
	}

	var errResponse ErrorResponseSchemaV2
	if err := json.Unmarshal(payload, &errResponse); err != nil {
		return err
	}
	return errors.New(errResponse.Description)
}

type Response struct {
	Body []byte
	Code int
}

type ErrorResponseSchemaV1 struct {
	Message string   `json:"message"`
	Name    string   `json:"name"`
	Details []Detail `json:"details"`
}

func (e ErrorResponseSchemaV1) String() string {
	var errs []string
	for _, d := range e.Details {
		errs = append(errs, d.DescriptionDetail)
	}
	return strings.Join(errs, "; ")
}

type Detail struct {
	Status            string `json:"status"`
	ErroCode          string `json:"error_code"`
	Description       string `json:"description"`
	DescriptionDetail string `json:"description_detail"`
}

type ErrorResponseSchemaV2 struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}
