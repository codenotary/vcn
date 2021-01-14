package login

import (
	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/cmd/internal/cli"
	"github.com/vchain-us/vcn/pkg/store"
)

// Execute the login action
func ExecuteLC(host, port, lcCert string, skipTlsVerify, lcNoTls bool) error {
	apiKey, err := cli.ProvideLcApiKey()
	if err != nil {
		return err
	}
	if apiKey != "" {
		_, err = api.NewLcUser(apiKey, host, port, lcCert, skipTlsVerify, lcNoTls)
		if err != nil {
			return err
		}
		// Store the new config
		if err := store.SaveConfig(); err != nil {
			return err
		}
	}
	return nil
}
