/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package api

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	sdk "github.com/vchain-us/ledger-compliance-go/grpcclient"
	"github.com/vchain-us/vcn/pkg/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"io/ioutil"
	"strconv"
	"time"
)

func NewLcClientByContext(context store.CurrentContext) (*sdk.LcClient, error) {
	return NewLcClient(context.LcApiKey, context.LcHost, context.LcPort, context.LcCert, context.LcSkipTlsVerify, context.LcNoTls)
}

func NewLcClient(lcApiKey, host, port, lcCertPath string, skipTlsVerify, noTls bool) (*sdk.LcClient, error) {
	if skipTlsVerify && noTls {
		return nil, errors.New("illegal parameters submitted: lc-skip-tls-verify and lc-no-tls arguments are both provided")
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		return nil, errors.New("ledger compliance port is invalid")
	}
	defaultOptions := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                20 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	}

	currentOptions := []grpc.DialOption{}
	currentOptions = append(currentOptions, defaultOptions...)
	if !skipTlsVerify {
		if lcCertPath != "" {
			tlsCredentials, err := loadTLSCertificate(lcCertPath)
			if err != nil {
				return nil, fmt.Errorf("cannot load TLS credentials: %s", err)
			}
			currentOptions = append(currentOptions, grpc.WithTransportCredentials(tlsCredentials))
		} else {
			// automatic loading of local CA in os
			currentOptions = append(currentOptions, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
		}
	} else {
		currentOptions = append(currentOptions, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	}

	if noTls {
		currentOptions = []grpc.DialOption{grpc.WithInsecure()}
		currentOptions = append(currentOptions, defaultOptions...)
	}

	return sdk.NewLcClient(sdk.ApiKey(lcApiKey), sdk.Host(host), sdk.Port(p), sdk.Dir(store.CurrentConfigFilePath()), sdk.DialOptions(currentOptions)), nil
}

func loadTLSCertificate(certPath string) (credentials.TransportCredentials, error) {
	cert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(cert) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}
	config := &tls.Config{
		RootCAs: certPool,
	}
	return credentials.NewTLS(config), nil
}
