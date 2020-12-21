# Changelog

## [Unreleased]

### Changed

- Change type of service keys from `string` to `interface{}`. [#4](https://github.com/go-nacelle/service/pull/4)

### Removed

- Removed MustGet and MustSet methods from the Container interface. [#3](https://github.com/go-nacelle/service/pull/3)

## [v1.0.2] - 2020-09-30

### Removed

- Removed dependency on [aphistic/sweet](https://github.com/aphistic/sweet) by rewriting tests to use [testify](https://github.com/stretchr/testify). [#2](https://github.com/go-nacelle/service/pull/2)

## [v1.0.1] - 2020-09-07

### Added

- Added overlay container. [#1](https://github.com/go-nacelle/service/pull/1)

## [v1.0.0] - 2019-06-17

### Changed

- Migrated from [efritz/bussard](https://github.com/efritz/bussard).

[Unreleased]: https://github.com/go-nacelle/service/compare/v1.0.2...HEAD
[v1.0.0]: https://github.com/go-nacelle/service/releases/tag/v1.0.0
[v1.0.1]: https://github.com/go-nacelle/service/compare/v1.0.0...v1.0.1
[v1.0.2]: https://github.com/go-nacelle/service/compare/v1.0.1...v1.0.2
