---
name: Bug report
description: File a bug report
about: Create a report to help us improve
title: '[Bug]: '
labels: ['bug', 'CLI']
assignees: ''

---

body:
  - type: markdown
    id: current
    attributes:
      label: Current behaviour
      placeholder: From a high level, describe the problem you ran into.
      value: |
  - type: markdown
    id: expected
    attributes:
      label: Expected behaviour
      placeholder: What did you expect to happen?
      value: |
  - type: markdown
    id: reproduce
    attributes:
      label: What can we do to reproduce your bug?
      placeholder: Please provide: Relevant CLI commands, relevant environment details
      value: |
  - type: textarea
    id: error
    attributes:
      label: Error messages
      description: Please provide any error messages that you encountered
      render: 'shell'
    validations:
      required: true
  - type: input
    id: machine-version
    attributes:
      label: Operating System
      description: Which machine did you use when testing this?
      placeholder: e.g. MacOS
    validations:
      required: true
  - type: input
    id: tool-version
    attributes:
      label: CLI version
      description: Which version of the CLI did you use?
      placeholder: e.g. 1.22.3
  - type: textarea
    id: screens
    attributes:
      label: Screenshots
      description: If applicable, add screenshots to help explain your problem. (Be sure to scrub any sensitive information.)
