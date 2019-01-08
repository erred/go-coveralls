# go-coveralls

[![Godoc](http://img.shields.io/badge/godoc-reference-5272b4.svg?maxAge=31536000&style=for-the-badge)](http://godoc.org/github.com/seankhliao/go-coveralls)
[![License](https://img.shields.io/github/license/seankhliao/go-coveralls.svg?style=for-the-badge)](githib.com/seankhliao/go-coveralls)
[![Build Status](https://img.shields.io/travis-ci/seankhliao/go-coveralls.svg?style=for-the-badge)](https://travis-ci.org/seankhliao/go-coveralls)
[![Codecov](https://img.shields.io/codecov/c/github/seankhliao/go-coveralls.svg?style=for-the-badge)](https://codecov.io/gh/seankhliao/go-coveralls)
![](https://img.shields.io/github/tag/seankhliao/go-coveralls.svg?style=for-the-badge)

[![Go Report Card](https://goreportcard.com/badge/github.com/seankhliao/go-coveralls?style=flat-square)](https://goreportcard.com/report/github.com/seankhliao/go-coveralls)

## API Structure

```
Logical Structure:
Provider -> Owner -> Repo -> Build -> Job

Endpoints

Add Repo

Submit Job


Latest build
Page of 10 builds

Single Build
Single Build Subset
Single Build Single File Coverage

Single Job
Single Job Single File Coverage

```
