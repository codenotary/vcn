package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

func ProvideKeystorePassword() (passphrase string, err error) {
	passphrase = os.Getenv("KEYSTORE_PASSWORD")
	if passphrase != "" {
		LOG.Trace("Keystore password provided (environment)")
		return passphrase, nil
	}
	fmt.Print("Keystore passphrase: ")
	passphraseBytes, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println(".")
	if err != nil {
		return "", nil
	}
	LOG.Trace("Keystore password provided (interactive)")
	return string(passphraseBytes), nil
}

func ProvideUsername() (user string, err error) {
	user = os.Getenv("VCN_USER")
	if user != "" {
		LOG.WithFields(logrus.Fields{
			"username": user,
		}).Trace("Username provided (environment)")
		return user, nil
	}
	fmt.Print("Email address: ")
	cnt, err := fmt.Scanln(&user)
	if err != nil {
		return "", err
	}
	if cnt <= 0 {
		return "", fmt.Errorf("username must not be empty")
	}
	user = strings.TrimSpace(user)
	LOG.WithFields(logrus.Fields{
		"username": user,
	}).Trace("Username provided (interactive)")
	return user, nil
}
