# Signatures

*TODO* Asset signing explanation

## Statuses

Code | Status | Color | Description | Error message | Explanation
------------ | ------------- | ------------- | ------------ | ------------- | -------------
0 | **TRUSTED** | *green* | The asset was signed. | *none* | The signature on the blockchain indicates that the signer trusts that asset.
2 | **UNKNOWN** | *yellow* | The asset was not signed. | *hash* was not signed *[by <key/list of keys/org>]* | There's no signature on the blockchain.
1 | **UNTRUSTED** | *red* | The asset is untrusted. | *hash* was is untrusted *[by <key/list of keys/org>]* | The signature on the blockchain indicates that the signer DOES NOT trust that asset.
3 | **UNSUPPORTED** | *red* | The asset is unsupported. | *hash* was is unsupported *[by <key/list of keys/org>]* | The signature on the blockchain indicates that the signer DOES NOT trust that asset because it is not supported anymore (eg. deprecated).

## Levels

*TODO* Explain levels

