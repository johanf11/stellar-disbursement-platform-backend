package owlting

import (
	"context"
	"fmt"

	"github.com/stellar/stellar-disbursement-platform-backend/internal/data"
	"github.com/stellar/stellar-disbursement-platform-backend/internal/integration"
	"github.com/stellar/stellar-disbursement-platform-backend/internal/services"
	"github.com/stellar/stellar-disbursement-platform-backend/internal/transactionsubmission/engine/signing"
	"github.com/stellar/stellar-disbursement-platform-backend/internal/utils"
)

// Service is the OwlPay Harbor implementation of integration.Provider.
//
// It mirrors bridge.Service so wiring is identical and the disbursement ramps
// (USD on-ramp, HTG/FX off-ramp, KYC) can be served by OwlTing. Every provider
// call is stubbed with a TODO pending the real OwlPay Harbor API contract.
//
// NOTE: the importer "HTG collection" rail (accepting gourdes from Haitian
// importers) is intentionally NOT part of this provider — OwlTing does not
// on-ramp HTG inside Haiti. That rail is Theo Haiti Ops and is tracked
// separately (see THEO_TREASURY_MODEL.md).
type Service struct {
	client                      ClientInterface
	baseURL                     string
	apiKey                      string
	models                      *data.Models
	distributionAccountResolver signing.DistributionAccountResolver
	distributionAccountService  services.DistributionAccountServiceInterface
	networkType                 utils.NetworkType
}

var _ integration.Provider = (*Service)(nil)

// ServiceOptions configures the OwlTing service. Mirrors bridge.ServiceOptions.
type ServiceOptions struct {
	BaseURL                     string
	APIKey                      string
	Models                      *data.Models
	DistributionAccountResolver signing.DistributionAccountResolver
	DistributionAccountService  services.DistributionAccountServiceInterface
	NetworkType                 utils.NetworkType
}

// Validate validates the OwlTing service options.
func (o ServiceOptions) Validate() error {
	if o.BaseURL == "" {
		return fmt.Errorf("baseURL is required")
	}
	if o.APIKey == "" {
		return fmt.Errorf("apiKey is required")
	}
	if o.Models == nil {
		return fmt.Errorf("models is required")
	}
	if o.DistributionAccountResolver == nil {
		return fmt.Errorf("distributionAccountResolver is required")
	}
	if o.DistributionAccountService == nil {
		return fmt.Errorf("distributionAccountService is required")
	}
	if err := o.NetworkType.Validate(); err != nil {
		return fmt.Errorf("validating NetworkType: %w", err)
	}
	return nil
}

// NewService creates a new OwlTing service instance.
func NewService(opts ServiceOptions) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validating owlting.Service options: %w", err)
	}

	client, err := NewClient(ClientOptions{BaseURL: opts.BaseURL, APIKey: opts.APIKey})
	if err != nil {
		return nil, fmt.Errorf("creating OwlTing client: %w", err)
	}

	return &Service{
		client:                      client,
		baseURL:                     opts.BaseURL,
		apiKey:                      opts.APIKey,
		models:                      opts.Models,
		distributionAccountResolver: opts.DistributionAccountResolver,
		distributionAccountService:  opts.DistributionAccountService,
		networkType:                 opts.NetworkType,
	}, nil
}

// OptIn onboards the organization. Mirrors bridge.Service.OptInToBridge.
func (s *Service) OptIn(ctx context.Context, opts integration.OptInOptions) (*integration.Info, error) {
	// TODO 1: validate USDC trustline on the distribution account
	//         (s.distributionAccountResolver + s.distributionAccountService + s.networkType).
	// TODO 2: check for an existing integration record via s.models.
	// TODO 3: s.client.PostKYCLink(...) and map KYCLinkResponse -> integration.KYCLinkInfo.
	// TODO 4: persist the integration record (generalize data.BridgeIntegration or add an
	//         owlting-backed table / provider column) and return integration.Info.
	return nil, fmt.Errorf("owlting.OptIn: %w", ErrNotImplemented)
}

// GetIntegration returns current status enriched with live OwlPay Harbor data.
// Mirrors bridge.Service.GetBridgeIntegration.
func (s *Service) GetIntegration(ctx context.Context) (*integration.Info, error) {
	// TODO: read the integration record; when present, enrich with
	//       s.client.GetCustomer / GetKYCLink / GetVirtualAccount and map to integration.Info.
	return nil, fmt.Errorf("owlting.GetIntegration: %w", ErrNotImplemented)
}

// CreateVirtualAccount provisions the USD->USDC on-ramp. Mirrors
// bridge.Service.CreateVirtualAccount.
func (s *Service) CreateVirtualAccount(ctx context.Context, userID, distributionAccountAddress string) (*integration.Info, error) {
	// TODO 1: load the integration record; require OPTED_IN and no existing virtual account.
	// TODO 2: verify KYC approved + TOS accepted + customer active via s.client.
	// TODO 3: build VirtualAccountRequest (usd -> usdc on stellar, tenant memo) and
	//         s.client.PostVirtualAccount(...).
	// TODO 4: persist the virtual account ID + READY_FOR_DEPOSIT status; map
	//         VirtualAccountResponse -> integration.VirtualAccountInfo.
	return nil, fmt.Errorf("owlting.CreateVirtualAccount: %w", ErrNotImplemented)
}

// OptInForExistingCustomer onboards with a pre-existing provider customer ID.
// Mirrors bridge.Service.OptInForExistingCustomer.
func (s *Service) OptInForExistingCustomer(ctx context.Context, customerID, userID string) (*integration.Info, error) {
	// TODO: validate the customer via s.client.GetCustomer; persist the integration record.
	return nil, fmt.Errorf("owlting.OptInForExistingCustomer: %w", ErrNotImplemented)
}
