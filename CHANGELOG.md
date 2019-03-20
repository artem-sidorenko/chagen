Changelog
=========

## [v0.0.3](https://github.com/artem-sidorenko/chagen/releases/v0.0.3) (20.03.2019)

Closed issues
-------------
- Add own tap for homebrew on MacOS [\#41](https://github.com/artem-sidorenko/chagen/issues/41)
- GitLab support [\#22](https://github.com/artem-sidorenko/chagen/issues/22)
- Introduce go routines for parallel fetching of information [\#39](https://github.com/artem-sidorenko/chagen/issues/39)
- Tag filter should be implemented [\#16](https://github.com/artem-sidorenko/chagen/issues/16)
- Integration tests [\#17](https://github.com/artem-sidorenko/chagen/issues/17)

Merged pull requests
--------------------
- Docs: installation instructions [\#65](https://github.com/artem-sidorenko/chagen/pull/65) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Use stderr for progress information [\#64](https://github.com/artem-sidorenko/chagen/pull/64) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Use the same testing data everywhere [\#63](https://github.com/artem-sidorenko/chagen/pull/63) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Refactoring: datasource is a better name for data sources [\#62](https://github.com/artem-sidorenko/chagen/pull/62) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Refactoring: lets have projectID generation once [\#61](https://github.com/artem-sidorenko/chagen/pull/61) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Renaming: better names following the conventions [\#60](https://github.com/artem-sidorenko/chagen/pull/60) ([artem-sidorenko](https://github.com/artem-sidorenko))
- GitLab connector [\#59](https://github.com/artem-sidorenko/chagen/pull/59) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Preparations for GitLab connector [\#58](https://github.com/artem-sidorenko/chagen/pull/58) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Improvement: lets have the max counters first [\#56](https://github.com/artem-sidorenko/chagen/pull/56) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Display fetching progress of PRs and issues [\#57](https://github.com/artem-sidorenko/chagen/pull/57) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Testing fix: use the same setup for coverage calc [\#55](https://github.com/artem-sidorenko/chagen/pull/55) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Fix: race condition with wait group [\#54](https://github.com/artem-sidorenko/chagen/pull/54) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Test settings adjustment [\#52](https://github.com/artem-sidorenko/chagen/pull/52) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Testing: more issues and PRs in the simulated GH client [\#53](https://github.com/artem-sidorenko/chagen/pull/53) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Feature: show the progress of tags [\#51](https://github.com/artem-sidorenko/chagen/pull/51) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Test improvements: much more tags for simulated GH client [\#50](https://github.com/artem-sidorenko/chagen/pull/50) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Fix: follow the good pipeline patterns [\#49](https://github.com/artem-sidorenko/chagen/pull/49) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Minor refactoring and fixes [\#48](https://github.com/artem-sidorenko/chagen/pull/48) ([artem-sidorenko](https://github.com/artem-sidorenko))
- First drop of goroutines for fetching of information [\#47](https://github.com/artem-sidorenko/chagen/pull/47) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Fix: display custom err message by missing repo [\#46](https://github.com/artem-sidorenko/chagen/pull/46) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Refactor setup cli context [\#44](https://github.com/artem-sidorenko/chagen/pull/44) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Tests: compare errros via reflect.DeepEqual [\#43](https://github.com/artem-sidorenko/chagen/pull/43) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Linting: making some new linter happy [\#45](https://github.com/artem-sidorenko/chagen/pull/45) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Feature: custom err message for missing repos [\#42](https://github.com/artem-sidorenko/chagen/pull/42) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Feature: filtering of labels [\#38](https://github.com/artem-sidorenko/chagen/pull/38) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Feature: enable tag filtering per default [\#36](https://github.com/artem-sidorenko/chagen/pull/36) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Refactoring: testing of github connector [\#35](https://github.com/artem-sidorenko/chagen/pull/35) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Generator: footer with chagen information [\#34](https://github.com/artem-sidorenko/chagen/pull/34) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Big refactoring [\#33](https://github.com/artem-sidorenko/chagen/pull/33) ([artem-sidorenko](https://github.com/artem-sidorenko))
- CI: small improvements [\#32](https://github.com/artem-sidorenko/chagen/pull/32) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Bugfix: don't print error if there is no text [\#31](https://github.com/artem-sidorenko/chagen/pull/31) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Bugfix: include issues/MRs in the same second like release [\#30](https://github.com/artem-sidorenko/chagen/pull/30) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Better error handling and tiny style improvements [\#29](https://github.com/artem-sidorenko/chagen/pull/29) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Update of dependencies to the latest state [\#28](https://github.com/artem-sidorenko/chagen/pull/28) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Get it back to live [\#27](https://github.com/artem-sidorenko/chagen/pull/27) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Update of go version within travis [\#26](https://github.com/artem-sidorenko/chagen/pull/26) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Refactoring: tests and common commands code [\#25](https://github.com/artem-sidorenko/chagen/pull/25) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Adding some badges [\#20](https://github.com/artem-sidorenko/chagen/pull/20) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Fix: type in the release file names [\#19](https://github.com/artem-sidorenko/chagen/pull/19) ([artem-sidorenko](https://github.com/artem-sidorenko))

## [v0.0.2](https://github.com/artem-sidorenko/chagen/releases/tag/v0.0.2) (04.08.2017)

Closed issues
-------------
- Upload of released builds and packages to GitHub and OBS from Travis [\#15](https://github.com/artem-sidorenko/chagen/issues/15)
- Signing of git release tags with gpg [\#14](https://github.com/artem-sidorenko/chagen/issues/14)

Merged pull requests
--------------------
- Upload and signing of releases [\#18](https://github.com/artem-sidorenko/chagen/pull/18) ([artem-sidorenko](https://github.com/artem-sidorenko))

## [v0.0.1](https://github.com/artem-sidorenko/chagen/releases/tag/v0.0.1) (04.08.2017)

Merged pull requests
--------------------
- Release prepararation logic in the Makefile [\#13](https://github.com/artem-sidorenko/chagen/pull/13) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Archiving of release builds [\#12](https://github.com/artem-sidorenko/chagen/pull/12) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Enable preparation of new release, which is not tagged yet [\#11](https://github.com/artem-sidorenko/chagen/pull/11) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Lets have reverse sorting everywhere [\#10](https://github.com/artem-sidorenko/chagen/pull/10) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Enabling CLI flags for specifying a github repo [\#9](https://github.com/artem-sidorenko/chagen/pull/9) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Refactoring of data types [\#8](https://github.com/artem-sidorenko/chagen/pull/8) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Sorting and filtering of release data [\#7](https://github.com/artem-sidorenko/chagen/pull/7) ([artem-sidorenko](https://github.com/artem-sidorenko))
- First integration between connector and changelog generator [\#6](https://github.com/artem-sidorenko/chagen/pull/6) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Support of GitHub API authentication via environment variable [\#5](https://github.com/artem-sidorenko/chagen/pull/5) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Basic GitHub connector [\#4](https://github.com/artem-sidorenko/chagen/pull/4) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Use gometalinter for linting [\#3](https://github.com/artem-sidorenko/chagen/pull/3) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Refactoring: lets register our commands via init() [\#2](https://github.com/artem-sidorenko/chagen/pull/2) ([artem-sidorenko](https://github.com/artem-sidorenko))
- Enable Travis CI [\#1](https://github.com/artem-sidorenko/chagen/pull/1) ([artem-sidorenko](https://github.com/artem-sidorenko))
