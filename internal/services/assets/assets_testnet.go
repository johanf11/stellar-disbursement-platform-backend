package assets

import "github.com/stellar/stellar-disbursement-platform-backend/internal/data"

var AllAssetsTestnet = []data.Asset{
	XLMAsset,
	USDCAssetTestnet,
	HTGCAssetTestnet, // TODO: replace with mainnet issuer on launch (see HTGCAssetIssuerTestnet below)
}

// USDC

const USDCAssetIssuerTestnet = "GBBD47IF6LWK7P7MDEVSCWR7DPUWV3NY3DTQEVFL4NAT4AQH3ZLLFLA5"

var USDCAssetTestnet = data.Asset{
	Code:   USDCAssetCode,
	Issuer: USDCAssetIssuerTestnet,
}

// HTGC — Theo Haitian Gourde Coin (testnet demo issuer; replace with BVI SPV issuer on mainnet launch)
const HTGCAssetIssuerTestnet = "GDSRYZWTLQLBECKCL4TV7ZRGBZGBMSPD4V47B7Y7JSQVDJRSEXQTFCQT"

var HTGCAssetTestnet = data.Asset{
	Code:   "HTGC",
	Issuer: HTGCAssetIssuerTestnet,
}

// EURC

const EURCAssetIssuerTestnet = "GB3Q6QDZYTHWT7E5PVS3W7FUT5GVAFC5KSZFFLPU25GO7VTC3NM2ZTVO"

var EURCAssetTestnet = data.Asset{
	Code:   EURCAssetCode,
	Issuer: EURCAssetIssuerTestnet,
}
