/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package bom_java

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/prometheus/common/log"
	"github.com/vchain-us/vcn/pkg/bom_component"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const mvn_version = 3
const mvn_artid = 0
const mvn_pkg_name = 1

// JavaMavenPackage implements Package interface
type JavaMavenPackage struct {
	folder string
}

// New returns new JavaMavenPackage object
func New(path string) *JavaMavenPackage {
	f, err := getPOMFolder(path)
	if err != nil {
		return nil
	}

	return &JavaMavenPackage{folder: f}
}

func (p *JavaMavenPackage) Type() string {
	return "JavaMaven"
}

func (p *JavaMavenPackage) Close() {}

// Components returns list of java packages used during the build
func (p *JavaMavenPackage) Components() ([]bom_component.Component, error) {
	fname, err := exec.LookPath("mvn")
	if err != nil {
		return nil, fmt.Errorf("please install mvn tool follwing this link: https://maven.apache.org/install.html. Error reported: %w", err)
	}

	res := make([]bom_component.Component, 0)

	xmlDepFn, err := ioutil.TempFile(filepath.Dir(p.folder), "xml")
	if err != nil {
		return nil, err
	}
	defer xmlDepFn.Close()
	command := exec.Command(fname, "dependency:tree", "-DoutputType=graphml", "-f="+p.folder, "-DappendOutput=true", "-DoutputFile="+xmlDepFn.Name())
	_, err = command.Output()
	if err != nil {
		return nil, err
	}

	xmlDep, err := ioutil.ReadFile(xmlDepFn.Name())
	if err != nil {
		return nil, err
	}
	graph := GraphML{}
	err = xml.Unmarshal(xmlDep, &graph)
	if err != nil {
		return nil, err
	}
	for _, n := range graph.Graph.Nodes {

		fields := strings.Split(n.Data.ShapeNode.NodeLabel, ":")
		if fields[mvn_pkg_name] == "" {
			log.Errorf("unable to retrieve package name of component %s", n.Data.ShapeNode.NodeLabel)
			continue
		}
		if fields[mvn_version] == "" {
			log.Errorf("unable to retrieve package version of component %s", n.Data.ShapeNode.NodeLabel)
			continue
		}
		if fields[mvn_artid] == "" {
			log.Errorf("unable to retrieve package artifact ID of component %s", n.Data.ShapeNode.NodeLabel)
			continue
		}
		var comp bom_component.Component

		mavenSHA1url := strings.Join([]string{maven_repo_base_path, strings.Replace(fields[mvn_artid], ".", "/", -1), fields[mvn_pkg_name], fields[mvn_version], fields[mvn_pkg_name] + "-" + fields[mvn_version] + ".jar.sha1"}, "/")
		resp, err := http.Get(mavenSHA1url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		hash := string(body[0:40])
		if len(hash) < 40 {
			continue
		}
		comp.Hash = hash
		comp.Version = fields[mvn_version]

		comp.Name = fields[mvn_pkg_name]

		res = append(res, comp)
	}

	return res, nil
}

func getPOMFolder(path string) (string, error) {

	if strings.ToLower(filepath.Base(path)) == "pom.xml" {
		return path, nil
	}
	r, err := zip.OpenReader(path)

	if err != nil {
		return "", err
	}

	defer r.Close()
	for _, f := range r.File {
		if filepath.Base(f.Name) != "pom.xml" {
			continue
		}
		hash := sha256.Sum256([]byte(f.Name))
		tmpDirName := base64.RawURLEncoding.EncodeToString(hash[:])
		tempDir := filepath.Join(os.TempDir(), "vcn", tmpDirName)
		err := os.MkdirAll(tempDir, 755)
		if err != nil {
			return "", err
		}
		tmp, err := ioutil.TempFile(tempDir, "vcn_bom_pom.xml")
		if err != nil {
			return "", err
		}
		defer tmp.Close()
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()

		_, err = io.Copy(tmp, rc)
		if err != nil {
			return "", err
		}
		return tmp.Name(), nil
	}
	return "", errors.New("no pom is founded")
}
