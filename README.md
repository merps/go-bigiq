
[//]: # (Original work Copyright © 2015 Scott Ware)
[//]: # (Modifications Copyright 2019 F5 Networks Inc)
[//]: # (Modifications Copyright © 2022 m.kennedy@f5.com - forked from go-bigip)
[//]: # (Licensed under the Apache License, Version 2.0 [the "License"];)
[//]: # (You may not use this file except in compliance with the License.)
[//]: # (You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0)
[//]: # (Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,)
[//]: # (WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.)
[//]: # (See the License for the specific language governing permissions and limitations under the License.)

## go-bigiq
[![GoDoc](https://godoc.org/github.com/merps/go-bigiq?status.svg)](https://godoc.org/github.com/merps/go-bigiq)
[![Build Status](https://app.travis-ci.com/merps/go-bigiq.svg?branch=main)](https://app.travis-ci.com/merps/go-bigiq)
[![Go Report Card](https://goreportcard.com/badge/github.com/merps/go-bigiq)](https://goreportcard.com/report/github.com/merps/go-bigiq)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/merps/go-bigiq/master/LICENSE)

A Go package that interacts with F5 BIG-IQ systems using the REST API.

Some of the tasks you can do are as follows:

* connect and authorise (basic-auth or token)
* license provisioning as per [**License Management**](https://clouddocs.f5.com/products/big-iq/mgmt-api/v0.0/HowToSamples/bigiq_public_api_wf/t_bigiq_public_api_workflows.html)
* basic BIG-IQ device system queries [**Device Management**](https://clouddocs.f5.com/products/big-iq/mgmt-api/v0.0/HowToSamples/bigiq_public_api_wf/t_bigiq_public_api_workflows.html)

> **Note**: You must be on version 8.0+! 

### Examples & Documentation
Initial examples are located within `examples/` path

### TODO
- [ ] Upload of License file based on manual/ccn activation.
- [ ] Additional inline TODO's as per code.
- [ ] Validate AS3 Upload/Download
- [ ] Correct GoPkg endpoint
- [ ] TaskId or Softlink output (Poll, Dossier, etc)

### Contributors
A very special thanks to the following who have helped contribute to earlier fork of this codebase.

* [Adam Burnett](https://github.com/aburnett)
* [Michael D. Ivey](https://github.com/ivey)
* [Scott Ware](https://github.com/scottdware/go-bigip)

And F5 A&O Ecosystems

[godoc-go-bigiq]: http://godoc.org/github.com/merps/go-bigiq
[license]: https://github.com/merps/go-bigiq/blob/master/LICENSE
