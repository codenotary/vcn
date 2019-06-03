/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package store

import (
	"os"

	"github.com/spf13/viper"
)

// Keystore holds the path of a user's keystore
type Keystore struct {
	Path string `json:"path"`
}

// User holds user configuration
type User struct {
	Email     string      `json:"email"`
	Token     string      `json:"token"`
	Keystores []*Keystore `json:"keystores"`
}

type config struct {
	Users          []*User `json:"users"`
	CurrentContext string  `json:"currentContext"`
}

var cfg *config
var v = viper.New()

// Config returns the global config instance
func Config() *config {
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
	c := config{}
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
		return err
	}

	return nil
}

// SaveConfig stores the current configuration to file
func SaveConfig() error {
	// Setup config file
	setupConfigFile()

	// Ensure working dir
	if err := ensureDir(dir); err != nil {
		return err
	}

	v.Set("users", cfg.Users)
	v.Set("currentContext", cfg.CurrentContext)
	return v.WriteConfig()
}

// User returns an User from the global config matching the given email
func (c *config) User(email string) *User {
	if c == nil {
		return nil
	}

	for _, u := range c.Users {
		if u.Email == email {
			return u
		}
	}

	// fixme(leogr): what if email is empty?
	u := User{
		Email: email,
	}

	c.Users = append(c.Users, &u)
	return &u
}

// RemoveUser removes an user from config matching the given email, if not found return false
func (c *config) RemoveUser(email string) bool {
	if c == nil {
		return false
	}

	for i, u := range c.Users {
		if u.Email == email {
			c.Users = append(c.Users[i:], c.Users[i+1:]...)
			return true
		}
	}
	return false
}

// ClearContext unsets the current context (ie. logouts all users)
func (c *config) ClearContext() {
	if c == nil {
		return
	}
	for _, u := range c.Users {
		u.Token = ""
	}
	c.CurrentContext = ""
}
