package escherhelper

import (
	"github.com/emartech/escher-go"
)

type CredentialConfig struct {
	Host             string
	DisableBodyCheck bool
	AccessKeyID      string
	APISecret        string
	CredentialScope  string
	Date             string
}

func (config *CredentialConfig) GetCredentialScope() string {
	if config.CredentialScope == "" {
		credentialScope := "eu/suite/ems_request"

		return credentialScope
	}

	return config.CredentialScope
}

func (config *CredentialConfig) GetEscherConfig() escher.EscherConfig {
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
