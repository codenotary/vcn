/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/dghubble/sling"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

type Wallet struct {
	Id      uint64 `json:"id"`
	Address string `json:"address"`
}

func CreateKeystore(password string) {
	if password == "" {
		var err error
		password, err = readPassword("Keystore passphrase:")
		if err != nil {
			log.Fatal(err)
		}
	}
	ks := keystore.NewKeyStore(WalletDirectory(), keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Public key:\t", account.Address.Hex())
	fmt.Println("Keystore:\t", WalletDirectory())
}

func HasKeystore() (bool, error) {
	files, err := ioutil.ReadDir(WalletDirectory())
	if err != nil {
		return false, err
	}
	return len(files) > 0, nil
}

func LoadPublicKeys() (addresses []string, err error) {
	authError := new(Error)
	wallets := new([]Wallet)
	token, err := LoadToken()
	if err != nil {
		log.Fatal(err)
	}
	r, err := sling.New().
		Add("Authorization", "Bearer "+token).
		Get(WalletEndpoint()).
		Receive(wallets, authError)
	if err != nil {
		log.Fatal(err)
	}
	if r.StatusCode != 200 {
		log.Fatalf("request failed: %s (%d)", authError.Message,
			authError.Status)
	}
	var result []string
	for _, wallet := range *wallets {
		result = append(result, wallet.Address)
	}
	return result, nil
}

func SyncKeys() {
	authError := new(Error)
	token, err := LoadToken()
	if err != nil {
		log.Fatal(err)
	}
	addresses, err := LoadPublicKeys()
	if err != nil {
		log.Fatal(err)
	}
	reader, err := firstFile(WalletDirectory())
	if err != nil {
		log.Fatal(err)
	}
	contents, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	var keyfile map[string]*json.RawMessage
	err = json.Unmarshal(contents, &keyfile)
	if err != nil {
		log.Fatal(err)
	}
	var localAddress string
	err = json.Unmarshal(*keyfile["address"], &localAddress)
	if err != nil {
		log.Fatal(err)
	}
	localAddress = "0x" + localAddress
	if contains(addresses, localAddress) {
		return
	}
	r, err := sling.New().
		Add("Authorization", "Bearer "+token).
		Post(WalletEndpoint()).
		BodyJSON(Wallet{Address: localAddress}).
		Receive(nil, authError)
	if err != nil {
		log.Fatal(err)
	}
	if r.StatusCode != 200 {
		log.Fatalf("request failed: %s (%d)", authError.Message,
			authError.Status)
	}
}
