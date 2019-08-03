# Updated

*Docker container to update and push files periodically to Github*

[![updated](https://github.com/qdm12/updated/raw/master/title.png)](https://hub.docker.com/r/qmcgaw/updated)

[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/updated.svg)](https://github.com/qdm12/updated/issues)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/updated.svg)](https://github.com/qdm12/updated/issues)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/updated.svg)](https://github.com/qdm12/updated/issues)

| Image size | RAM usage | CPU usage |
| --- | --- | --- |
| 65.2MB | 1MB | Very low |

It is based on:

- [Alpine 3.10](https://alpinelinux.org) with Git, OpenSSH, Wget, Bind tools, Perl XMLPathg

It updates the following every 24 hours:

- Malicous IPs
- Malicious Hostnames
- NSA Hostnames
- DNS Root Anchors
- DNS Named root

## Setup

1. Download the *Dockerfile* and the *entrypoint.sh* script

    ```bash
    wget https://raw.githubusercontent.com/qdm12/updated/master/Dockerfile
    wget https://raw.githubusercontent.com/qdm12/updated/master/entrypoint.sh
    ```
  
1. Place your SSH Github private key as file `key` in your working directory
1. Build the image

    ```bash
    docker build -t qmcgaw/updated .
    ```

1. Run the container
  
    ```bash
    docker run -d qmcgaw/updated
    ```

### Environment variables

- `VERBOSE` set to `1` or `0`, defaults to `1`

## TODOs

- [ ] Healthcheck
- [ ] More checks i.e. bad wget
- [ ] Use lists from Blockada

## License

This repository is under an [MIT license](https://github.com/qdm12/updated/master/license)
