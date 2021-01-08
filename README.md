# vcn  <img align="right" src="https://github.com/vchain-us/vcn/blob/master/docs/img/cn-color.eeadbabe.svg" width="160px"/>
> **_vChain CodeNotary_**

[![CircleCI](https://circleci.com/gh/vchain-us/vcn.svg?style=shield)](https://circleci.com/gh/vchain-us/vcn)
[![Go Report Card](https://goreportcard.com/badge/github.com/vchain-us/vcn?style=flat-square)](https://goreportcard.com/report/github.com/vchain-us/vcn)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/vchain-us/vcn)
[![Docker pulls](https://img.shields.io/docker/pulls/codenotary/vcn?style=flat-square)](https://hub.docker.com/r/codenotary/vcn)
[![Changelog](https://img.shields.io/badge/CHANGELOG-.md-blue?style=flat-square)](https://github.com/vchain-us/vcn/blob/master/CHANGELOG.md)
[![Release](https://img.shields.io/github/release/vchain-us/vcn.svg?style=flat-square)](https://github.com/vchain-us/vcn/releases/latest)
## The Trust and Integrity platform for the Cloud native environment
Give any digital asset a meaningful, globally-unique, immutable identity that is authentic, verifiable, traceable from anywhere.

<img align="right" src="https://github.com/vchain-us/vcn/blob/master/docs/img/codenotary_mascot.png" width="256px"/>
When using CodeNotary vcn in source code, release, deployment or at runtime, you allow a continuous trust verification that can be used to detect unusual or unwanted activity in your workload and act on it.
<br/>
Powered by CodeNotary's digital identity infrastructure, vcn lets you notarize all of your digital assets that add a trust level of your choice, custom attributes and a meaningful status without touching or appending anything (unlike digital certificates).
That allows change and revocation post-release without breaking any customer environment.
<br/>
Everything is done in a global, collaborative way to break the common silo solution architecture. Leveraging an immutable, always-on DLT platform allows you to avoid complex setup of Certificate authorities or digital certificates (that are unfit for DevOps anyway).

## Table of contents

- [Quick start](#quick-start)
- [DevSecOps in mind](#devsecops-in-mind)
- [What kind of behaviors can CodeNotary vcn detect](#what-kind-of-behaviors-can-codenotary-vcn-detect)
- [How it works](#how-it-works)
- [Installation](#installation)
- [Usage](#usage)
- [Integrations](#integrations)
- [Documentation](#documentation)
- [Testing](#testing)
- [Cross-compiling for various platforms](#cross-compiling-for-various-platforms)
- [CodeNotary Ledger Compliance](#codenotary-ledger-compliance)

- [License](#license)

## Quick start

1. **Download CodeNotary vcn.** There are releases for different platforms:

- [Download the latest release](https://github.com/vchain-us/vcn/releases/latest) and then read the [Usage](#usage) section below.
- We recommend storing `vcn` in your `PATH` - Linux example:
   ```bash
   cp vcn-v<version>-linux-amd64 /usr/local/bin/vcn
   ```

2. **Authenticate digital objects** You can use the command as a starting point.

   ```bash
   vcn authenticate <file|dir://directory|docker://dockerimage|git://gitdirectory>
   ```

3. [**Create your identity (free)**](https://dashboard.codenotary.io/auth/signup) You need an identity at CodeNotary to notarize objects yourself (btw, we're a data-minimum company and only ask for data that is really required)


4. **Notarize existing digital objects** Once you have an account you can start notarizing digital assets to give them an identity.

   ```bash
   vcn login
   vcn notarize <file|dir://directory|docker://dockerimage|git://gitdirectory>
   ```


## DevSecOps in mind
Codenotary vcn is a solution written by a devops-obsessed engineers for Devops engineers to bring better trust and security to the the CloudNative source to deployment process

## What kind of behaviors can CodeNotary vcn detect
vcn (and its extensions for Docker, Kubernetes, documents or CI/CD) can detect, authenticate and alert on any behavior that involves using unauthentic digital assets. vcn verification can be embedded anywhere and can be used to trigger alerts, updates or workflows.

vcn is so versatile, it can help detecting or acting on the following (but not limited to):
* Immutable tagging of source code, builds, and container images with version number, owner, timestamp, organization, trust level, and much more
* Simple and tamper-proof extraction of notarized tags like version number, owner, timestamp, organization, and trust level from any source code, build and container (based on the related image)
* Quickly discover and identify untrusted, revoked or obsolete libraries, builds, and containers in your application
* Detect the launch of an authorized or unknown container immediately
* Prevent untrusted or revoked containers from starting in production
* Verify the integrity and the publisher of all the data received over any channel

and more
* Enable application version checks and actions
* Buggy or rogue libraries can be traced by simple revoke or unsupport
* Revoke or unsupport your build or build version post-deployment (no complex certificate revocation that includes delivery of newly signed builds)
* Stop unwanted containers from being launched
* Make revocation part of the remediation process
* Use revocation without impairing customer environments
* Trace source code to build to deployment by integration into CI/CD or manual workflow
* Tag your applications for specific use cases (alpha, beta - non-commercial aso).

not just containers, also virtual machines -  [check out vCenter Connector, in case you're running VMware vSphere](https://github.com/openfaas-incubator/vcenter-connector)
* Newly created or existing virtual machines automatically get a unique identity that can be trusted or untrusted
* Prevent launch of untrusted VMs
* Stop or suspend running outdated or untrusted VMs
* Detect the cloning or export of VMs and alert


## How it works
![vcn How it works](https://raw.githubusercontent.com/vchain-us/vcn/master/docs/vcn_hiwb.png "How it works")

## Installation

### Download binary

It's easiest to download the latest version for your platform from the [release page](
https://github.com/vchain-us/vcn/releases).

Once downloaded, you can rename the binary to `vcn`, then run it from anywhere.
> For Linux and macOS you need to mark the file as executable: `chmod +x vcn`

### Homebrew / Linuxbrew

If you are on macOS and using [Homebrew](https://brew.sh/) (or on Linux and using [Linuxbrew](https://linuxbrew.sh/)), you can install `vcn` with the following:

```
brew tap vchain-us/brew
brew install vcn
```

### Build from Source

After having installed [golang](https://golang.org/doc/install) 1.12 or newer clone this
repository into your working directory.

Now, you can build `vcn` in the working directory by using `make vcn` and then run `./vcn`.

Alternatively, you can install `vcn` in your system simply by running `make install`. This will put the `vcn` executable into `GOBIN` which is
accessible throughout the system.

## Usage

Basically, `vcn` can notarize or authenticate any of the following kind of assets:

- a **file**
- an entire **directory** (by prefixing the directory path with `dir://`)
- a **git commit** (by prefixing the local git working directory path with `git://`)
- a **container image** (by using `docker://` or `podman://` followed by the name of an image present in the local registry of docker or podman, respectively)

> It's possible to provide a hash value directly by using the `--hash` flag.

> It's also possible to notarize assets using wildcard.
> With `--recursive` flag is possible to iterate over inner directories.
>```shell script
>./vcn n "*.md" --recursive
>```

For detailed **command line usage** see [docs/cmd/vcn.md](https://github.com/vchain-us/vcn/blob/master/docs/cmd/vcn.md) or just run `vcn help`.

### Local api server

It's possible to start a local API server. All commands are supported.
The notarization password can be submitted with the `x-notarization-password` header.
Examples:

```bash
curl --location --request GET '127.0.0.1:8080/inspect/e2b58ab102dbadb3b1fd5139c8d2a937dc622b1b0d0907075edea163fe2cd093' \
--header 'x-notarization-password: *********' \
--header 'Authorization: Basic ****' \
--header 'Content-Type: application/json' \
--data-raw '{
	"Kind":		"file",
	"Name":		"CONTRIBUTING.md",
	"Hash":		"e2b58ab102dbadb3b1fd5139c8d2a937dc622b1b0d0907075edea163fe2cd093",
	"Size":		1400,
	"ContentType":	"text/plain; charset=utf-8"
}'
```
### Notarization

Register an account with [codernotary.io](https://codenotary.io) first.

Then start with the `login` command. `vcn` will walk you through login and importing up your secret upon initial use.
```
vcn login
```

Once your secret is set you can notarize assets like in the following examples:

```
vcn notarize <file>
vcn notarize dir://<directory>
vcn notarize docker://<imageId>
vcn notarize podman://<imageId>
vcn notarize git://<path_to_git_repo>
vcn notarize --hash <hash>
```

By default all assets are notarized private, so not much information is disclosed about the asset. If you want to make that public and therefore, more trusted, please use the `--public` flag.

```
vcn notarize --public <asset>
```

Change the asset's status:

```
vcn unsupport <asset>
vcn untrust <asset>
```

Finally, to fetch all assets you've notarized:

```
vcn list
```

### Authentication

```
vcn authenticate <file>
vcn authenticate dir://<directory>
vcn authenticate docker://<imageId>
vcn authenticate podman://<imageId>
vcn authenticate git://<path_to_git_repo>
vcn authenticate --hash <hash>
```
> You can use `vcn authenticate` even without a [codernotary.io](https://codenotary.io) account.

To output results in `json` or `yaml` formats:
```
vcn authenticate --output=json <asset>
vcn authenticate --output=yaml <asset>
```
> Check out the [user guide](https://github.com/vchain-us/vcn/blob/master/docs/user-guide/formatted-output.md) for further details.


## Integrations

* [Github Action](https://github.com/marketplace/actions/verify-commit) - An action to verify the authenticity of your commits within your Github workflow
* [docker](https://github.com/vchain-us/vcn/blob/master/docs/user-guide/schemes/docker.md) - Out of the box support for notarizing and authenticating Docker images.
* [hub.docker.com/r/codenotary/vcn](https://hub.docker.com/r/codenotary/vcn) - The `vcn`'s DockerHub repository.
* [kube-notary](https://github.com/vchain-us/kube-notary) - A Kubernetes watchdog for verifying image trust with CodeNotary.
* [vcn-watchdog](https://github.com/vchain-us/vcn-watchdog) - Continuous authentication with CodeNotary for Docker.
* [jsvcn](https://github.com/vchain-us/jsvcn) - CodeNotary JavaScript Client.
* [jvcn](https://github.com/vchain-us/jvcn) - CodeNotary Java Bindings.
* [jvcn-maven-plugin](https://github.com/vchain-us/jvcn-maven-plugin) - Maven dependency authentication and enforcement.

## Documentation

* [Command line usage](https://github.com/vchain-us/vcn/blob/master/docs/cmd/vcn.md)
* [Configuration](https://github.com/vchain-us/vcn/blob/master/docs/user-guide/configuration.md)
* [Environments](https://github.com/vchain-us/vcn/blob/master/docs/user-guide/environments.md)
* [Formatted output (json/yaml)](https://github.com/vchain-us/vcn/blob/master/docs/user-guide/formatted-output.md)
* [Notarization explained](https://github.com/vchain-us/vcn/blob/master/docs/user-guide/notarization.md)

## Examples

#### Authenticate a Docker image automatically prior to running it

First, you’ll need to pull the image by using:

```
docker pull hello-world
```

Then use the below command to put in place an automatic safety check. It allows only verified images to run.

```
vcn authenticate docker://hello-world && docker run hello-world
```
If an image was not verified, it will not run and nothing will execute.


#### Authenticate multiple assets
You can authenticate multiple assets by piping other command outputs into `vcn`:
```
ls | xargs vcn authenticate
```
> The exit code will be `0` only if all the assets in you other command outputs are verified.

#### Authenticate by a specific signer
By adding `--signerID`, you can authenticate that your asset has been signed by a specific SignerID.
> A SignerID is the signer public address (represented as a 40 hex characters long string prefixed with `0x`).

```
vcn authenticate --signerID 0x8f2d1422aed72df1dba90cf9a924f2f3eb3ccd87 docker://hello-world
```

#### Authenticate by a list of signers

If an asset you or your organization wants to trust needs to be verified against a list of signers as a prerequisite, then use the `vcn authenticate` command and the following syntax:

- Add a `--signerID` flag in front of each SignerID you want to add
(eg. `--signerID 0x0...1 --signerID 0x0...2`)
- Or set the env var `VCN_SIGNERID` correctly by using a space to separate each SignerID (eg. `VCN_SIGNERID=0x0...1 0x0...2`)
> Be aware that using the `--signerID` flag will take precedence over `VCN_SIGNERID`.

The asset authentication will succeed only if the asset has been signed by at least one of the signers.

#### Authenticate using the asset's hash

If you want to authenticate an asset using only its hash, you can do so by using the command as shown below:

```
vcn authenticate --hash fce289e99eb9bca977dae136fbe2a82b6b7d4c372474c9235adc1741675f587e
```

#### Unsupport/untrust an asset you do not have anymore

In case you want to unsupport/untrust an asset of yours that you no longer have, you can do so using the asset hash(es) with the following steps below.

First, you’ll need to get the hash of the asset from your CodeNotary [dashboard](https://dashboard.codenotary.io/) or alternatively you can use the `vcn list` command. Then, in the CLI, use:

```
vcn untrust --hash <asset's hash>
# or
vcn unsupport --hash <asset's hash>
```

#### Notarization within automated environments

First, you’ll need to make `vcn` have access to the `${HOME}/.vcn` folder that holds your secret (the private key).
Then, set up your environment accordingly using the following commands:
```
export VCN_USER=<email>
export VCN_PASSWORD=<login password>
export VCN_NOTARIZATION_PASSWORD=<notarization password>
```
> It's possible to disable one time password requirement with:
>```bash
> export VCN_OTP_EMPTY=true
> ```
Once done, you can use `vcn` in your non-interactive environment using:

```
vcn login
vcn notarize <asset>
```
> Other commands like `untrust` and `unsupport` will also work.

## Testing
```
make test
```

## Cross-compiling for various platforms

The C libraries of [go-ethereum](https://github.com/ethereum/go-ethereum) make a more sophisticated cross-compilation
necessary.
The `make dist` target takes care of all steps by using [xgo](https://github.com/techknowlogick/xgo) and [docker](https://github.com/docker).


## CodeNotary Ledger Compliance

Vcn was extended in order to be compatible with [CodeNotary Ledger Compliance](https://codenotary.com/) .
Notarized assets informations are stored in a tamperproof ledger with cryptographic verification backed by [immudb](https://codenotary.com/technologies/immudb/), the immutable database.
Thanks to this `vcn` is faster and provides more powerful functionalities like the local data inclusion and consistency verification and an enhanced cli filter system.

### Obtain an api key
To provide access to Ledger Compliance a valid api key is required.
This api key is bound to a specific ledger and it's required during vcn login.
To obtain a valid key you need to get access to a licensed CodeNotary Ledger Compliance platform.

### Login

To login in Ledger Compliance provides --lc-port and --lc-host flag and submit api key when requested.
Once host port and api key are provided it's possible to omit them in following commands; it's also possible to provide them in other commands like `notarize`, `verify` or `inspect`.
```shell script
vcn login --lc-port 3324 --lc-host 127.0.0.1
```
> One time password (otp) is not mandatory

> To set up a secure connection (tls) with Ledger Compliance server it's possible to provide a certificate
>```shell script
>vcn login --lc-port 3324 --lc-host 127.0.0.1  --lc-cert mycert.pem
>```

### Commands
All commands reference didn't change.

### Inspect
Inspect is extended with the addition of new filter: `--last`, `--first`, `--start` and `--end`.
With `--last` and `--first` are returned the N first or last respectively.
```shell script
vcn inspect document.pdf --last 10
```
With `--start` and `--end` it's possible to use a time range filter
```shell script
vcn inspect document.pdf --start 2020/10/28-08:00:00 --end 2020/10/28-17:00:00
```
If no filters are provided only maximum 100 items are returned.

### Signer Identifier
It's possible to filter results with a single signer identifier
```shell script
vcn inspect document.pdf --signerID CygBE_zb8XnprkkO6ncIrbbwYoUq5T1zfyEF6DhqcAI=
```

### Local api server

Local API server is supported.
The `api key` can be submitted with the `x-notarization-lc-api-key` header.

Notarize example:
```bash
curl --location --request POST '127.0.0.1:8082/notarize' \
--header 'x-notarization-lc-api-key: oikfnlbjinhhclvjiotckgwfuyfjxntxmcau' \
--header 'Content-Type: application/json' \
--data-raw '{
"Kind":"file",
"Name":"CONTRIBUTING.md",
"Hash":"e2b58ab102dbadb3b1fd5139c8d2a937dc622b1b0d0907075edea163fe2cd093",
"Size":1400,
"ContentType":"text/plain; charset=utf-8"
}'
```
Authenticate example:
```bash
curl --location --request GET '127.0.0.1:8081/authenticate/e2b58ab102dbadb3b1fd5139c8d2a937dc622b1b0d0907075edea163fe2cd093' \
--header 'x-notarization-lc-api-key: oikfnlbjinhhclvjiotckgwfuyfjxntxmcau'
```
Inspect example:
```bash
curl --location --request GET '127.0.0.1:8082/authenticate/e2b58ab102dbadb3b1fd5139c8d2a937dc622b1b0d0907075edea163fe2cd093' \
--header 'x-notarization-lc-api-key: oikfnlbjinhhclvjiotckgwfuyfjxntxmcau'
```
Inspect with signerID example:
```bash
curl --location --request GET '127.0.0.1:8081/inspect/e2b58ab102dbadb3b1fd5139c8d2a937dc622b1b0d0907075edea163fe2cd093?signerid=yZtm26ZgmZr37NQ41TXbJ2jStMVWZhE-3cp4Wb7gKQo=' \
--header 'x-notarization-lc-api-key: oikfnlbjinhhclvjiotckgwfuyfjxntxmcau'
```
Untrust example:
```bash
curl --location --request POST '127.0.0.1:8081/untrust' \
--header 'x-notarization-lc-api-key: oikfnlbjinhhclvjiotckgwfuyfjxntxmcau' \
--header 'Content-Type: application/json' \
--data-raw '{
	"Kind":		"file",
	"Name":		"CONTRIBUTING.md",
	"Hash":		"e2b58ab102dbadb3b1fd5139c8d2a937dc622b1b0d0907075edea163fe2cd093",
	"Size":		1400,
	"ContentType":	"text/plain; charset=utf-8"
}'
```
## Generating smart contracts on linux

Clone https://github.com/ethereum/go-ethereum and compile `abigen` command in ./cmd/abigen

Download solc at https://github.com/ethereum/solidity/releases?after=v0.5.0 and chmod +x on it and copy in /usr/bin
solc version must be:
```bash
solc --version
solc, the solidity compiler commandline interface
Version: 0.4.24+commit.e67f0147.Linux.g++
```

Place at the root of the contracts repo and generate go files with:
```
abigen --sol organisations/contracts/OrganisationsRelay.sol --pkg blockchain --out organisationsrelay.go
abigen --sol assets/contracts/AssetsRelay.sol --pkg blockchain --out assetsrelay.go
```

## License

This software is released under [GPL3](https://www.gnu.org/licenses/gpl-3.0.en.html).
