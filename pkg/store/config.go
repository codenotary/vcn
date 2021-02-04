/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package store

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	configSchemaVer uint = 3
)

// User holds user's configuration.
type User struct {
	Email    string `json:"email,omitempty"`
	Token    string `json:"token,omitempty"`
	KeyStore string `json:"keystore,omitempty"`
	LcApiKey string `json:"lcApiKey,omitempty"`
	LcCert   string `json:"lcCert,omitempty"`
}

// ConfigRoot holds root fields of the configuration file.
type ConfigRoot struct {
	SchemaVersion  uint           `json:"schemaVersion"`
	Users          []*User        `json:"users"`
	CurrentContext CurrentContext `json:"currentContext"`
}

type CurrentContext struct {
	Email           string `json:"email,omitempty"`
	LcApiKey        string `json:"lcApiKey,omitempty"`
	LcHost          string `json:"LcHost,omitempty"`
	LcPort          string `json:"LcPort,omitempty"`
	LcCert          string `json:"LcCert,omitempty"`
	LcSkipTlsVerify bool   `json:"LcSkipTlsVerify,omitempty"`
	LcNoTls         bool   `json:"LcNoTls,omitempty"`
}

func (cc *CurrentContext) Clear() {
	cc.Email = ""
	cc.LcApiKey = ""
	cc.LcHost = ""
	cc.LcPort = ""
	cc.LcCert = ""
	cc.LcSkipTlsVerify = false
	cc.LcNoTls = false
}

var cfg *ConfigRoot
var v = viper.New()

// Config returns the global config instance
func Config() *ConfigRoot {
	return cfg
}

func setupConfigFile() string {
	cfgFile := ConfigFile()
	v.SetConfigFile(cfgFile)
	v.SetConfigPermissions(FilePerm)
	return cfgFile
}

// LoadConfig loads the global configuration from file
func LoadConfig() error {

	// Make default config
	c := ConfigRoot{
		SchemaVersion: configSchemaVer,
	}
	cfg = &c

	// Setup config file
	cfgFile := setupConfigFile()

	// Ensure working dir
	if err := ensureDir(dir); err != nil {
		return err
	}

	// Create default file if it does not exist yet
	if ConfigFile() == defaultConfigFilepath() {
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			sErr := SaveConfig()
			return sErr
		}
	}

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(&c); err != nil {
		oldFormat := ConfigRootV2{
			SchemaVersion: 2,
		}
		if err := v.Unmarshal(&oldFormat); err != nil {
			return errors.New("unable to parse config file")
		}
		fmt.Println("Upgrading config to new format. Old sessions will expire")
		c.Users = []*User{}
		c.CurrentContext.Email = oldFormat.CurrentContext
		c.SchemaVersion = 3
	}

	ul := 0
	akr := false
	for _, u := range c.Users {
		if u.LcApiKey != "" {
			c.Users = append(c.Users[:ul], c.Users[ul+1:]...)
			ul--
			akr = true
		}
		ul++
	}
	if akr {
		fmt.Println("Configuration file cleaned.")
	}
	return SaveConfig()
}

// SaveConfig stores the current configuration to file
func SaveConfig() error {
	// Setup config file
	setupConfigFile()

	// Ensure working dir
	if err := ensureDir(dir); err != nil {
		return err
	}

	cfg.SchemaVersion = configSchemaVer
	v.Set("users", cfg.Users)
	v.Set("currentContext", cfg.CurrentContext)
	v.Set("schemaVersion", cfg.SchemaVersion)
	return v.WriteConfig()
}

// User returns an User from the global config matching the given email.
// User returns nil when an empty email is given or c is nil.
func (c *ConfigRoot) UserByMail(email string) *User {
	defer func() {
		if cfg != nil {
			cfg.CurrentContext.Clear()
			cfg.CurrentContext.Email = email
		}
	}()
	if c == nil || email == "" {
		return nil
	}

	for _, u := range c.Users {
		if u.Email == email {
			u.LcApiKey = ""
			return u
		}
	}

	u := User{
		Email: email,
	}

	c.Users = append(c.Users, &u)
	return &u
}

// User returns an User from the global config matching the given email.
// User returns nil when an empty email is given or c is nil.
func (c *ConfigRoot) UserByLcApiKey(lcApiKey string) (u *User) {
	defer func() {
		cfg.CurrentContext.Clear()
		cfg.CurrentContext.LcApiKey = lcApiKey
	}()
	if c == nil || lcApiKey == "" {
		return nil
	}

	for _, u := range c.Users {
		if u.LcApiKey == lcApiKey {
			return u
		}
	}

	u = &User{
		LcApiKey: lcApiKey,
	}

	c.Users = append(c.Users, u)
	return u
}

// User returns an User from the global config matching the given email.
// User returns nil when an empty email is given or c is nil.
func (c *ConfigRoot) NewLcUser(lcApiKey, host, port, lcCert string, lcSkipTlsVerify, lcNoTls bool) (u *User) {
	defer func() {
		cfg.CurrentContext.Clear()
		cfg.CurrentContext.LcApiKey = lcApiKey
		cfg.CurrentContext.LcHost = host
		cfg.CurrentContext.LcPort = port
		cfg.CurrentContext.LcCert = lcCert
		cfg.CurrentContext.LcSkipTlsVerify = lcSkipTlsVerify
		cfg.CurrentContext.LcNoTls = lcNoTls
	}()
	if c == nil || lcApiKey == "" {
		return nil
	}

	u = &User{
		LcApiKey: lcApiKey,
	}

	return u
}

// RemoveUser removes an user from config matching the given email, if not found return false
func (c *ConfigRoot) RemoveUserByMail(email string) bool {
	if c == nil {
		return false
	}

	for i, u := range c.Users {
		if u.Email == email {
			c.Users = append(c.Users[:i], c.Users[i+1:]...)
			return true
		}
	}
	return false
}

// ClearContext clean up all auth token for all users and set an empty context.
func (c *ConfigRoot) ClearContext() {
	if c == nil {
		return
	}
	for _, u := range c.Users {
		u.Token = ""
	}
	c.CurrentContext = CurrentContext{}
}
