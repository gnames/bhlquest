# Changelog

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

- Add parameters for maximum number of returns and score threshold (#13)
- Change limiting factor from main_taxa to main_class of an item (#12)
- Add better hanling of initiation process (#11)
- Add text of matched chunk to the output (if possible) (#10)
- Add an option to limit BHL items intake by a list of taxa (#9)

### Added v0.0.1

- Create API endpoint (#8).
- Ask a question and get an answer (#7).
- Embed chunks and save to the database (#6).
- Break items into chunks (#5).
- Connect to llmutil RESTful service (#4).
- Create database reset with pgvector for the data (#3).
- Return list of item ids (#2).


## [v0.0.0] - 2023-11-11

### Added v0.0.0

- General "plumbing": filestructure makefile etc (#1)

## Footnotes

This document follows [changelog guidelines]

[v0.0.2]: https://github.com/gnames/bhlquest/compare/v0.0.1...v0.0.2
[v0.0.1]: https://github.com/gnames/bhlquest/compare/v0.0.0...v0.0.1
[v0.0.0]: https://github.com/gnames/bhlquest/tree/v0.0.0

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
