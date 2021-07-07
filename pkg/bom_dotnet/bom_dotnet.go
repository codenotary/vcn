/*
 * Copyright (c) 2021 CodeNotary, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package bom_dotnet

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/vchain-us/vcn/pkg/bom_component"
)

// DotnetPackage implements Package interface
type DotnetPackage struct {
	path     string
	projects []string
}

type PackageFile struct {
	Targets map[string]Deps `json:"dependencies"`
}

type Deps map[string]Dep

type Dep struct {
	Type    string
	Version string `json:"resolved"`
	Hash    string `json:"contentHash"`
}

type mapKey struct {
	name    string
	version string
}

const (
	packagesFileName   = "packages.lock.json"
	fallbackDir        = "NuGetFallbackFolder"
	dotnet             = "dotnet"
	checksumFileSuffix = ".nupkg.sha512"
)

var restoreArgs = []string{"restore", "--use-lock-file"}
var localCacheArgs = []string{"nuget", "locals", "global-packages", "--list"}
var listSdkArgs = []string{"--list-sdks"}
var versionArgs = []string{"--version"}

var projFileExtensions = []string{".csproj", ".vbproj"}

// New returns new NugetPackage object if Nuget packages file is found in directory
func New(path string) *DotnetPackage {
	fi, err := os.Stat(path)
	if err != nil {
		return nil
	}
	if fi.IsDir() {
		var projFiles []string
		files, err := filepath.Glob(filepath.Join(path, "*.sln"))
		if err == nil && len(files) > 0 {
			for _, file := range files {
				projFiles = append(projFiles, getProjFromSln(file)...)
			}
		} else {
			for _, ext := range projFileExtensions {
				files, err = filepath.Glob(filepath.Join(path, "*"+ext))
				if err == nil && len(files) > 0 {
					projFiles = append(projFiles, files...)
				}
			}
		}

		if len(projFiles) > 0 {
			return &DotnetPackage{path: path, projects: projFiles}
		}
	} else {
		// if path is file, it should point to either Solution or Project file
		if strings.HasSuffix(path, ".sln") {
			projFiles := getProjFromSln(path)
			if len(projFiles) == 0 {
				return nil
			}
			return &DotnetPackage{path: filepath.Dir(path), projects: projFiles}
		}
		for _, ext := range projFileExtensions {
			if strings.HasSuffix(path, ext) {
				return &DotnetPackage{path: filepath.Dir(path), projects: []string{filepath.Base(path)}}
			}
		}
	}

	return nil
}

// Components returns all dependencies for all known projects
func (p *DotnetPackage) Components() ([]bom_component.Component, error) {
	// restore all dependencies
	absPath, err := filepath.Abs(p.path)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(dotnet, restoreArgs...)
	cmd.Dir = absPath
	_, err = cmd.Output()
	if err != nil {
		return nil, errors.New("cannot download dependencies - 'dotnet restore' failed")
	}

	// find local (user) and global package cache directories
	pkgCacheDirs := make([]string, 2)
	pkgCacheDirs[0], err = userCacheDir()
	if err != nil {
		return nil, err
	}
	pkgCacheDirs[1], err = fallbackCacheDir()
	if err != nil {
		return nil, err
	}

	res := make([]bom_component.Component, 0)
	seen := map[mapKey]struct{}{} // use common hashmap for all projects to avoid dupes
	for _, proj := range p.projects {
		pkgFileName := filepath.Join(absPath, filepath.Dir(proj), packagesFileName)

		projDeps, err := processPackageFile(pkgFileName, seen, pkgCacheDirs)
		if err != nil {
			return nil, err
		}
		res = append(res, projDeps...)
	}

	return res, nil
}

func (p *DotnetPackage) Type() string {
	return ".Net"
}

func (p *DotnetPackage) Path() string {
	return p.path
}

func (p *DotnetPackage) Close() {
	// do nothing, function needed to comply with Package interface
}

func getProjFromSln(slnFile string) []string {
	file, err := os.Open(slnFile)
	if err != nil {
		return nil
	}
	defer file.Close()

	var res []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "Project(") {
			continue
		}
		// check if project isn't a dummy one - it should point to project file
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		fileName := strings.Trim(fields[3], ",\"")
		for _, ext := range projFileExtensions {
			if strings.HasSuffix(fileName, ext) {
				res = append(res, fileName)
				break
			}
		}
	}

	return res
}

// process single package file
// it uses map, passed by caller, to avoid dupes. Map access isn't synced so this function isn't thread-safe
func processPackageFile(fileName string, seen map[mapKey]struct{}, pkgCacheDirs []string) ([]bom_component.Component, error) {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("cannot read Nuget dependencies file: %w", err)
	}

	var lockFile PackageFile
	err = json.Unmarshal(buf, &lockFile)
	if err != nil {
		return nil, fmt.Errorf("cannot parse Nuget dependencies file: %w", err)
	}

	res := make([]bom_component.Component, 0)
	for _, target := range lockFile.Targets {
		for name, details := range target {
			if details.Type == "Project" {
				continue // skip entry for the parent project
			}
			_, ok := seen[mapKey{name: name, version: details.Version}]
			if ok {
				continue // module already processed
			}
			seen[mapKey{name: name, version: details.Version}] = struct{}{} // mark as seen

			hash, err := getPackageHash(name, details.Version, pkgCacheDirs)
			if err != nil {
				return nil, err
			}
			res = append(res, bom_component.Component{
				Name:     name,
				Version:  details.Version,
				Hash:     hash,
				HashType: bom_component.HashSHA512})
		}
	}

	return res, nil
}

func userCacheDir() (string, error) {
	cmd := exec.Command(dotnet, localCacheArgs...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	fields := strings.Fields(string(out))
	return fields[1], nil
}

func fallbackCacheDir() (string, error) {
	cmd := exec.Command(dotnet, versionArgs...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	ver := strings.TrimSpace(string(out))

	cmd = exec.Command(dotnet, listSdkArgs...)
	out, err = cmd.Output()
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(out))
	for scanner.Scan() {
		fields := strings.SplitN(scanner.Text(), " ", 2)
		if len(fields) < 2 || fields[0] != string(ver) {
			continue
		}
		return filepath.Join(strings.Trim(fields[1], "[]"), fallbackDir), nil
	}

	return "", errors.New("cannot find Nuget fallback directory")
}

func getPackageHash(name, version string, pkgCacheDirs []string) (string, error) {
	lowerName := strings.ToLower(name)
	for _, cache := range pkgCacheDirs {
		hashBase64, err := ioutil.ReadFile(filepath.Join(
			cache,
			lowerName,
			version,
			lowerName+"."+version+checksumFileSuffix))
		if err != nil {
			continue
		}
		hashRaw, err := base64.StdEncoding.DecodeString(string(hashBase64))
		if err != nil {
			return "", fmt.Errorf("cannot decode hash for package '%s'", name)
		}
		return hex.EncodeToString(hashRaw), nil
	}

	return "", fmt.Errorf("cannot find package '%s' in package cache", name)
}
