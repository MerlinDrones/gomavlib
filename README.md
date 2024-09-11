# gomavlib

<!--
[![Test](https://github.com/bluenviron/gomavlib/workflows/test/badge.svg)](https://github.com/bluenviron/gomavlib/actions?query=workflow:test)
[![Lint](https://github.com/bluenviron/gomavlib/workflows/lint/badge.svg)](https://github.com/bluenviron/gomavlib/actions?query=workflow:lint)
[![Dialects](https://github.com/bluenviron/gomavlib/workflows/dialects/badge.svg)](https://github.com/bluenviron/gomavlib/actions?query=workflow:dialects)
[![Go Report Card](https://goreportcard.com/badge/github.com/bluenviron/gomavlib)](https://goreportcard.com/report/github.com/bluenviron/gomavlib)
[![CodeCov](https://codecov.io/gh/bluenviron/gomavlib/branch/main/graph/badge.svg)](https://app.codecov.io/gh/bluenviron/gomavlib/branch/main)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/bluenviron/gomavlib/v3)](https://pkg.go.dev/github.com/bluenviron/gomavlib/v3#pkg-index)
-->

gomavlib is a library that implements the Mavlink protocol (2.0 and 1.0) in the Go programming language. It can interact with Mavlink-capable devices through a serial port, UDP, TCP or a custom transport, and it can be used to power UGVs, UAVs, ground stations, monitoring systems or routers.

Mavlink is a lightweight and transport-independent protocol that is mostly used to communicate with unmanned ground vehicles (UGV) and unmanned aerial vehicles (UAV, drones, quadcopters, multirotors). It is supported by the most popular open-source flight controllers (Ardupilot and PX4).

First I want to give a big shoutout to the folks at bluenviron for doing an amazing job in writing this in the first 
place.  Thank you.

This library is a modified fork of the gomavlib at [https://github.com/bluenviron/gomavlib](https://github.com/bluenviron/gomavlib)

## Table of contents
* [Differences](#differences)
* [Features](#features)
* [Installation](#installation)
* [API Documentation](#api-documentation)
* [Dialect generation](#dialect-generation)
* [Testing](#testing)
* [Specifications](#specifications)
* [Links](#links)

## Differences
This implementation just focuses on the messages and frames of the mavlink protocol. All endpoints, and transport 
decisions have been dropped.  We felt it was best to decouple any networking and transport from the actual message.
Rather, these features should be implemented in their own library or application.  There is a high probability that
we will add another library with all the nodes and endpoints.  By decoupling these features, we now have the freedom 
to just create a message and let us worry about what we then want to do with it or where we want to send it.

Secondly, this library only supports Mavlink V2.

## Features:

* Decode and encode Mavlink v2.0 and v1.0. Supports checksums, empty-byte truncation (v2.0), signatures (v2.0), message extensions (v2.0).
* Dialects are optional, the library can work with standard dialects (ready-to-use standard dialects are provided in directory `dialects/`), custom dialects or no dialects at all. In case of custom dialects, a dialect generator is available in order to convert XML definitions into their Go representation.
* Examples provided for every feature, comprehensive test suite, continuous integration
* Easily create messages and frames


## Installation

By default only the minimal, standard and common dialects are generated and in the repo. After cloning the repo you
should generate the dialects you are interested in. See [Dialect Generation](#dialect-generation).

```go
import "github.com/merlindrones/gomavlib"
```

## API Documentation

[Click to open the API Documentation](https://pkg.go.dev/github.com/merlindrones/gomavlib/#pkg-index)

## Dialect generation

There are no dialects with the source code. These must be generated for us with the library.
There is a tool to generate these; just follow the commands below. the tool can take 1 flag:
`--dialects=[, sep list of dialect names]`

Without any flag it will generate:
* minimal.xml
* standard.xml
* common.xml

These will be generated in a dir called dialects in the current dir. So its probably a good idea to do this in the
pkg dir of the git repo.

To get a list of all available dialects
see [https://mavlink.io/en/messages/#dialects](https://mavlink.io/en/messages/#dialects)

```bash
go install github.com/merlindrones/gomavlib/cmd/gen-mavlink-dialects@latest
gen-mavlink-dialects
```

## Testing

If you want to hack the library and test the results, unit tests can be launched with:

```bash
make [test | testwithcoverage]
```

## Building Source

There is a make file with many targets available. To get a list of the targets:

```bash
make help
make debug
```

## Specifications

|name|area|
|----|----|
|[main website](https://mavlink.io/en/)|protocol|
|[packet format](https://mavlink.io/en/guide/serialization.html)|protocol|
|[common dialect](https://github.com/mavlink/mavlink/blob/master/message_definitions/v1.0/common.xml)|dialects|
|[Golang project layout](https://github.com/golang-standards/project-layout)|project layout|

## Links

Related projects

* [mavp2p](https://github.com/bluenviron/mavp2p)

Other Go libraries

* [gobot](https://github.com/hybridgroup/gobot/tree/master/platforms/mavlink)
* [liamstask/go-mavlink](https://github.com/liamstask/go-mavlink)
* [ungerik/go-mavlink](https://github.com/ungerik/go-mavlink)
* [SpaceLeap/go-mavlink](https://github.com/SpaceLeap/go-mavlink)
* [MAVSDK-Go](https://github.com/mavlink/MAVSDK-Go)

Other non-Go libraries

* [official library (C)](https://github.com/mavlink/c_library_v2)
* [pymavlink (Python)](https://github.com/ArduPilot/pymavlink)
* [mavlink.net (C#)](https://github.com/asvol/mavlink.net)
* [rust-mavlink (Rust)](https://github.com/3drobotics/rust-mavlink)
* [node-mavlink (JS)](https://github.com/omcaree/node-mavlink)
