package wallets

import (
	"github.com/stellar/stellar-disbursement-platform-backend/internal/data"
	"github.com/stellar/stellar-disbursement-platform-backend/internal/services/assets"
)

var TestnetWallets = []data.Wallet{
	{
		Name:              "Demo Wallet",
		Homepage:          "https://demo-wallet.stellar.org",
		DeepLinkSchema:    "https://demo-wallet.stellar.org",
		SEP10ClientDomain: "demo-wallet-server.stellar.org",
		Assets: []data.Asset{
			assets.USDCAssetTestnet,
			assets.XLMAsset,
		},
	},
	{
		Name:              "Decaf",
		Homepage:          "https://decaf.so",
		DeepLinkSchema:    "https://decafwallet.app.link",
		SEP10ClientDomain: "decaf.so",
		Assets:            assets.AllAssetsTestnet,
	},
	{
		Name:              "Vesseo",
		Homepage:          "https://vesseoapp.com",
		DeepLinkSchema:    "https://vesseoapp.com/disbursement",
		SEP10ClientDomain: "vesseoapp.com",
		Assets:            assets.AllAssetsTestnet,
	},
	{
		Name:              "Via Wallet",
		Homepage:          "https://www.solvewithvia.com/wallet/",
		DeepLinkSchema:    "https://www.solvewithvia.com/wallet/",
		SEP10ClientDomain: "solvewithvia.com",
		Assets:            assets.AllAssetsTestnet,
	},
	{
		Name:        "User Managed Wallet",
		Assets:      assets.AllAssetsTestnet,
		UserManaged: true,
	},
	{
		Name:           "Embedded Wallet",
		DeepLinkSchema: "SELF",
		Homepage:       "https://stellar.org",
		Assets: []data.Asset{
			assets.XLMAsset,
			assets.USDCAssetTestnet,
			assets.EURCAssetTestnet,
		},
		Embedded: true,
	},
}
