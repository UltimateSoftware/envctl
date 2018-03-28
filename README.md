# envctl

Every codebase has a set of tools developers need to work with it. Managing that
set of tools is usually tricky for a bunch of reasons.

But it really doesn't have to be complicated.

Envctl manages a codebase's tools by allowing its authors to maintain a sandbox
environment with all the stuff they need to work on it.

## Quick Start

**NOTE:** This assumes you have [Docker](https://www.docker.com/) installed.

1. Install `envctl`.
```
$ go install github.com/UltimateSoftware/envctl
```
2. Change directory to your repo.
```
$ cd $HOME/src/my-repo
```
3. Generate a configuration file
```
$ envctl init
```
4. After editing the configuration file, create your environment.
```
$ envctl create
```
5. Log in and poke around
```
$ envctl login
```
6. Destroy it when you don't need it anymore.
```
$ envctl destroy
```
