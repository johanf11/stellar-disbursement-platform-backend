// Package integration defines the vendor-agnostic interface for a fiat
// compliance + on/off-ramp provider used by the SDP: KYC onboarding plus a
// virtual account (fiat deposit instructions) that settles to the Stellar
// distribution account.
//
// It was extracted from the original bridge.ServiceInterface so multiple
// providers can be plugged in behind one interface. bridge.Service is the
// reference implementation; owlting.Service (internal/owlting) is a second
// implementation. The concrete provider is selected at wiring time in
// cmd/serve.go via configuration.
//
// TODO(bridge-adapter): bridge.Service does not yet satisfy integration.Provider
// — its methods are named OptInToBridge/GetBridgeIntegration and it returns
// *bridge.BridgeIntegrationInfo. Either rename + return *integration.Info, or
// add a thin adapter in this package that maps bridge types onto these.
package integration

import (
	"context"
	"time"
)

// Provider is the vendor-agnostic compliance + fiat on/off-ramp integration.
//
//go:generate mockery --name=Provider --case=underscore --structname=MockProvider --output=. --filename=provider_mock.go --inpackage
type Provider interface {
	// OptIn onboards the organization: creates a KYC link with the provider and
	// persists the integration record.
	OptIn(ctx context.Context, opts OptInOptions) (*Info, error)
	// GetIntegration returns the current integration status, enriched with live
	// data from the provider API when available.
	GetIntegration(ctx context.Context) (*Info, error)
	// CreateVirtualAccount provisions a fiat deposit account (USD -> USDC) that
	// settles to the given Stellar distribution account address.
	CreateVirtualAccount(ctx context.Context, userID, distributionAccountAddress string) (*Info, error)
	// OptInForExistingCustomer onboards using a customer ID that already exists
	// on the provider side (skips KYC link creation).
	OptInForExistingCustomer(ctx context.Context, customerID, userID string) (*Info, error)
}

// CustomerType is the onboarding entity type.
type CustomerType string

const (
	CustomerTypeIndividual CustomerType = "individual"
	CustomerTypeBusiness   CustomerType = "business"
)

// Status is the lifecycle state of a provider integration.
type Status string

const (
	StatusNotOptedIn      Status = "NOT_OPTED_IN"
	StatusOptedIn         Status = "OPTED_IN"
	StatusReadyForDeposit Status = "READY_FOR_DEPOSIT"
	StatusError           Status = "ERROR"
)

// OptInOptions are the inputs required to onboard an organization.
type OptInOptions struct {
	UserID      string
	FullName    string
	Email       string
	RedirectURL string
	KYCType     CustomerType
}

// Info is the composite, provider-agnostic view of an integration.
type Info struct {
	Status                  Status              `json:"status"`
	CustomerID              *string             `json:"customer_id,omitempty"`
	KYCLinkInfo             *KYCLinkInfo        `json:"kyc_status,omitempty"`
	VirtualAccountDetails   *VirtualAccountInfo `json:"virtual_account,omitempty"`
	OptedInBy               *string             `json:"opted_in_by,omitempty"`
	OptedInAt               *time.Time          `json:"opted_in_at,omitempty"`
	VirtualAccountCreatedBy *string             `json:"virtual_account_created_by,omitempty"`
	VirtualAccountCreatedAt *time.Time          `json:"virtual_account_created_at,omitempty"`
}

// KYCLinkInfo is the provider-agnostic KYC link / status.
type KYCLinkInfo struct {
	ID               string       `json:"id"`
	FullName         string       `json:"full_name"`
	Email            string       `json:"email"`
	Type             CustomerType `json:"type"`
	KYCLink          string       `json:"kyc_link"`
	TOSLink          string       `json:"tos_link"`
	KYCStatus        string       `json:"kyc_status"`
	TOSStatus        string       `json:"tos_status"`
	RejectionReasons []string     `json:"rejection_reasons,omitempty"`
	CustomerID       string       `json:"customer_id"`
}

// VirtualAccountInfo is the provider-agnostic fiat deposit account.
type VirtualAccountInfo struct {
	ID                  string               `json:"id"`
	Status              string               `json:"status"`
	CustomerID          string               `json:"customer_id"`
	DepositInstructions *DepositInstructions `json:"deposit_instructions,omitempty"`
	Destination         *Destination         `json:"destination,omitempty"`
}

// DepositInstructions are the fiat (bank) deposit details the depositor uses to
// fund the account.
type DepositInstructions struct {
	BankBeneficiaryName string   `json:"bank_beneficiary_name"`
	Currency            string   `json:"currency"`
	BankName            string   `json:"bank_name"`
	BankAddress         string   `json:"bank_address"`
	BankAccountNumber   string   `json:"bank_account_number"`
	BankRoutingNumber   string   `json:"bank_routing_number"`
	PaymentRails        []string `json:"payment_rails"`
}

// Destination is the on-chain settlement target for the on-ramp.
type Destination struct {
	PaymentRail    string `json:"payment_rail"`
	Currency       string `json:"currency"`
	Address        string `json:"address"`
	BlockchainMemo string `json:"blockchain_memo,omitempty"`
}
