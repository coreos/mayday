# Change Log

All notable changes to the project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/) 
and this project adheres to [Semantic Versioning](http://semver.org/).

## [1.0.0]
### Added
- Data collection can now be specified through a config file
  - Config file can be specified via the arguments `-c` & `--config`
  - Config files specify a list of files and commands to be collected
  - A config file can be referenced from it's default locations as a "profile"
    (`-p`, `--profile`) allowing for users to keep custom sets of telemetry for
    different products
- systemd style directory preference support for config files/profiles
  - `/etc/mayday`
  - `/usr/share/mayday`
- "danger" flag (`-d`, `--danger`) which will gather more logfiles and 
  _potentially_ sensitive information
- Output file selection (`-o`, `--output`)
- Collection / Commands
  - slabtop
  - iptables / ip6tables
  - rkt support
  - docker support
- Collection of systemd info
- Journal collection for system supplied units

### Changed
- Improved ASCII artwork
- Collection specification moved into configuration file
- Removed godeps
- Dependency vendoring is now done through [glide](https://github.com/Masterminds/glide)
- Improved error messaging
- Fille collection and tar creation now executes in memory
- Updated yaml v2 vendored package
- Flag mechanism now uses [viper](https://github.com/spf13/viper)

## [0.1.0] 2014-10-20
### Added
- Collection of disk statistics from host
- Collection of memory statistics from host
- Added network related commands
- Adding general OS collection commands

[Unreleased]: https://github.com/coreos/mayday/compare/v1.0.0...HEAD
[1.0.0]:      https://github.com/coreos/mayday/compare/v0.1.0...v1.0.0
[0.1.0]:      https://github.com/coreos/mayday/compare/455bde42...v0.1.0

<!--
> vim:set ts=2 sw=2 expandtab:
-->
