# Changelog

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [ v0.0.19] - 2024-12-02 Mon

- Add citation.cff.

## [v0.0.15] - 2024-01.18 Thu

- Fix typo

## [v0.0.14] - 2024-01-18 Thu

- Clean up API

## [v0.0.13] - 2024-01-18 Thu

- Fix naming, config

## [v0.0.12] - 2024-01-18 Thu

- Fix: gnbhl reference

## [v0.0.11] - 2024-01-18 Thu

- Improve chunks in result (#26)
- Add references to results (#25)

## [v0.0.10] - 2024-01-04 Thu

- Add bigger texts for JSON
- Prepare configuration of new features for Kubernetes.

## [v0.0.9] - 2023-12-24 Sat

### Added v0.0.9

- Add Cohere to rerank resuts (#24)

## [v0.0.8] - 2023-12-20 Wed

### Added v0.0.8

- Add Qdrant for data initiation (#23)
- Add Qdrant for querying (#21)
- Add Cross-Embedding (#22)

## [v0.0.7] - 2023-12-13 Wed

### Added v0.0.7

- Summary using OpenAi chatGPT 3.5 Turbo API (#20).

## [v0.0.6] - 2023-11-30

### Added v0.0.6

- Use model with bigger vectors

## [v0.0.5] - 2023-11-29

### Added v0.0.5

- Example questions for the web GUI (#19)

## [v0.0.4] - 2023-11-27

### Added v0.0.4

- Flag to show or hide matched texts (#18)
- Update web GUI form with GET data (#17)

## [v0.0.3] - 2023-11-23

### Added v0.0.3

- Add steps needed for Kubernetes deployment (#16)
- Add web GUI (#15)

## [v0.0.2] - 2023-11-17

### Added v0.0.2

- Add max-results and score-threshold parameters to API (#14)
- Add parameters for maximum number of returns and score threshold (#13)
- Change limiting factor from main_taxa to main_class of an item (#12)
- Add better handling of initiation process (#11)
- Add text of matched chunk to the output (if possible) (#10)
- Add an option to limit BHL items intake by a list of taxa (#9)

## [v0.0.1] - 2023-11-15

### Added v0.0.1

- Create API endpoint (#8).
- Ask a question and get an answer (#7).
- Embed chunks and save to the database (#6).
- Break items into chunks (#5).
- Connect to `llmutil` RESTful service (#4).
- Create database reset with `pgvector` for the data (#3).
- Return list of item IDs (#2).

## [v0.0.0] - 2023-11-11

### Added v0.0.0

- General "plumbing": files structure, makefile etc. (#1)

## Footnotes

This document follows [changelog guidelines]

[v0.0.9]: https://github.com/gnames/bhlquest/compare/v0.0.8...v0.0.9
[v0.0.8]: https://github.com/gnames/bhlquest/compare/v0.0.7...v0.0.8
[v0.0.7]: https://github.com/gnames/bhlquest/compare/v0.0.6...v0.0.7
[v0.0.6]: https://github.com/gnames/bhlquest/compare/v0.0.5...v0.0.6
[v0.0.5]: https://github.com/gnames/bhlquest/compare/v0.0.4...v0.0.5
[v0.0.4]: https://github.com/gnames/bhlquest/compare/v0.0.3...v0.0.4
[v0.0.3]: https://github.com/gnames/bhlquest/compare/v0.0.2...v0.0.3
[v0.0.2]: https://github.com/gnames/bhlquest/compare/v0.0.1...v0.0.2
[v0.0.1]: https://github.com/gnames/bhlquest/compare/v0.0.0...v0.0.1
[v0.0.0]: https://github.com/gnames/bhlquest/tree/v0.0.0
[#30]: https://codeberg.org/dimus/madcow/issues/30
[#29]: https://codeberg.org/dimus/madcow/issues/29
[#28]: https://codeberg.org/dimus/madcow/issues/28
[#27]: https://codeberg.org/dimus/madcow/issues/27
[#26]: https://codeberg.org/dimus/madcow/issues/26
[#25]: https://codeberg.org/dimus/madcow/issues/25
[#24]: https://codeberg.org/dimus/madcow/issues/24
[#23]: https://codeberg.org/dimus/madcow/issues/23
[#22]: https://codeberg.org/dimus/madcow/issues/22
[#21]: https://codeberg.org/dimus/madcow/issues/21
[#20]: https://codeberg.org/dimus/madcow/issues/20
[#19]: https://codeberg.org/dimus/madcow/issues/19
[#18]: https://codeberg.org/dimus/madcow/issues/18
[#17]: https://codeberg.org/dimus/madcow/issues/17
[#16]: https://codeberg.org/dimus/madcow/issues/16
[#15]: https://codeberg.org/dimus/madcow/issues/15
[#14]: https://codeberg.org/dimus/madcow/issues/14
[#13]: https://codeberg.org/dimus/madcow/issues/13
[#12]: https://codeberg.org/dimus/madcow/issues/12
[#11]: https://codeberg.org/dimus/madcow/issues/11
[#10]: https://codeberg.org/dimus/madcow/issues/10
[#9]: https://codeberg.org/dimus/madcow/issues/9
[#8]: https://codeberg.org/dimus/madcow/issues/8
[#7]: https://codeberg.org/dimus/madcow/issues/7
[#6]: https://codeberg.org/dimus/madcow/issues/6
[#5]: https://codeberg.org/dimus/madcow/issues/5
[#4]: https://codeberg.org/dimus/madcow/issues/4
[#3]: https://codeberg.org/dimus/madcow/issues/3
[#2]: https://codeberg.org/dimus/madcow/issues/2
[#1]: https://codeberg.org/dimus/madcow/issues/1
[changelog guidelines]: https://keepachangelog.com/en/1.0.0/
