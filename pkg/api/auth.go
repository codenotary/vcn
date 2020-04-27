/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 */

package api

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/sirupsen/logrus"
	"github.com/vchain-us/vcn/internal/errors"
	"github.com/vchain-us/vcn/pkg/meta"
)

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	defer r.Body.Close()

	if r.StatusCode == http.StatusOK {
		return response.Exists, nil
	}
	return false, fmt.Errorf("check publisher failed: %+v", restError)
}

func checkToken(token string) (success bool, err error) {
	restError := new(Error)
	response, err := newSling(token).
		Get(publisherEndpoint()+"/auth/check").
		Receive(nil, restError)
	if err != nil {
		return false, fmt.Errorf("check token failed: %s", err)
	}
	logger().WithFields(logrus.Fields{
		"response":  response,
		"err":       err,
		"restError": restError,
	}).Trace("checkToken")
	if response != nil {
		defer response.Body.Close()
		switch response.StatusCode {
		case http.StatusOK:
			return true, nil
		case http.StatusUnauthorized:
			fallthrough
		case http.StatusForbidden:
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

func authenticateUser(email, password string) (token string, err error) {
	response := new(tokenResponse)
	restError := new(Error)
	r, err := sling.New().
		Post(publisherEndpoint()+"/auth").
		BodyJSON(authRequest{Email: email, Password: password}).
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
	r.Body.Close()
	switch r.StatusCode {
	case http.StatusOK:
		return response.Token, nil
	case http.StatusBadRequest:
		return "", fmt.Errorf(errors.UnconfirmedEmail, email, meta.DashboardURL())
	case http.StatusUnauthorized:
		return "", fmt.Errorf("invalid password")
	}
	if restError.Error != "" {
		err = fmt.Errorf("%+v", restError)
	}
	return "", fmt.Errorf("authentication failed: %s", err)
}
