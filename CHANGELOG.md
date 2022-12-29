# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased
### Added
- Add `akash_providers` datasource
### Changed
- `akash_deployment` output `services` now has a new structure and more information regarding replicas
### Development
- Integration with Praetor's API through custom caching service
- Start using Akash's type definitions
- Tracing enabled on deployment creation

## [0.0.6] - 2022-11-16
### Added
- Add support to mainnet-4
### Fixed
- Bug where deployments did not close when no bids were made

## [0.0.5] - 2022-10-27
### Added
- Transaction memo stating the Terraform provider was used and which version
### Fixed
- Bug where destroying deployments did not really close them
- Owner not being set in Go client for bid listing
- Issue on temporary folder on some OS
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
