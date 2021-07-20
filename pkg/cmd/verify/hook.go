/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package verify

import (
	"fmt"
	"path/filepath"

	"github.com/vchain-us/vcn/pkg/store"

	"github.com/spf13/cobra"
	"github.com/vchain-us/vcn/pkg/bundle"

	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/extractor/dir"
)

type hook struct {
	a       api.Artifact
	rawDiff bool
}

func newHook(cmd *cobra.Command, a *api.Artifact) *hook {
	if a != nil {
		h := hook{
			a: a.Copy(),
		}
		h.rawDiff, _ = cmd.Flags().GetBool("raw-diff")
		dir.RemoveMetadata(a)
		return &h
	}
	return nil
}

func (h *hook) finalize(alertConfig *api.AlertConfig, output string) error {
	if h != nil && output == "" {
		manifest, path := dir.Metadata(h.a)
		if manifest != nil && path != "" {
			oldManifest, err := store.ReadManifest(h.a.Kind, path)
			if err != nil {
				oldManifest, err = bundle.ReadManifest(filepath.Join(path, bundle.ManifestFilename))
			}
			if err != nil {
				return nil // ignore missing manifest
			}
			// check old manifest integrity
			oldDigest, err := oldManifest.Digest()
			if err != nil {
				fmt.Printf("Diff is unavailable because '%s' is invalid.\n\n", bundle.ManifestFilename)
				return nil // ignore bad manifest
			}
			v, err := api.Verify(oldDigest.Encoded())
			if err != nil {
				return err
			}
			if v != nil && !v.Unknown() {
				var report string
				var equal bool
				var err error
				if h.rawDiff {
					report, equal, err = manifest.Diff(*oldManifest)
				} else {
					report, equal, err = manifest.DiffByPath(*oldManifest)
				}
				if err != nil {
					return err
				}
				if !equal {
					fmt.Printf("Diff since %s\n\n%s\n\n", v.Date(), report)
					if alertConfig != nil {
						alertConfig.Metadata["diff"] = report
					}
				}
			} else {
				fmt.Printf("Diff is unavailable because '%s' has been tampered.\n\n", bundle.ManifestFilename)
			}
		}
	}
	return nil
}

func (h *hook) lcFinalizeWithoutAlert(user *api.LcUser, output string, txId uint64) error {
	if h != nil && output == "" {
		manifest, path := dir.Metadata(h.a)
		if manifest != nil && path != "" {
			oldManifest, err := store.ReadManifest(h.a.Kind, path)
			if err != nil {
				oldManifest, err = bundle.ReadManifest(filepath.Join(path, bundle.ManifestFilename))
			}
			if err != nil {
				return nil // ignore missing manifest
			}
			// check old manifest integrity
			oldDigest, err := oldManifest.Digest()
			if err != nil {
				fmt.Printf("Diff is unavailable because '%s' is invalid.\n\n", bundle.ManifestFilename)
				return nil // ignore bad manifest
			}
			oldArtifact, _, err := user.LoadArtifact(oldDigest.Encoded(), "", "", txId, nil)

			if err != nil {
				if err == api.ErrNotFound {
					fmt.Printf("%s was not notarized", oldDigest.Encoded())
				} else {
					return err
				}
			}
			if oldArtifact != nil {
				var report string
				var equal bool
				var err error
				if h.rawDiff {
					report, equal, err = manifest.Diff(*oldManifest)
				} else {
					report, equal, err = manifest.DiffByPath(*oldManifest)
				}
				if err != nil {
					return err
				}
				if !equal {
					fmt.Printf("Diff since %s\n\n%s\n\n", oldArtifact.Timestamp, report)
				}
			} else {
				fmt.Printf("Diff is unavailable because '%s' has been tampered.\n\n", bundle.ManifestFilename)
			}
		}
	}
	return nil
}
