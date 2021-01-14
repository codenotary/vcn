package login

import (
	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/cmd/internal/cli"
	"github.com/vchain-us/vcn/pkg/store"
)

// Execute the login action
func ExecuteLC(host, port, lcCert string, skipTlsVerify bool) error {
	apiKey, err := cli.ProvideLcApiKey()
	if err != nil {
		return err
	}
	if apiKey != "" {
		_, err = api.NewLcUser(apiKey, host, port, lcCert, skipTlsVerify)
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
