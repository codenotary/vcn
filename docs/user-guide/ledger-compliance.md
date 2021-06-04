# Codenotary Immutable Ledger

`vcn` tool is able to work with a standalone [CNIL] instance.

In such scenario, `vcn` is using an API key for authenticatation.
This key is connected to a specific ledger
and can be created in ledger's settings page
or through the API of the CNIL instance.

Once the API key is created, `vcn` can authenticate to the CNIL instance:

```sh
# Store the API key in the `VCN_LC_API_KEY` env variable
read -s VCN_LC_API_KEY

# Login to the CNIL server
vcn login --lc-port 443 --lc-host "<cnil instance host or ip>" --lc-cert "<path to cnil's certificate file>"

# Notarize / verify files
vcn notarize "<resource to be notarized>"
vcn authenticate "<resource to be verified>"
```

Note that either the `VCN_LC_API_KEY` env variable must be set or the `--lc-api-key` option specified
to perform any `vcn` operation using a standalone CNIL instance.

[CNIL]: https://www.codenotary.com/products/immutable-ledger
