// +build integration

/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package sign_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/extractor"
	"github.com/vchain-us/vcn/pkg/extractor/file"
	"github.com/vchain-us/vcn/pkg/meta"
	"github.com/vchain-us/vcn/pkg/store"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"
	"testing"
)

func TestParallelNotarization(t *testing.T) {

	userEmail := os.Getenv(meta.VcnUserEnv)
	if userEmail == "" {
		log.Fatalf("to run this test please set %s environment variable", meta.VcnUserEnv)
	}
	userPw := os.Getenv(meta.VcnPasswordEnv)
	if userPw == "" {
		log.Fatalf("to run this test please set %s environment variable", meta.VcnPasswordEnv)
	}

	extractor.Register(file.Scheme, file.Artifact)
	user := &api.User{}

	su := &store.User{Email: userEmail}

	user.UserByCfg(su)

	err := user.Authenticate(userPw, "")
	assert.NoError(t, err)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		a, _ := extractor.Extract([]string{"file://../../../CONTRIBUTING.md"})
		_, err := psign(user, a, userPw)
		assert.NoError(t, err)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		a, _ := extractor.Extract([]string{"file://../../../README.md"})
		_, err := psign(user, a, userPw)
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		a, _ := extractor.Extract([]string{"file://../../../LICENSE"})
		_, err := psign(user, a, userPw)
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		a, _ := extractor.Extract([]string{"file://../../../go.mod"})
		_, err := psign(user, a, userPw)
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		a, _ := extractor.Extract([]string{"file://../../../Makefile"})
		_, err := psign(user, a, userPw)
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		a, _ := extractor.Extract([]string{"file://../../../go.mod"})
		_, err := psign(user, a, userPw)
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		a, _ := extractor.Extract([]string{"file://../../../go.sum"})
		_, err := psign(user, a, userPw)
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		a, _ := extractor.Extract([]string{"file://../../../Dockerfile"})
		_, err := psign(user, a, userPw)
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		a, _ := extractor.Extract([]string{"file://../../../Dockerfile.docker"})
		_, err := psign(user, a, userPw)
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		a, _ := extractor.Extract([]string{"file://../../../CONTRIBUTING.md"})
		_, err := psign(user, a, userPw)
		assert.NoError(t, err)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		a, _ := extractor.Extract([]string{"file://../../../README.md"})
		_, err := psign(user, a, userPw)
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		a, _ := extractor.Extract([]string{"file://../../../LICENSE"})
		_, err := psign(user, a, userPw)
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		a, _ := extractor.Extract([]string{"file://../../../go.mod"})
		_, err := psign(user, a, userPw)
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		a, _ := extractor.Extract([]string{"file://../../../Makefile"})
		_, err := psign(user, a, userPw)
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		a, _ := extractor.Extract([]string{"file://../../../go.mod"})
		_, err := psign(user, a, userPw)
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		a, _ := extractor.Extract([]string{"file://../../../go.sum"})
		_, err := psign(user, a, userPw)
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		a, _ := extractor.Extract([]string{"file://../../../Dockerfile"})
		_, err := psign(user, a, userPw)
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		a, _ := extractor.Extract([]string{"file://../../../Dockerfile.docker"})
		_, err := psign(user, a, userPw)
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("done")
}

func psign(user *api.User, a []*api.Artifact, userPw string) (*api.BlockchainVerification, error) {
	keyin, _, _, err := user.Secret()
	if err != nil {
		log.Fatal(err)
	}
	opts := []api.SignOption{
		api.SignWithKey(keyin, userPw),
		api.SignWithStatus(meta.StatusTrusted),
	}

	v, err := user.Sign(*a[0], opts...)
	fmt.Printf("%s\n", v.Date())
	return v, err
}

func TestParallelNotarizationCLI(t *testing.T) {

	if !isCommandAvailable("vcn") {
		log.Fatalf("please install vcn in your $PATH")
	}
	userEmail := os.Getenv(meta.VcnUserEnv)
	if userEmail == "" {
		log.Fatalf("to run this test please set %s environment variable", meta.VcnUserEnv)
	}
	userPw := os.Getenv(meta.VcnPasswordEnv)
	if userPw == "" {
		log.Fatalf("to run this test please set %s environment variable", meta.VcnPasswordEnv)
	}

	file1, err := ioutil.TempFile("", "_vcn_test_")
	file1.Write([]byte(`1`))
	file1.Close()
	if err != nil {
		log.Fatal(err)
	}
	file2, err := ioutil.TempFile("", "_vcn_test_")
	file2.Write([]byte(`2`))
	file2.Close()
	if err != nil {
		log.Fatal(err)
	}
	file3, err := ioutil.TempFile("", "_vcn_test_")
	file3.Write([]byte(`3`))
	file3.Close()
	if err != nil {
		log.Fatal(err)
	}
	file4, err := ioutil.TempFile("", "_vcn_test_")
	file4.Write([]byte(`4`))
	file4.Close()
	if err != nil {
		log.Fatal(err)
	}
	file5, err := ioutil.TempFile("", "_vcn_test_")
	file5.Write([]byte(`5`))
	file5.Close()
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file1.Name())
	defer os.Remove(file2.Name())
	defer os.Remove(file3.Name())
	defer os.Remove(file4.Name())
	defer os.Remove(file5.Name())

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := cliSign(fmt.Sprintf("file://%s", file1.Name()))
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := cliSign(fmt.Sprintf("file://%s", file2.Name()))
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := cliSign(fmt.Sprintf("file://%s", file3.Name()))
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := cliSign(fmt.Sprintf("file://%s", file4.Name()))
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := cliSign(fmt.Sprintf("file://%s", file5.Name()))
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("done")
}

func cliSign(fn string) error {
	vcn := exec.Command("vcn", "n", fn)

	stdoutStderr, err := vcn.CombinedOutput()

	fmt.Printf("%s\n", stdoutStderr)

	return err
}

func isCommandAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
