/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package meta

import (
	"fmt"
	"log"
	"runtime"

	"github.com/fatih/color"
)

// Level is the type for all possible signature levels
type Level int64

// Status is the type for all possible asset statuses
type Status int64

// Visibility is the type for all visibility values
type Visibility int64

// Allowed Level values
const (
	LevelDisabled         Level = -1
	LevelUnknown          Level = 0
	LevelEmailVerified    Level = 1
	LevelSocialVerified   Level = 2
	LevelIDVerified       Level = 3
	LevelLocationVerified Level = 4
	LevelVchain           Level = 99
)

// Allowed Status values
const (
	StatusTrusted     Status = 0
	StatusUntrusted   Status = 1
	StatusUnknown     Status = 2
	StatusUnsupported Status = 3
)

// Allowed Visibility values
const (
	VisibilityPublic  Visibility = 0
	VisibilityPrivate Visibility = 1
)

// Event tracking related consts
const (
	VcnLoginEvent  string = "VCN_LOGIN"
	VcnSignEvent   string = "VCN_SIGN"
	VcnVerifyEvent string = "VCN_VERIFY"
)

// vcn environment variable names
const (
	VcnUserEnv              string = "VCN_USER"
	VcnPasswordEnv          string = "VCN_PASSWORD"
	VcnNotarizationPassword string = "VCN_NOTARIZATION_PASSWORD"
    VcnWriteIgnoreFile      string = "VCN_WRITE_VCNIGNORE"
)

// UserAgent returns the vcn's User-Agent string
func UserAgent() string {
	// Syntax reference: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/User-Agent#Syntax
	return fmt.Sprintf("vcn/%s (%s; %s)", Version(), runtime.GOOS, runtime.GOARCH)
}

// LevelName returns the name of the given level as string
func LevelName(level Level) string {
	switch level {
	case LevelDisabled:
		return "DISABLED"
	case LevelUnknown:
		return "0 - UNKNOWN"
	case LevelEmailVerified:
		return "1 - EMAIL_VERIFIED"
	case LevelSocialVerified:
		return "2 - SOCIAL_VERIFIED"
	case LevelIDVerified:
		return "3 - ID_VERIFIED"
	case LevelLocationVerified:
		return "4 - LOCATION_VERIFIED"
	case LevelVchain:
		return "99 - VCHAIN"
	default:
		log.Fatal("unsupported level", level)
		return ""
	}
}

// StatusName returns the name of the given status as string
func StatusName(status Status) string {
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
		log.Fatal("unsupported status", status)
		return ""
	}
}

// StatusNameStyled returns the colorized name of the given status as string
func StatusNameStyled(status Status) string {
	c, s := StatusColor(status)
	return color.New(c, s).Sprintf(StatusName(status))
}

// VisibilityName returns the name of the given visibility as string
func VisibilityName(visibility Visibility) string {
	switch visibility {
	case VisibilityPublic:
		return "PUBLIC"
	case VisibilityPrivate:
		return "PRIVATE"
	default:
		log.Fatal("unsupported visibility", visibility)
		return ""
	}
}

// VisibilityForFlag returns VisibilityPublic if public is true, otherwise VisibilityPrivate
func VisibilityForFlag(public bool) Visibility {
	if public {
		return VisibilityPublic
	}
	return VisibilityPrivate
}
