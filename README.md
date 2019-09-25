PicFight coin regression testing
=======
[![Build Status](http://img.shields.io/travis/picfight/pfcregtest.svg)](https://travis-ci.org/picfight/pfcregtest)
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)

Harbours a pre-configured test setup and unit-tests to run RPC-driven node tests.

Builds a testing harness crafting and executing integration tests by driving a `pfcd` and `pfcwallet` instances via the `RPC` interface.

## Build 

```
set GO111MODULE=on
go build ./...
go clean -testcache
go test ./...
 ```
 
 ## License
 This code is licensed under the [copyfree](http://copyfree.org) ISC License.
