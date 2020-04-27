/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package api

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/vchain-us/vcn/pkg/meta"
)

type verifyArtifactTrackingEventRequest struct {
	Client   string `json:"client"`
	Filename string `json:"filename"`
	Hash     string `json:"hash"`
	URL      string `json:"url"`
}

type signArtifactTrackingEventRequest struct {
	Filename string `json:"filename"`
	Hash     string `json:"hash"`
	Name     string `json:"name"`
	URL      string `json:"url"`
}

type publisherTrackingEventRequest struct {
	Name string `json:"name"`
}

func trackingEvent() string {
	return meta.APIEndpoint("tracking-event")
}

// TrackVerify is deprecated and will be removed
func TrackVerify(user *User, hash, filename string) (err error) {
	logger().WithFields(logrus.Fields{
		"hash":     hash,
		"filename": filename,
	}).Trace("TrackVerify")
	restError := new(Error)
	r, err := newSling(user.token()).
		Post(trackingEvent()+"/verify").
		BodyJSON(verifyArtifactTrackingEventRequest{
			Client:   meta.UserAgent(),
			Filename: filename,
			Hash:     hash,
		}).Receive(nil, restError)
	if err != nil {
		return err
	}
	r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("TrackVerify failed: %+v", restError)
	}
	return nil
}

// TrackPublisher is deprecated and will be removed
func TrackPublisher(user *User, event string) (err error) {
	logger().WithFields(logrus.Fields{
		"event": event,
	}).Trace("TrackPublisher")
	restError := new(Error)
	r, err := newSling(user.token()).
		Post(trackingEvent()+"/publisher").
		BodyJSON(publisherTrackingEventRequest{
			Name: event,
		}).Receive(nil, restError)
	if err != nil {
		return err
	}
	r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("TrackPublisher failed: %+v", restError)
	}
	return nil
}

// TrackSign is deprecated and will be removed
func TrackSign(user *User, hash, filename string, status meta.Status) (err error) {
	logger().WithFields(logrus.Fields{
		"hash":     hash,
		"filename": filename,
		"status":   status,
	}).Trace("TrackSign")
	restError := new(Error)
	r, err := newSling(user.token()).
		Post(trackingEvent()+"/sign").
		BodyJSON(signArtifactTrackingEventRequest{
			Name:     status.String(),
			Hash:     hash,
			Filename: filename,
		}).Receive(nil, restError)
	if err != nil {
		return err
	}
	r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("TrackSign failed: %+v", restError)
	}
	return nil
}
