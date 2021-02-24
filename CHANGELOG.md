# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic
Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased

### Added
- Added --use-ec2-region for getting region when ran on an EC2 instance

## [0.3.0] - 2021-01-20

### Changed
- Moved to sensu org
- README updates regarding annotations, remove empty sections, general cleanup
- Ran 'go get -u' and 'go mod tidy' to update all modules
- GitHub Actions: add lint, add pull_request for test
- Capture and output published Message ID, note for future use

## [0.2.0] - 2020-12-04

### Changed
- Added assume role arn support

### Changed
- Updated to use SDK templating
- Updated SDK version to v0.11.0
- Updated dependent modules
- Changed types import to corev2

## [0.1.0] - 2020-04-16

### Added
- Initial release
