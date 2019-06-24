/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */
package api

import (
	"github.com/vchain-us/vcn/pkg/store"
)

type User struct {
	cfg *store.User
}

func NewUser(email string) *User {
	return &User{
		cfg: store.Config().User(email),
	}
}

func (u User) Email() string {
	if u.cfg != nil {
		return u.cfg.Email
	}
	return ""
}

func (u User) Keystores() []*store.Keystore {
	return u.cfg.Keystores
}

func (u User) DefaultKeystore() (*store.Keystore, error) {
	return u.cfg.DefaultKeystore()
}

func (u User) DefaultKey() string {
	if u.cfg != nil {
		return u.cfg.LastPubKey()
	}
	return ""
}

func (u User) HasKey() bool {
	if u.cfg != nil {
		return u.cfg.HasPubKey()
	}
	return false
}

func (u User) Keys() []string {
	if u.cfg != nil {
		return u.cfg.PubKeys()
	}
	return nil
}

func (u *User) Authenticate(password string) (err error) {
	if u == nil || u.Email() == "" {
		return makeFatal("user not initialized", nil)
	}

	token, err := authenticateUser(u.Email(), password)
	if err != nil {
		return err
	}

	u.cfg.Token = token
	return nil
}

func (u *User) ClearAuth() {
	if u != nil && u.cfg != nil {
		u.cfg.Token = ""
	}
}

func (u *User) token() string {
	if u != nil && u.cfg != nil {
		return u.cfg.Token
	}
	return ""
}

func (u User) IsAuthenticated() (bool, error) {
	if u.cfg == nil || u.cfg.Token == "" {
		return false, nil
	}

	return checkToken(u.cfg.Token)
}

func (u User) IsExist() (bool, error) {
	email := u.Email()
	if email != "" {
		return checkUserExists(email)
	}
	return false, nil
}
