# Updated

*Go program to update and push files periodically to a Git repository*

[![updated](https://github.com/qdm12/updated/raw/master/title.png)](https://hub.docker.com/r/qmcgaw/updated)

[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/updated.svg)](https://github.com/qdm12/updated/issues)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/updated.svg)](https://github.com/qdm12/updated/issues)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/updated.svg)](https://github.com/qdm12/updated/issues)

## Features

- Periodically builds
    - A list of unique malicious hostnames
    - A list of unique malicious IP addresses
    - A list of unique ads hostnames
    - A list of unique surveillance hostnames
    - InterNIC's named roots for DNS resolvers
    - Root anchors XML for DNS resolvers
    - Root keys to be used by Unbound
- Optionally upload the changes to a Git repository using an SSH key
- Based on the Scratch image with a total uncompressed size of 15.4MB
- Compatible with amd64 only although it is easily cross CPU compiled if needed

## Setup

### Using Docker (recommended)

1. For bind mounting, create a `files` directory with the right permissions:

    ```sh
    mkdir files
    chown 1000 files
    chmod 700 files
    ```

1. Use the following command:

    ```sh
    docker run -d -v $(pwd)/files:/files qmcgaw/updated
    ```

    You can also use [docker-compose.yml](https://github.com/qdm12/updated/blob/master/docker-compose.yml) with:

    ```sh
    docker-compose up -d
    ```

    To use with Git, you will also need to bind mount some files:
    - SSH key file at `/key` by default
    - SSH key passphrase optionally at `/passphrase`, if your SSH key is encrypted
    - SSH known hosts at `/known_hosts` by default, the default contains only Github key
    And set their ownership to user ID `1000` also.

1. Check logs with `docker logs updated` and update the image with `docker pull qmcgaw/updated`

### Environment variables

This Go program only reads parameters from environment variables for ease of use with Docker.

- Commonly used

    | Environment variable | Default | Possible values | Description |
    | --- | --- | --- | --- |
    | `OUTPUT_DIR` | `./files` | Any absolute or relative directory path | Directory where files are written to |
    | `PERIOD` | `24h` | Integer from `1` | Period in minutes between each run |
    | `RESOLVE_HOSTNAMES` | `no` | `yes` or `no` | Resolve hostnames found to obtain IP addresses |
    | `HTTP_TIMEOUT` | `3s` | *integer* from 1 | Default HTTP client timeout in milliseconds |
    | `LOG_LEVEL` | `info` | `debug`, `info`, `warning`, `error` | Logging level |
    | `TZ` | `America/Montreal` | *string* | Timezone |

- Git operation

    | Environment variable | Default | Possible values | Description |
    | --- | --- | --- | --- |
    | `GIT` | `no` | `yes` or `no` | Do Git operations or not |
    | `GIT_URL` | | SSH Git URL address | SSH URL to the remote repository, used only if `GIT=yes` |
    | `SSH_KEY` | `./key` | File path | File path containing your SSH key, only used if `GIT=yes` |
    | `SSH_KEY_PASSPHRASE` | | File path | Optional file path containing your SSH key passphrase, only used if `GIT=yes`. **Does not work with OpenSSH keys** |

- Checksums

    | Environment variable | Default | Possible values | Description |
    | --- | --- | --- | --- |
    | `NAMED_ROOT_MD5` | `9cdf6a12db2dc558e299601015c37948` | MD5 hexadecimal sum or `""` | Named root MD5 sum. Disables checking if empty |
    | `ROOT_ANCHORS_SHA256` | `45336725f9126db810a59896ae93819de743c416262f79c4444042c92e520770` | SHA256 hexadecimal sum or `""` | Root anchors SHA256 sum. Disables checking if empty |

- Extras

    | Environment variable | Default | Possible values | Description |
    | --- | --- | --- | --- |
    | `SHOUTRRR_SERVICES` | | One or more comma separated Shoutrrr URLs | Comma separated list of [Shoutrrr URLs](https://containrrr.dev/shoutrrr/services/overview/) |

### Using Go

1. Build the program

    ```sh
    go build cmd/updated/main.go -o updated
    ```

1. Depending on your system, change its permissions `chmod +x updated`
1. Run the program `./updated`

## Why

This container is used to periodically update files at [github.com/qdm12/files](https://github.com/qdm12/files) which are used by several other projects.

## TODOs

- [ ] Compress long repetitive files
- [ ] Unit tests (no time sorry)
- [ ] Version in json file with updated files
- [ ] Use lists from Blockada

## License

This repository is under an [MIT license](https://github.com/qdm12/updated/master/license)
