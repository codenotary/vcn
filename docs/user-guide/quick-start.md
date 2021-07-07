# Quick start

## Overview
`vcn` is the *Command Line Interface* for the CodeNotary platform. Basically, it can [notarize and authenticate](notarization.md) any of the following kind of assets:

- a **file**
- an entire **directory** (by prefixing the directory path with `dir://`)
- a **git commit** (by prefixing the local git working directory path with `git://`)
- a [**container image**](schemes/docker.md) (by using `docker://` or `podman://` followed by the name of an image present in the local registry of docker or podman, respectively)
- a **java BoM component** (by using `javacom://javacomponent` followed by the name of a maven jar or pom.xml file)

> It's also possible to provide a hash value directly by using the `--hash` flag.

## Install the CLI

The easiest way to get `vcn` is to download the latest version for your platform from the [release page](
https://github.com/vchain-us/vcn/releases).

Once downloaded, you can rename the binary to `vcn` and store it in your `PATH`, then run it from anywhere.
> For Linux and macOS you need to mark the file as executable: `chmod +x vcn`

## Authenticate an asset

Authentication is always free and can be performed by anyone, anywhere, at any time. You can use `vcn authenticate` even without a [codernotary.io](https://codenotary.io) account.

Examples:
```
vcn authenticate <file>
vcn authenticate dir://<directory>
vcn authenticate docker://<imageId>
vcn authenticate podman://<imageId>
vcn authenticate git://<path_to_git_repo>
vcn authenticate --hash <hash>
vcn authenticate javacom://<java maven jar or pom.xml>
```

To output results in `json` or `yaml` formats:
```
vcn authenticate --output=json <asset>
vcn authenticate --output=yaml <asset>
```

## Notarize an asset

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
vcn notarize javacom://<java maven jar or pom.xml>
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
