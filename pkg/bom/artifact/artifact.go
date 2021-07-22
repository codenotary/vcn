/*
 * Copyright (c) 2021 CodeNotary, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package artifact

import (
	"errors"
	"fmt"

	"github.com/vchain-us/vcn/pkg/api"
)

// Artifact is a result of build process.
// It is a language- and/or environment-specific interface which finds dependencies
type Artifact interface {
	Path() string
	Type() string
	Dependencies() ([]Dependency, error)
}

type GenericArtifact struct {
	Deps []Dependency
}

// AuthDependencies checks if all dependencies have sufficient trust level, optionally notarizing unsupported ones.
// If trustlevel is Trusted and threshold > 0, up to threshold (in %) number of dependencies are allowed to
// have Unsupported status
func AuthDependencies(lcUser *api.LcUser, a Artifact, signerID string, trustLevel TrustLevel, autoNotarize bool, threshold uint) error {
	deps, err := a.Dependencies()
	if err != nil {
		return err
	}
	if len(deps) == 0 {
		return nil // nothing to do
	}
	if trustLevel < MinTrustLevel || trustLevel > MaxTrustLevel {
		return errors.New("invalid trust level")
	}
	// doesn't matter if threshold is > 100, it works the same as 100

	fmt.Printf("Authenticating dependencies...\n")

	unsupportedCount := 0
	messages := make([]string, 0)
	failed := false
	for i := range deps { // use index i to iterate because Authenticate() mutates the object
		err = deps[i].Authenticate(lcUser, signerID)
		if err != nil {
			return err
		}
		if deps[i].TrustLevel < trustLevel {
			messages = append(messages, fmt.Sprintf("Dependency %s trust level is %s", deps[i].Name, TrustLevelName(deps[i].TrustLevel)))
			if deps[i].TrustLevel == Unsupported {
				unsupportedCount++
			} else {
				failed = true
			}
		}
	}

	if threshold < 100 && unsupportedCount > len(deps)*int(threshold)/100 {
		failed = true
	}

	if failed {
		fmt.Printf("Some dependencies have insufficient trust level:\n")
		for _, m := range messages {
			fmt.Println(m)
		}
		return errors.New("some dependencies have insufficient trust level")
	}

	// notarize the components with Unsupported trust level
	if autoNotarize && unsupportedCount > 0 {
		fmt.Printf("Notarizing dependencies...\n")
		for i := range deps { // use index i to iterate because Notarize() mutates the object
			if deps[i].TrustLevel == Unsupported {
				err = deps[i].Notarize(lcUser, a.Type())
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
