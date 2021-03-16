# Contributing to the Meroxa CLI

Thank you so much for contributing to our CLI. We appreciate your time and help.
Here are some guidelines to help you get started.

## Filing a bug or feature

1. Before filing an issue, please check the existing issues to see if a
   similar one was already opened. If there is one already opened, feel free
   to comment on it.
1. If you believe you've found a bug, please provide detailed steps of
   reproduction, the version of Meroxa and anything else you believe will be
   useful to help troubleshoot it (e.g. OS environment, environment variables,
   etc...). Also state the current behavior vs. the expected behavior. Open an issue directly using [our template to file a bug](https://github.com/meroxa/cli/issues/new?assignees=&labels=&template=bug_report.md&title=).
1. If you'd like to see a feature or an enhancement please open an issue using [our template to file a feature request](https://github.com/meroxa/cli/issues/new?assignees=&labels=&template=feature_request.md&title=).

## Submitting changes

1. Tests: If you are submitting code, please ensure you have adequate tests
   for the feature. Tests can be run via `go test ./...` or `make test`.
1. Since this is golang project, ensure the new code is properly formatted to
   ensure code consistency. Run `go fmt`.

### Quick steps to contribute

1. Fork the project.
1. Download your fork to your machine (`git clone https://github.com/your_username/cli && cd cli`)
1. Create your feature branch (`git checkout -b my-new-feature`)
1. Make changes and run tests (`make test`)
1. Add them to staging (`git add .`)
1. Commit your changes (`git commit -m 'Add some feature'`)
1. Push to the branch (`git push origin my-new-feature`)
1. Create new pull request

## License

Apache 2.0, see [LICENSE](LICENSE).

## Code of Conduct

### Our Pledge

In the interest of fostering an open and welcoming environment, we as
contributors and maintainers pledge to making participation in our project and
our community a harassment-free experience for everyone, regardless of age, body
size, disability, ethnicity, gender identity and expression, level of experience,
nationality, personal appearance, race, religion, or sexual identity and
orientation.

### Our Standards

Examples of behavior that contributes to creating a positive environment
include:

* Using welcoming and inclusive language
* Being respectful of differing viewpoints and experiences
* Gracefully accepting constructive criticism
* Focusing on what is best for the community
* Showing empathy towards other community members

Examples of unacceptable behavior by participants include:

* The use of sexualized language or imagery and unwelcome sexual attention or
  advances
* Trolling, insulting/derogatory comments, and personal or political attacks
* Public or private harassment
* Publishing others' private information, such as a physical or electronic
  address, without explicit permission
* Other conduct which could reasonably be considered inappropriate in a
  professional setting

### Our Responsibilities

Project maintainers are responsible for clarifying the standards of acceptable
behavior and are expected to take appropriate and fair corrective action in
response to any instances of unacceptable behavior.

Project maintainers have the right and responsibility to remove, edit, or
reject comments, commits, code, wiki edits, issues, and other contributions
that are not aligned to this Code of Conduct, or to ban temporarily or
permanently any contributor for other behaviors that they deem inappropriate,
threatening, offensive, or harmful.

### Scope

This Code of Conduct applies both within project spaces and in public spaces
when an individual is representing the project or its community. Examples of
representing a project or community include using an official project e-mail
address, posting via an official social media account, or acting as an appointed
representative at an online or offline event. Representation of a project may be
further defined and clarified by project maintainers.

### Enforcement

Instances of abusive, harassing, or otherwise unacceptable behavior may be
reported by contacting our team at [support@meroxa.io](mailto:support@meroxa.io). All
complaints will be reviewed and investigated and will result in a response that
is deemed necessary and appropriate to the circumstances. Meroxa Inc is
obligated to maintain confidentiality with regard to the reporter of an incident.
Further details of specific enforcement policies may be posted separately.

Project maintainers who do not follow or enforce the Code of Conduct in good
faith may face temporary or permanent repercussions as determined by other
members of the project's leadership.

### Attribution

This Code of Conduct is adapted from the Contributor Covenant, version 1.4, available at http://contributor-covenant.org/version/1/4.

