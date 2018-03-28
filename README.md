# envctl

[![Build Status](https://travis-ci.org/UltimateSoftware/envctl.svg?branch=master)](https://travis-ci.org/UltimateSoftware/envctl)

-----

Every codebase has a set of tools developers need to work with it. Managing that
set of tools is usually tricky for a bunch of reasons.

But it really doesn't have to be complicated.

Envctl manages a codebase's tools by allowing its authors to maintain a sandbox
environment with all the stuff they need to work on it.

## Quick Start

**NOTE:** This assumes you have [Docker](https://www.docker.com/) installed.

```bash
$ go install github.com/UltimateSoftware/envctl
$ cd $HOME/src/my-repo
$ envctl init
$ envctl create
$ $EDITOR envctl.yaml
$ envctl login # do stuff, then exit
$ envctl destroy
```

## Configuration Guide

The configuration takes the following format:
```yaml
---
# Required - the base container image for the environment
image: ubuntu:latest

# Required - the shell to use when logged in
shell: /bin/bash

# The mount directory inside the container for the repo
mount: /mnt/repo

# An array of commands to run in the specified shell when creating the
# environment.
bootstrap:
- ./bootstrap.sh
- ./extra-config.sh

# An array of environment variables. Anything with a $ will be evaluated against
# the current set of exported variables being used by the current session.
variables:
- FOO=bar
- SECRET=$SECRET
```
