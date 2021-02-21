package escherhelper

import (
	"github.com/emartech/escher-go"
)

type CredentialsConfig struct {
	DisableBodyCheck bool
	AccessKeyID      string
	APISecret        string
	CredentialScope  string
	Date             string
}

func (config CredentialsConfig) GetCredentialScope() string {
	if config.CredentialScope == "" {
		return "eu/suite/ems_request"
	}

	return config.CredentialScope
}

func (config CredentialsConfig) GetEscherConfig() escher.EscherConfig {
	return escher.EscherConfig{
		VendorKey:       "Escher",
		AlgoPrefix:      "EMS",
		HashAlgo:        "SHA256",
		AuthHeaderName:  "X-Ems-Auth",
		DateHeaderName:  "X-Ems-Date",
		AccessKeyId:     config.AccessKeyID,
		ApiSecret:       config.APISecret,
		CredentialScope: config.GetCredentialScope(),
		Date:            config.Date,
	}
}
