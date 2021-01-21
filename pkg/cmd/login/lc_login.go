package login

import (
	"context"
	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/cmd/internal/cli"
	"github.com/vchain-us/vcn/pkg/meta"
	"github.com/vchain-us/vcn/pkg/store"
	"google.golang.org/grpc/metadata"
)

// Execute the login action
func ExecuteLC(host, port, lcCert string, skipTlsVerify, lcNoTls bool) error {
	apiKey, err := cli.ProvideLcApiKey()
	if err != nil {
		return err
	}
	if apiKey != "" {
		u, err := api.NewLcUser(apiKey, host, port, lcCert, skipTlsVerify, lcNoTls)
		if err != nil {
			return err
		}
		if u != nil {
			err = u.Client.Connect()
			if err != nil {
				return err
			}
			md := metadata.Pairs(meta.VcnLCPluginTypeHeaderName, meta.VcnLCPluginTypeHeaderValue)
			ctx := metadata.NewOutgoingContext(context.Background(), md)
			_, err = u.Client.Health(ctx)
			if err != nil {
				return err
			}
			// Store the new config
			if err := store.SaveConfig(); err != nil {
				return err
			}
		}
	}
	// shouldn't happen
	return nil
}
