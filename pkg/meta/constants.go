/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package meta

import (
	"log"

	"github.com/fatih/color"
)

type Level int64
type Status int64
type Visibility int64

const (
	LevelDisabled         Level = -1
	LevelUnknown          Level = 0
	LevelEmailVerified    Level = 1
	LevelSocialVerified   Level = 2
	LevelIdVerified       Level = 3
	LevelLocationVerified Level = 4
	LevelVchain           Level = 99
)

const (
	StatusTrusted     Status = 0
	StatusUntrusted   Status = 1
	StatusUnknown     Status = 2
	StatusUnsupported Status = 3
)

const (
	VisibilityPublic  Visibility = 0
	VisibilityPrivate Visibility = 1
)

const (
	VcnLoginEvent  string = "VCN_LOGIN"
	VcnSignEvent   string = "VCN_SIGN"
	VcnVerifyEvent string = "VCN_VERIFY"
	// todo(leogr): backend needs to be update prior to fix event strings
	KeyCreatedEvent  string = "KEYSTORE_CREATED"
	KeyUploadedEvent string = "KEYSTORE_UPLOADED"
)

const (
	VcnUserEnv          string = "VCN_USER"
	VcnPasswordEnv      string = "VCN_PASSWORD"
	KeyStorePasswordEnv string = "KEYSTORE_PASSWORD"
)

func VcnClientName() (name string) {
	return "VCN:" + Version()
}

func LevelName(level Level) (name string) {
	switch level {
	case LevelDisabled:
		return "DISABLED"
	case LevelUnknown:
		return "0 - UNKNOWN"
	case LevelEmailVerified:
		return "1 - EMAIL_VERIFIED"
	case LevelSocialVerified:
		return "2 - SOCIAL_VERIFIED"
	case LevelIdVerified:
		return "3 - ID_VERIFIED"
	case LevelLocationVerified:
		return "4 - LOCATION_VERIFIED"
	case LevelVchain:
		return "99 - VCHAIN"
	default:
		log.Fatal("unsupported level", name)
		return ""
	}
}

func StatusName(status Status) (name string) {
	switch status {
	case StatusTrusted:
		return "TRUSTED"
	case StatusUntrusted:
		return "UNTRUSTED"
	case StatusUnknown:
		return "UNKNOWN"
	case StatusUnsupported:
		return "UNSUPPORTED"
	default:
		log.Fatal("unsupported status", name)
		return ""
	}
}

func StatusNameStyled(status Status) string {
	c, s := StatusColor(status)
	return color.New(c, s).Sprintf(StatusName(status))
}

func VisibilityName(visibility Visibility) (name string) {
	switch visibility {
	case VisibilityPublic:
		return "PUBLIC"
	case VisibilityPrivate:
		return "PRIVATE"
	default:
		log.Fatal("unsupported visibility", name)
		return ""
	}
}

func VisibilityFromString(text string) (visibility Visibility) {
	switch text {
	case "0":
		return VisibilityPublic
	case "1":
		return VisibilityPrivate
	case "PUBLIC":
		return VisibilityPublic
	case "PRIVATE":
		return VisibilityPrivate
	default:
		return VisibilityPrivate
	}
}

func VisibilityForFlag(public bool) (visibility Visibility) {
	if public {
		return VisibilityPublic
	} else {
		return VisibilityPrivate
	}
}
