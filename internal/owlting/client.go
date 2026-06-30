package owlting

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// ErrNotImplemented is returned by every stubbed OwlPay Harbor call until the
// real endpoints are wired in.
var ErrNotImplemented = errors.New("owlting: not implemented")

// ClientInterface is the low-level OwlPay Harbor HTTP client.
//
// TODO: map each method to the real OwlPay Harbor API endpoint + auth header.
// Reference: https://www.owlting.com/owlpay/harbor (request partner/developer API docs).
//
//go:generate mockery --name=ClientInterface --case=underscore --structname=MockClient --output=. --filename=client_mock.go --inpackage
type ClientInterface interface {
	PostKYCLink(ctx context.Context, req KYCLinkRequest) (*KYCLinkResponse, error)
	GetKYCLink(ctx context.Context, kycLinkID string) (*KYCLinkResponse, error)
	GetCustomer(ctx context.Context, customerID string) (*CustomerResponse, error)
	PostVirtualAccount(ctx context.Context, customerID string, req VirtualAccountRequest) (*VirtualAccountResponse, error)
	GetVirtualAccount(ctx context.Context, customerID, virtualAccountID string) (*VirtualAccountResponse, error)
}

// Client is the OwlPay Harbor HTTP client.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

var _ ClientInterface = (*Client)(nil)

// ClientOptions configures the OwlPay Harbor client.
type ClientOptions struct {
	BaseURL string
	APIKey  string
}

// NewClient builds an OwlPay Harbor client.
func NewClient(opts ClientOptions) (*Client, error) {
	if opts.BaseURL == "" {
		return nil, fmt.Errorf("owlting: baseURL is required")
	}
	if opts.APIKey == "" {
		return nil, fmt.Errorf("owlting: apiKey is required")
	}
	return &Client{
		baseURL:    opts.BaseURL,
		apiKey:     opts.APIKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// --- wire types: TODO align field names/enums with the OwlPay Harbor API schema ---

// KYCLinkRequest is the payload to create a hosted KYC link.
type KYCLinkRequest struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Type        string `json:"type"` // individual | business
	RedirectURI string `json:"redirect_uri,omitempty"`
	// TODO: add OwlPay Harbor-specific fields (country, corridor, endorsements, etc.).
}

// KYCLinkResponse is the hosted KYC link + status.
type KYCLinkResponse struct {
	ID         string   `json:"id"`
	FullName   string   `json:"full_name"`
	Email      string   `json:"email"`
	KYCLink    string   `json:"kyc_link"`
	TOSLink    string   `json:"tos_link"`
	KYCStatus  string   `json:"kyc_status"`
	TOSStatus  string   `json:"tos_status"`
	CustomerID string   `json:"customer_id"`
	Rejections []string `json:"rejection_reasons,omitempty"`
	// TODO: align with the OwlPay Harbor KYC response (status enum values may differ).
}

// CustomerResponse is the provider customer record.
type CustomerResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	// TODO: align with the OwlPay Harbor customer schema.
}

// VirtualAccountRequest provisions a fiat deposit account that settles on-chain.
type VirtualAccountRequest struct {
	SourceCurrency      string `json:"source_currency"`      // e.g. usd
	DestinationRail     string `json:"destination_rail"`     // e.g. stellar
	DestinationCurrency string `json:"destination_currency"` // e.g. usdc
	DestinationAddress  string `json:"destination_address"`
	BlockchainMemo      string `json:"blockchain_memo,omitempty"`
	// TODO: align with the OwlPay Harbor on-ramp / virtual account request.
}

// VirtualAccountResponse is the provisioned deposit account + instructions.
type VirtualAccountResponse struct {
	ID                string   `json:"id"`
	Status            string   `json:"status"`
	CustomerID        string   `json:"customer_id"`
	BankBeneficiary   string   `json:"bank_beneficiary_name"`
	BankName          string   `json:"bank_name"`
	BankAddress       string   `json:"bank_address"`
	BankAccountNumber string   `json:"bank_account_number"`
	BankRoutingNumber string   `json:"bank_routing_number"`
	Currency          string   `json:"currency"`
	PaymentRails      []string `json:"payment_rails"`
	// TODO: align with the OwlPay Harbor virtual account / deposit-instructions schema.
}

// --- stubbed client methods: TODO implement against the real OwlPay Harbor API ---

func (c *Client) PostKYCLink(ctx context.Context, req KYCLinkRequest) (*KYCLinkResponse, error) {
	// TODO: POST {baseURL}/<kyc-link-endpoint> with apiKey auth; decode KYCLinkResponse.
	return nil, fmt.Errorf("owlting.PostKYCLink: %w", ErrNotImplemented)
}

func (c *Client) GetKYCLink(ctx context.Context, kycLinkID string) (*KYCLinkResponse, error) {
	// TODO: GET {baseURL}/<kyc-link-endpoint>/{kycLinkID}.
	return nil, fmt.Errorf("owlting.GetKYCLink: %w", ErrNotImplemented)
}

func (c *Client) GetCustomer(ctx context.Context, customerID string) (*CustomerResponse, error) {
	// TODO: GET {baseURL}/<customer-endpoint>/{customerID}.
	return nil, fmt.Errorf("owlting.GetCustomer: %w", ErrNotImplemented)
}

func (c *Client) PostVirtualAccount(ctx context.Context, customerID string, req VirtualAccountRequest) (*VirtualAccountResponse, error) {
	// TODO: POST {baseURL}/<virtual-account-endpoint> scoped to customerID.
	return nil, fmt.Errorf("owlting.PostVirtualAccount: %w", ErrNotImplemented)
}

func (c *Client) GetVirtualAccount(ctx context.Context, customerID, virtualAccountID string) (*VirtualAccountResponse, error) {
	// TODO: GET {baseURL}/<virtual-account-endpoint>/{virtualAccountID} scoped to customerID.
	return nil, fmt.Errorf("owlting.GetVirtualAccount: %w", ErrNotImplemented)
}
