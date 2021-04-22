package nearapi

import "os"

// A Config for the NEAR network.
type Config struct {
	NetworkID string
	NodeURL   string
}

// GetConfig returns the NEAR network config depending on the setting of the
// environment variable NEAR_ENV.
func GetConfig() *Config {
	switch os.Getenv("NEAR_ENV") {
	case "production":
		fallthrough
	case "mainnet":
		return &Config{
			NetworkID: "mainnet",
			NodeURL:   "https://rpc.mainnet.near.org",
		}
	case "betanet":
		return &Config{
			NetworkID: "betanet",
			NodeURL:   "https://rpc.betanet.near.org",
		}
	case "local":
		return &Config{
			NetworkID: "local",
			NodeURL:   "http://localhost:3030",
		}
	case "development":
		fallthrough
	case "testnet":
		fallthrough
	default:
		return &Config{
			NetworkID: "default",
			NodeURL:   "https://rpc.testnet.near.org",
		}
	}
}
