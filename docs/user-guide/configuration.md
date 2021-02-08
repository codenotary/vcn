# Configuration

By default, the `vcn` command line stores its config file (`config.json`) in a directory called `.vcn` within your `$HOME` directory.
> If the `STAGE` environment variable has been set, the default configuration directory can be different. See [environments](environments.md).

However, you can specify a different location for the config file via the `--vcnpath` command line option. For example:

```
vcn --vcnpath /path/to/your/config.json
```

<!-- The config file contains paths to keystore directories, and stores credentials of the current authenticated user.

`vcn` manages these files and directories and you should not modify them.
However, *you can modify* the config file to control where keys are stored. -->

## Config file

### Example of `config.json`

```
{
  "currentcontext": "example@example.net",
  "users": [
    {
      "email": "example@example.net",
      "token": "<authentication_bearer_token>"
      ]
    }
  ]
}
```

### Breakdown of `config.json`'s components

#### currentcontext

The property `currentcontext` holds the reference (user's email) to the current authenticated user.

#### users

The property `users` is an array of objects (one entry per user). Each object holds:

 - `email` the email address that identifies a specific user
 - `token` a bearer token used obtained by using `vcn login`
 <!-- - `keystore` path to the actual directory that store private keys -->

<!-- ### Storing secret keys

Secret keys are stored as encrypted JSON files according to the Web3 Secret Storage specification.
See https://github.com/ethereum/wiki/wiki/Web3-Secret-Storage-Definition for more information.

You can modify the `keystore` property according to your needs in order to store secret keys in a different location. -->
