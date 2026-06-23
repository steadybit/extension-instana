# Changelog

## v1.1.18

- chore(deps): bump github.com/steadybit/extension-kit
- chore(deps): bump golang.org/x/net to v0.55.0 (CVE-2026-39821) (#80)

## v1.1.17

- chore(deps): bump alpine from 3.23 to 3.24

## v1.1.16

- chore: update to go 1.26.4
- feat: add weekly auto patch-release workflow

## v1.1.15

- Support discovery group attribute via `STEADYBIT_EXTENSION_DISCOVERY_GROUP` env var (or `discovery.group` Helm value) — when set, the extension adds `steadybit.group=<value>` to every discovered target
- Update dependencies

## v1.1.14

- Bump Go to 1.26.3
- Update dependencies

## v1.1.13

- Bump Go to 1.25.9
- Support if-none-match for the extension list endpoint
- Update dependencies

## v1.1.12

- feat(chart): split image.name into image.registry + image.name
- Support global.priorityClassName
- Update alpine packages in Docker image to address CVEs
- Update dependencies

## v1.1.11

- Update dependencies

## v1.1.10

- Update dependencies

## v1.1.9

- Update dependencies

## v1.1.8

- Updated dependencies

## v1.1.7

- Updated dependencies
- add insecureSkipVerify option

##  v1.1.6

- update dependencies
- Use uid instead of name for user statement in Dockerfile

##  v1.1.5

- Set new `Technology` property in extension description
- Update dependencies (go 1.23)

## v1.1.4

- Update dependencies (go 1.22)

## v1.1.3

 - Update dependencies

## v1.1.2

 - Update dependencies

## v1.1.1

 - Update dependencies

## v1.1.0

 - Added a discovery for application perspectives
 - Added an action to create a maintenance window
 - Filter event check based on application perspective(s)
 - Events are shown in a timeline with clickable links to the event details

## v1.0.0

 - Initial release
