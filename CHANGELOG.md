# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Fixed
- Bug where destroying deployments did not really close them
### Development
- Replaced HTTP API call for deployments with CLI

## [0.0.4] - 2022-08-17
### Added
- Introduce provider filters with `enforce` and `providers` filters
### Changed
- Make temporary deployment file location cross-platform
### Fixed
- `net` field in provider configuration had wrong default value
- Bug where cheapest bids were not being selected
- Issue where gas adjustment was not enough on deployment update
### Development
- More unit tests
- Implemented several string utilities including `CointainsAny` and `FindAll` functions
