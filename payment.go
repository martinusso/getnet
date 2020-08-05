package getnet

import (
	"encoding/json"
	"time"
)

type ProductType string
type TransactionType string
type Currency string

const (
	endpointPaymentCredit = "/v1/payments/credit"

	// Identificador do tipo de produto vendido dentre as opções (product_type)
	CashCarry       ProductType = "cash_carry"
	DigitalContent  ProductType = "digital_content"
	DigitalGoods    ProductType = "digital_goods"
	DigitalPhysical ProductType = "digital_physical"
	GiftCard        ProductType = "gift_card"
	PhysicalGoods   ProductType = "physical_goods"
	RenewSubs       ProductType = "renew_subs"
	Shareware       ProductType = "shareware"
	Service         ProductType = "service"

	// Tipo de transação (transaction_type)
	// Pagamento completo à vista, parcelado sem juros, parcelado com juros.
	Full                TransactionType = "FULL"
	InstallNoInterest   TransactionType = "INSTALL_NO_INTEREST"
	InstallWithInterest TransactionType = "INSTALL_WITH_INTEREST"

	// Status da transação (status)
	PaymentCanceled   = "CANCELED"
	PaymentApproved   = "APPROVED"
	PaymentDenied     = "DENIED"
	PaymentAuthorized = "AUTHORIZED"
	PaymentConfirmed  = "CONFIRMED"

	RealBrazilian Currency = "BRL"
	DollarUS      Currency = "USD"
)

type Payment struct {
	SellerID  string     `json:"seller_id,omitempty"`
	Amount    float64    `json:"amount"`
	Currency  Currency   `json:"currency"`
	Order     Order      `json:"order"`
	Customer  Customer   `json:"customer"`
	Device    Device     `json:"device,omitempty"`
	Shippings []Shipping `json:"shippings,omitempty"`
	Credit    Credit     `json:"credit,omitempty"`
	Debit     Debit      `json:"debit,omitempty"`
}

func (p Payment) MarshalJSON() ([]byte, error) {
	type Alias Payment
	return json.Marshal(&struct {
		Alias
		Amount int `json:"amount"`
	}{
		Alias:  (Alias)(p),
		Amount: int(p.Amount * 100),
	})
}

func (p Payment) Pay(c ClientCredentials) (PaymentResponse, error) {
	if p.Currency == "" {
		p.Currency = RealBrazilian
	}
	if p.Credit.TransactionType == "" {
		p.Credit.TransactionType = Full
	}
	if p.Credit.NumberInstallments < 1 {
		p.Credit.NumberInstallments = 1
	}
	res, err := NewRestClient(c).Post(endpointPaymentCredit, p)
	if err != nil {
		return PaymentResponse{}, err
	}

	var pr PaymentResponse
	err = json.Unmarshal(res.Body, &pr)
	return pr, err
}

type Credit struct {
	Delayed            bool            `json:"delayed"`
	Authenticated      bool            `json:"authenticated"`
	PreAuthorization   bool            `json:"pre_authorization"`
	SaveCardData       bool            `json:"save_card_data"`
	TransactionType    TransactionType `json:"transaction_type"`
	NumberInstallments int             `json:"number_installments"`
	SoftDescriptor     string          `json:"soft_descriptor"`
	DynamicMCC         int             `json:"dynamic_mcc"`
	Card               Card            `json:"card"`
}

func (c Credit) MarshalJSON() ([]byte, error) {
	type Alias Credit
	return json.Marshal(&struct {
		Alias
		SoftDescriptor string `json:"soft_descriptor,omitempty"`
	}{
		Alias:          (Alias)(c),
		SoftDescriptor: maxLength(c.SoftDescriptor, 22),
	})
}

type Debit struct {
	CardHolderMobile string `json:"cardholder_mobile"`
	SoftDescriptor   string `json:"soft_descriptor"`
	DynamicMCC       int    `json:"dynamic_mcc"`
	Authenticated    bool   `json:"authenticated"`
	Card             Card   `json:"card"`
}

func (d Debit) MarshalJSON() ([]byte, error) {
	type Alias Debit
	return json.Marshal(&struct {
		Alias
		SoftDescriptor string `json:"soft_descriptor,omitempty"`
	}{
		Alias:          (Alias)(d),
		SoftDescriptor: maxLength(d.SoftDescriptor, 22),
	})
}

