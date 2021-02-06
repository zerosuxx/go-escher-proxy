package escherhelper

import "github.com/emartech/escher-go"

type CredentialConfig struct {
	Host            string
	AccessKeyId     string
	ApiSecret       string
	CredentialScope string
}

func (config *CredentialConfig) GetCredentialScope() *string {
	if config.CredentialScope == "" {
		credentialScope := "eu/suite/ems_request"

		return &credentialScope
	}

	return &config.CredentialScope
}

func (config *CredentialConfig) GetEscherConfig() escher.EscherConfig {
	return escher.EscherConfig{
		VendorKey:       "Escher",
		AlgoPrefix:      "EMS",
		HashAlgo:        "SHA256",
		AuthHeaderName:  "X-Ems-Auth",
		DateHeaderName:  "X-Ems-Date",
		AccessKeyId:     config.AccessKeyId,
		ApiSecret:       config.ApiSecret,
		CredentialScope: *config.GetCredentialScope(),
	}
}
