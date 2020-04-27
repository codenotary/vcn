/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package bundle

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/google/go-cmp/cmp"
	digest "github.com/opencontainers/go-digest"
)

// Diff returns a human-readable report as string of the raw differences between m and x.
//
// Do not depend on this output being stable.
func (m Manifest) Diff(x Manifest) (report string, equal bool, err error) {
	var r diffReporter
	equal = cmp.Equal(x, m, cmp.Reporter(&r))
	report = r.String()
	return
}

type diffReporter struct {
	path  cmp.Path
	lines []string
}

func (r *diffReporter) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

func (r *diffReporter) Report(rs cmp.Result) {
	if !rs.Equal() {
		vx, vy := r.path.Last().Values()
		var line string
		switch {
		case !vx.IsValid():
			line = fmt.Sprintf("%#v:\n\t+: %+v\n", r.path, vy)
		case !vy.IsValid():
			line = fmt.Sprintf("%#v:\n\t-: %+v\n", r.path, vx)
		default:
			line = fmt.Sprintf("%#v:\n\t-: %+v\n\t+: %+v\n", r.path, vx, vy)
		}
		r.lines = append(r.lines, line)
	}
}

func (r *diffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *diffReporter) String() string {
	return strings.Join(r.lines, "\n")
}

type modDiff struct {
	path string
	from Descriptor
	to   Descriptor
}

type pathDiff struct {
	desc Descriptor
	from string
	to   string
}

type itemDiff struct {
	path string
	desc Descriptor
}

// DiffByPath returns a human-readable report as string containing
// additions, modifications, renamings, deletions of x.Items relative to m.Items
// listed by path.
//
// Do not depend on this output being stable.
func (m Manifest) DiffByPath(x Manifest) (report string, equal bool, err error) {
	adds := make([]itemDiff, 0)
	mods := make([]modDiff, 0)
	rens := make([]pathDiff, 0)
	dels := make([]itemDiff, 0)

	mByPath := make(map[string]Descriptor)
	xByPath := make(map[string]Descriptor)
	newPaths := make(map[digest.Digest][]string)

	for _, d := range x.Items {
		for _, path := range d.Paths {
			xByPath[path] = d
		}
	}
	for _, d := range m.Items {
		for _, path := range d.Paths {
			mByPath[path] = d
			if _, ok := xByPath[path]; !ok {
				newPaths[d.Digest] = append(newPaths[d.Digest], path)
			}
		}
	}

	for path, xd := range xByPath {
		// try by path
		if md, ok := mByPath[path]; ok {
			if md.Digest != xd.Digest {
				// modified
				mods = append(mods, modDiff{
					path: path,
					from: xd,
					to:   md,
				})
			}
			// else:
			// same content, so no diff
		} else { // try by digest
			byDig := xd.Digest
			if mPaths, ok := newPaths[byDig]; ok && len(mPaths) > 0 {
				// renamed
				newPath := mPaths[0]
				newPaths[byDig] = mPaths[1:]
				rens = append(rens, pathDiff{
					desc: xd,
					from: path,
					to:   newPath,
				})
				delete(mByPath, newPath)
			} else {
				// deleted
				dels = append(dels, itemDiff{
					path: path,
					desc: xd,
				})
			}
		}
		delete(mByPath, path)
	}

	// finally, arrange new items
	for path, d := range mByPath {
		adds = append(adds, itemDiff{
			path: path,
			desc: d,
		})
	}

	equal = len(adds) == 0 && len(mods) == 0 && len(rens) == 0 && len(dels) == 0
	if equal {
		return // empty diff, no need to format lines
	}

	lines := []string{}

	sprintf := func(format string, a ...interface{}) {
		lines = append(lines, fmt.Sprintf(format, a...))
	}

	for _, d := range adds {
		sprintf(
			"\tnew item:   %s (%s)\n\t            + %s",
			d.path,
			humanize.Bytes(d.desc.Size),
			d.desc.Digest.String(),
		)
	}
	for _, d := range mods {
		sprintf(
			"\tmodified:   %s (%s -> %s)\n\t            - %s\n\t            + %s",
			d.path,
			humanize.Bytes(d.from.Size),
			humanize.Bytes(d.to.Size),
			d.from.Digest.String(),
			d.to.Digest.String(),
		)
	}
	for _, d := range rens {
		sprintf(
			"\trenamed:    %s -> %s (%s)\n\t            = %s",
			d.from,
			d.to,
			humanize.Bytes(d.desc.Size),
			d.desc.Digest.String(),
		)
	}
	for _, d := range dels {
		sprintf(
			"\tdeleted:    %s (%s)\n\t            - %s",
			d.path,
			humanize.Bytes(d.desc.Size),
			d.desc.Digest.String(),
		)
	}

	report = strings.Join(lines, "\n\n") + "\n"
	return
}
