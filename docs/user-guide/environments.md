# Environments

By default `vcn` will put the config file and secrets within the a directory called `.vcn` within your [home directory](https://en.wikipedia.org/wiki/Home_directory) (e.g. `$HOME/.vcn` or `%USERPROFILE%\.vcn` on Windows).

However, `vcn` can work with distinct environments (eg. for testing purpose).

The following environments are supported by setting the `STAGE` environment var:

Stage | Directory | Note
------------ | ------------- | -------------
`STAGE=PRODUCTION` | `.vcn` | *default* 
`STAGE=STAGING` | `.vcn.staging` |
`STAGE=TEST` | `.vcn.test` | *`VCN_TEST_DASHBOARD`, `VCN_TEST_NET`, `VCN_TEST_CONTRACT`, `VCN_TEST_API` must be set accordingly to your test environment*


## Other environment variables

Name | Description | Example 
------------ | ------------- | -------------
`VCN_USER`, `VCN_PASSWORD` | Credentials for non-interactive user login | `VCN_USER=example@example.net VCN_PASSWORD=<your_password> vcn login`
`VCN_SIGNERID` | For `vcn authenticate` acts as a list of SignerID(s) (separated by space) to authenticate against | `VCN_SIGNERID="0x0...0 0x0...1" vcn authenticate <asset>` or `VCN_SIGNERID="0x0...0 <asset>` 
`VCN_ORG` | Organization's ID to authenticate against | `VCN_ORG="vchain.us" vcn authenticate <asset>`
`VCN_NOTARIZATION_PASSWORD` | Notarization password for non-interactive notarization | `VCN_NOTARIZATION_PASSWORD=<your_notarization_passphrase> vcn notarize <asset>`
`LOG_LEVEL` | Logging verbosity. Accepted values: `TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC`  | `LOG_LEVEL=TRACE vcn login` 
`HTTP_PROXY` | HTTP Proxy configuration | `HTTP_PROXY=http://localhost:3128 vcn authenticate <asset>`
 VCN_WRITE_VCNIGNORE | Write the default .vcnignore while interacting with a directory if it doesn't exist (checks the variable is set, not the value) | VCN_WRITE_VCIGNORE=1
