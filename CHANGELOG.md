# Changelog


## [Unreleased]

### Removed

### Added

### Improved

### Deprecated


## [0.4.0](https://github.com/mozillazg/kube-audit-mcp/compare/v0.3.1...v0.4.0) (2025-09-06)

### Added

- Add new subcommand `test` ([#15](https://github.com/mozillazg/kube-audit-mcp/pull/15))
- Add gcp-cloud-logging provider ([#10](https://github.com/mozillazg/kube-audit-mcp/pull/10))
- Add Dockerfile and release workflow for docker image ([#11](https://github.com/mozillazg/kube-audit-mcp/pull/11))

### Improved

- Improve query tool and update documentation ([#17](https://github.com/mozillazg/kube-audit-mcp/pull/17))
- Update go dependencies ([#8](https://github.com/mozillazg/kube-audit-mcp/pull/8))


## [0.3.1](https://github.com/mozillazg/kube-audit-mcp/compare/v0.3.0...v0.3.1) (2025-08-11)

### Added

- Add README.zh-CN.md

### Fixed

- Ensure default config file path is `~/.config/kube-audit-mcp/config.yaml`

## [0.3.0](https://github.com/mozillazg/kube-audit-mcp/compare/v0.2.0...v0.3.0) (2025-08-10)

### Added

- Add new subcommand `version`
- Add new tool `list_clusters`
- Add new param `cluster_name` for `query_audit_log` tool

### Improved

- Update provider documents
- Skip saving existing file in `sample-config`
- Update dependencies

### Fixed

- Remove linkmode in CI

## [0.2.0](https://github.com/mozillazg/kube-audit-mcp/compare/v0.1.0...v0.2.0) (2025-08-10)

### Added

- Add Installation and Configurations to docs
- Add `make build` and `make test` to CI
- Support for multiple clusters
- Load configuration from `~/.config/kube-audit-mcp/config.yaml` by default
- Add new provider `aws-cloudwatch-logs`

### Improved

- Fix documentation details
- Support setting endpoint via region for Alibaba Cloud provider
- Rename `alibaba_sls` to `alibaba-sls`

### Fixed

- Fix tests



## 0.1.0 (2025-08-03)

This is the first release of the project.
