/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 */

package api

import (
	"fmt"

	"github.com/dghubble/sling"
	"github.com/sirupsen/logrus"
	"github.com/vchain-us/vcn/internal/errors"
	"github.com/vchain-us/vcn/pkg/meta"
)

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Otp      string `json:"otp"`
}

type tokenResponse struct {
	Token string `token:"token"`
}

type publisherExistsResponse struct {
	Exists bool `json:"exists"`
}

type publisherExistsParams struct {
	Email string `url:"email"`
}

func publisherEndpoint() string {
	return meta.APIEndpoint("publisher")
}

func checkUserExists(email string) (success bool, err error) {
	response := new(publisherExistsResponse)
	restError := new(Error)
	r, err := sling.New().
		Get(publisherEndpoint()+"/exists").
		QueryStruct(&publisherExistsParams{Email: email}).
		Receive(&response, restError)
	logger().WithFields(logrus.Fields{
		"response":  response,
		"err":       err,
		"restError": restError,
	}).Trace("checkUserExists")
	if err != nil {
		return false, err
	}
	if r.StatusCode == 200 {
		return response.Exists, nil
	}
	return false, fmt.Errorf("check publisher failed: %+v", restError)
}

func checkToken(token string) (success bool, err error) {
	restError := new(Error)
	response, err := newSling(token).
		Get(publisherEndpoint()+"/auth/check").
		Receive(nil, restError)
	logger().WithFields(logrus.Fields{
		"response":  response,
		"err":       err,
		"restError": restError,
	}).Trace("checkToken")
	if response != nil {
		switch response.StatusCode {
		case 200:
			return true, nil
		case 401:
			fallthrough
		case 403:
			fallthrough
		case 419:
			return false, nil
		}
	}
	if restError.Error != "" {
		err = fmt.Errorf("%+v", restError)
	}
	return false, fmt.Errorf("check token failed: %s", err)
}

func authenticateUser(email string, password string, otp string) (token string, err error) {
	response := new(tokenResponse)
	restError := new(Error)
	r, err := sling.New().
		Post(publisherEndpoint()+"/auth").
		BodyJSON(authRequest{Email: email, Password: password, Otp: otp}).
		Receive(response, restError)
	logger().WithFields(logrus.Fields{
		"email":     email,
		"response":  response,
		"err":       err,
		"restError": restError,
	}).Trace("authenticateUser")
	if err != nil {
		return "", err
	}
	switch r.StatusCode {
	case 200:
		return response.Token, nil
	case 400:
		return "", fmt.Errorf(errors.UnconfirmedEmail, email, meta.DashboardURL())
	case 401:
		return "", fmt.Errorf("invalid password")
	}
	if restError.Error != "" {
		err = fmt.Errorf("%+v", restError)
	}
	return "", fmt.Errorf("authentication failed: %s", err)
}