type BillingAddress struct {
	Street     string `json:"street,omitempty"`
	Number     string `json:"number,omitempty"`
	Complement string `json:"complement,omitempty"`
	District   string `json:"district,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	Country    string `json:"country,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
}

type Customer struct {
	CustomerID     string         `json:"customer_id"`
	FirstName      string         `json:"first_name,omitempty"`
	LastName       string         `json:"last_name,omitempty"`
	Name           string         `json:"name,omitempty"`
	Email          string         `json:"email,omitempty"`
	DocumentType   string         `json:"document_type,omitempty"`
	DocumentNumber string         `json:"document_number,omitempty"`
	PhoneNumber    string         `json:"phone_number,omitempty"`
	BillingAddress BillingAddress `json:"billing_address,omitempty"`
}

type Device struct {
	DeviceID  string `json:"device_id,omitempty"`
	IPAddress string `json:"ip_address,omitempty"`
}

type Order struct {
	OrderID     string      `json:"order_id"`
	SalesTax    int         `json:"sales_tax"`
	ProductType ProductType `json:"product_type"`
}

type Shipping struct {
	FirstName      string  `json:"first_name,omitempty"`
	Name           string  `json:"name,omitempty"`
	Email          string  `json:"email,omitempty"`
	PhoneNumber    string  `json:"phone_number,omitempty"`
	ShippingAmount float64 `json:"shipping_amount,omitempty"`
	Address        Address `json:"address,omitempty"`
}

func (s Shipping) MarshalJSON() ([]byte, error) {
	type Alias Shipping
	return json.Marshal(&struct {
		Alias
		ShippingAmount int `json:"shipping_amount"`
	}{
		Alias:          (Alias)(s),
		ShippingAmount: int(s.ShippingAmount * 100),
	})
}

type Address struct {
	Street     string `json:"street,omitempty"`
	Number     string `json:"number,omitempty"`
	Complement string `json:"complement,omitempty"`
	District   string `json:"district,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	Country    string `json:"country,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
}

type PaymentResponse struct {
	PaymentID  string         `json:"payment_id"`
	SellerID   string         `json:"seller_id"`
	Amount     float64        `json:"amount"`
	Currency   Currency       `json:"currency"`
	OrderID    string         `json:"order_id"`
	Status     string         `json:"status"`
	ReceivedAt time.Time      `json:"received_at"`
	Credit     CreditResponse `json:"credit"`
}

func (p PaymentResponse) Canceled() bool {
	return p.Status == PaymentCanceled
}

func (p PaymentResponse) Approved() bool {
	return p.Status == PaymentApproved
}

func (p PaymentResponse) Denied() bool {
	return p.Status == PaymentDenied
}

func (p PaymentResponse) Authorized() bool {
	return p.Status == PaymentAuthorized
}

func (p PaymentResponse) Confirmed() bool {
	return p.Status == PaymentConfirmed
}

func (p *PaymentResponse) UnmarshalJSON(data []byte) error {
	type Alias PaymentResponse
	aux := &struct {
		*Alias
		Amount     int    `json:"amount"`
		ReceivedAt string `json:"received_at"`
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	p.Amount = float64(aux.Amount) / 100
	p.ReceivedAt, _ = time.Parse("2006-01-02T15:04:05.000Z", aux.ReceivedAt)
	return nil
}

type CreditResponse struct {
	Delayed               bool      `json:"delayed"`
	AuthorizationCode     string    `json:"authorization_code"`
	AuthorizedAt          time.Time `json:"authorized_at"`
	ReasonCode            string    `json:"reason_code"`
	ReasonMessage         string    `json:"reason_message"`
	Acquirer              string    `json:"acquirer"`
	SoftDescriptor        string    `json:"soft_descriptor"`
	Brand                 Brand     `json:"brand"`
	TerminalNSU           string    `json:"terminal_nsu"`
	AcquirerTransactionID string    `json:"acquirer_transaction_id"`
	TransactionID         string    `json:"transaction_id"`
}

func (cr *CreditResponse) UnmarshalJSON(data []byte) error {
	type Alias CreditResponse
	aux := &struct {
		*Alias
		AuthorizedAt string `json:"authorized_at"`
	}{
		Alias: (*Alias)(cr),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	cr.AuthorizedAt, _ = time.Parse("2006-01-02T15:04:05Z", aux.AuthorizedAt)
	return nil
}
