/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package docker

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vchain-us/vcn/pkg/uri"
)

func TestDocker(t *testing.T) {
	_, err := exec.Command("docker", "pull", "hello-world").Output()
	if err != nil {
		t.Skip("docker not available")
	}

	u, _ := uri.Parse("docker://hello-world")
	a, err := Artifact(u)
	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.Equal(t, "docker://hello-world:latest", a.Name)
	assert.Regexp(t, "[0-9a-f]{64}", a.Hash)
	assert.NotZero(t, a.Size)
}

func TestInferVer(t *testing.T) {
	testCases := map[string]string{
		"golang:1.12-stretch": "1.12-stretch",
		"golang:latest":       "",
	}

	for tag, ver := range testCases {
		i := image{
			RepoTags: []string{tag},
		}
		assert.Equal(
			t,
			ver,
			i.inferVer(),
			"wrong version for %s", tag,
		)
	}
}
