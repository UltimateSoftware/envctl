# envctl

[![Build Status](https://travis-ci.org/UltimateSoftware/envctl.svg?branch=master)](https://travis-ci.org/UltimateSoftware/envctl)

-----

Every codebase has a set of tools developers need to work with it. Managing that
set of tools is usually tricky for a bunch of reasons. Getting new devs set up with a
codebase can often get complicated.

But it really doesn't have to be complicated.

Envctl manages a codebase's tools by allowing its authors to maintain a sandbox
environment with all the stuff they need to work on it.

You can find a sample repo with an environment already set up [here](https://github.com/juicemia/envctl-sample).

## Installation Guide

To install `envctl`, just download the [current release](https://github.com/UltimateSoftware/envctl/releases/tag/2.0.0) and extract the binary to somewhere in your `$PATH`.

Alternatively, if you have `go` installed, you can compile from source.

## Quick Start

```bash
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

# Specifies whether the base image should be cached. Defaults to true.
cache_image: false

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
# the current set of exported variables being used by the current session. If
# any of them evaluate to nothing, envctl will fail to create the environment.
variables:
  FOO: bar
  SECRET: $SECRET

# A map of layer 3 protocols to ports that can be exposed by Docker.
ports:
  tcp:
  - 4567
```

## Contributing Guide

- If you're new to Go, or don't know quite where to start, feel free to ask for
help. Check the issues for things labeled "good first issue".
- Pull requests are always welcome, no matter how crazy they are.
- Write tests, and make sure `go test ./...` passes.
