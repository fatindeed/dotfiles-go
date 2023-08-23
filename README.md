# dotfiles-go

## Supported Cipher Methods

- AES128_GCM
- AES256_GCM
- AES256_GCM_RAW
- AES128_CTR_HMAC_SHA256
- AES256_CTR_HMAC_SHA256
- CHACHA20_POLY1305
- XCHACHA20_POLY1305

## Supported Storage

-   Local Storage

    Example: `file://path/to/file`

-   HCP Vault Secrets

    Example: `hcp://app-name/secret-name`

    Environment Variables:

    - `HCP_CLIENT_ID`
    - `HCP_CLIENT_SECRET`
    - `HCP_ORGANIZATION_ID`
    - `HCP_PROJECT_ID`

-   S3

    Example: `s3://bucket/key`

    Environment Variables:

    - `DOTFILES_S3_ENDPOINT`
    - `DOTFILES_S3_REGION`
    - `DOTFILES_S3_ACCESS_KEY_ID`
    - `DOTFILES_S3_SECRET_ACCESS_KEY`

<!-- -   OnePassword

    Example: `op://vault/name/prop` -->

## Inspiration
- [chezmoi](https://github.com/twpayne/chezmoi)
- [yadm: Yet Another Dotfiles Manager](https://yadm.io/)
- [Homemaker](https://github.com/FooSoft/homemaker)
- [Tink](https://github.com/google/tink)