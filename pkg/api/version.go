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
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/vchain-us/vcn/pkg/meta"
)

type version struct {
	Release string `json:"release"`
	Host    string `json:"host"`
	Message string `json:"message"`
	Latest  bool   `json:"latest"`
}

func getLatestVersion() (*version, error) {
	response := new(struct {
		Content []version `json:"content"`
	})
	restError := new(Error)

	url := meta.APIEndpoint("version/vcn") + "?sort=latest%2Cdesc"
	r, err := newSling("").
		Get(url).
		Receive(&response, restError)
	if err != nil {
		return nil, err
	}

	if r.StatusCode == 200 && len(response.Content) > 0 {
		return &response.Content[0], nil
	}

	logger().WithFields(logrus.Fields{
		"response":  response,
		"err":       err,
		"restError": restError,
	}).Trace("getLatestVersion")
	return nil, fmt.Errorf(
		"request failed: %s (%d)",
		restError.Message,
		restError.Status,
	)
}

// LatestCLIVersion returns the version string and message of the latest `vcn` CLI release.
func LatestCLIVersion() (string, string, error) {
	v, err := getLatestVersion()
	if err != nil || v == nil {
		return "", "", err
	}
	return v.Release, v.Message, nil
}

// LatestCLIDownloadURL returns the download URL of the latest `vcn` CLI release for the current architecture and OS.
func LatestCLIDownloadURL() (string, error) {
	response := new(struct {
		URL string `json:"url"`
	})
	restError := new(Error)
	url := fmt.Sprintf(
		"%s?arch=%s&os=%s",
		meta.APIEndpoint("version/vcn/latest-download-url"),
		runtime.GOARCH,
		runtime.GOOS,
	)
	if meta.StaticBuild() {
		url += "&static=true"
	}
	r, err := newSling("").
		Get(url).
		Receive(&response, restError)
	if err != nil {
		return "", err
	}

	if r.StatusCode != 200 {
		logger().WithFields(logrus.Fields{
			"response":  response,
			"err":       err,
			"restError": restError,
		}).Trace("LatestCLIDownloadURL")
		return "", fmt.Errorf(
			"request failed: %s (%d)",
			restError.Message,
			restError.Status,
		)
	}

	return response.URL, nil
}
