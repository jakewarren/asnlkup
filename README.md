# asnlkup
[![Build Status](https://travis-ci.org/jakewarren/asnlkup.svg?branch=master)](https://travis-ci.org/jakewarren/asnlkup/)
[![GitHub release](http://img.shields.io/github/release/jakewarren/asnlkup.svg?style=flat-square)](https://github.com/jakewarren/asnlkup/releases])
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/jakewarren/asnlkup/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/jakewarren/asnlkup)](https://goreportcard.com/report/github.com/jakewarren/asnlkup)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=shields)](http://makeapullrequest.com)
> bulk lookup of AS info for IP addresses using IP2Location's ASN database

## Install
### Option 1: Binary

Download the latest release from [https://github.com/jakewarren/asnlkup/releases/latest](https://github.com/jakewarren/asnlkup/releases/latest)

### Option 2: From source

```
go get github.com/jakewarren/asnlkup
```

### Prerequisites

This program relies upon the IP2Location IP-ASN database https://lite.ip2location.com/database/ip-asn. Download the IPV4 CSV file from this page.

I recommend placing the database file in `/home/username/.cache/asnlkup/`.

## Example

```
❯ echo "8.8.8.8" | asnlkup 
IP      |ASN   |ISP
8.8.8.8 |15169 |Google Inc.
```

## Usage

`asnlkup` reads newline separated IP addresses from a file or STDIN.

```
❯ asnlkup -h
Usage: bulkiplkup [<flags>] [FILE]

Optional flags:

  -c, --csv=false: output in CSV format
  -d, --db="/home/jake/.cache/asnlkup/IP2LOCATION-LITE-ASN.CSV": db file name
  -h, --help=false: display help
  -j, --json=false: output in JSON format
  -o, --output="": output file name
```
## Changes

All notable changes to this project will be documented in the [changelog].

The format is based on [Keep a Changelog](http://keepachangelog.com/) and this project adheres to [Semantic Versioning](http://semver.org/).

## License

MIT © 2018 Jake Warren

[changelog]: https://github.com/jakewarren/asnlkup/blob/master/CHANGELOG.md
