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

	"github.com/vchain-us/vcn/pkg/meta"
)

type alert struct {
	ArtifactHash     string   `json:"artifactHash,omitempty"`
	ArtifactMetaHash string   `json:"artifactMetaHash,omitempty"`
	Email            string   `json:"email,omitempty"`
	Enabled          bool     `json:"enabled"`
	Metadata         Metadata `json:"metadata,omitempty"`
	Name             string   `json:"name"`
	UUID             string   `json:"uuid,omitempty"`
}

// AlertResponse holds alert values returned by the platform.
type AlertResponse struct {
	ArtifactHash     string   `json:"artifactHash,omitempty"`
	ArtifactMetaHash string   `json:"artifactMetaHash,omitempty"`
	Email            string   `json:"email,omitempty"`
	Enabled          bool     `json:"enabled"`
	Metadata         Metadata `json:"metadata,omitempty"`
	Name             string   `json:"name"`
	UnAcknowledged   bool     `json:"unAcknowledgedNotification"`
	UUID             string   `json:"uuid,omitempty"`
}

// AlertConfig represents a platform alert configuration.
type AlertConfig struct {
	AlertUUID string   `json:"alertUuid" yaml:"alertUUID"`
	Metadata  Metadata `json:"metadata,omitempty" yaml:"metadata"`
}

// CreateAlert creates a platform alert and returns its UUID.
func (u *User) CreateAlert(name, email string, a *Artifact, v BlockchainVerification, m Metadata) (alertConfig *AlertConfig, err error) {
	restError := new(Error)
	alertResponse := &alert{}
	r, err := newSling(u.token()).
		Post(meta.APIEndpoint("alert")).
		BodyJSON(alert{
			ArtifactHash:     a.Hash,
			ArtifactMetaHash: v.MetaHash(),
			Email:            email,
			Enabled:          true,
			Metadata:         m,
			Name:             name,
		}).Receive(alertResponse, restError)
	if err != nil {
		return
	}
	defer r.Body.Close()

	switch r.StatusCode {
	case http.StatusOK:
		alertConfig = &AlertConfig{
			AlertUUID: alertResponse.UUID,
			Metadata:  m,
		}
	case http.StatusBadRequest:
		err = fmt.Errorf("%s", restError.Description)
	case http.StatusRequestEntityTooLarge:
		err = fmt.Errorf("%s is not allowed, only email addresses using the same domain are allowed", email)
	default:
		err = fmt.Errorf("alert API request failed: %s (%d)", restError.Message,
			restError.Status)
	}
	return
}

// ModifyAlert modifies the settings of an already existing alert.
// func (u *User) ModifyAlert(config *AlertConfig, enabled bool) error {

// 	if config == nil {
// 		return fmt.Errorf("alert config cannot be nil")
// 	}

// 	restError := new(Error)
// 	alertResponse := &alert{}
// 	alertRequest := alert{
// 		Enabled:  enabled,
// 		Metadata: config.Metadata,
// 	}

// 	r, err := newSling(u.token()).
// 		Patch(meta.APIEndpoint("alert")+"?uuid="+config.AlertUUID).
// 		BodyJSON(alertRequest).Receive(alertResponse, restError)

// 	if err != nil {
// 		return err
// 	}

// 	switch r.StatusCode {
// 	case 200:
// 		return nil
// 	case 403:
// 		return fmt.Errorf("illegal alert access: %s", restError.Message)
// 	case 404:
// 		return fmt.Errorf(`no such alert found matching "%s"`, config.AlertUUID)
// 	default:
// 		return fmt.Errorf("alert API request failed: %s (%d)", restError.Message,
// 			restError.Status)
// 	}
// }

// GetAlert returns an AlertResponse for a given alert uuid.
func (u *User) GetAlert(uuid string) (*AlertResponse, error) {
	restError := new(Error)
	response := &AlertResponse{}
	r, err := newSling(u.token()).
		Get(meta.APIEndpoint("alert/"+uuid)).
		Receive(&response, restError)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	switch r.StatusCode {
	case http.StatusOK:
		return response, nil
	case http.StatusForbidden:
		return nil, fmt.Errorf("illegal alert access: %s", restError.Message)
	case http.StatusNotFound:
		return nil, fmt.Errorf(`no such alert found matching "%s"`, uuid)
	default:
		return nil, fmt.Errorf("alert API request failed: %s (%d)", restError.Message,
			restError.Status)
	}
}

func (u *User) alertMessage(config AlertConfig, what string) (err error) {
	restError := new(Error)
	r, err := newSling(u.token()).
		Post(meta.APIEndpoint("alert/"+what)).
		BodyJSON(config).Receive(nil, restError)
	if err != nil {
		return
	}
	defer r.Body.Close()

	switch r.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusForbidden:
		return fmt.Errorf("illegal alert access: %s", restError.Message)
	case http.StatusNotFound:
		return fmt.Errorf(`no such alert found matching "%s"`, config.AlertUUID)
	case http.StatusPreconditionFailed:
		return fmt.Errorf(`notification already triggered for alert "%s"`, config.AlertUUID)
	default:
		return fmt.Errorf("alert API request failed: %s (%d)", restError.Message,
			restError.Status)
	}
}

// PingAlert sends a ping for the given alert _config_.
// Once the first ping goes through, the platform starts a server-side watcher and will trigger a notification
// after some amount of time if no further pings for the alert are received.
func (u *User) PingAlert(config AlertConfig) error {
	return u.alertMessage(config, "ping")
}

// TriggerAlert triggers a notification immediately for the given alert _config_.
func (u *User) TriggerAlert(config AlertConfig) error {
	return u.alertMessage(config, "notify")
}
