# envctl

Every codebase has a set of tools developers need to work with it. Managing that
set of tools is usually tricky for a bunch of reasons.

But it really doesn't have to be complicated.

Envctl manages a codebase's tools by allowing its authors to maintain a sandbox
environment with all the stuff they need to work on it.

## Quick Start

**NOTE:** This assumes you have [Docker](https://www.docker.com/) installed.

```
$ go install github.com/UltimateSoftware/envctl
$ cd $HOME/src/my-repo
$ envctl init
$ envctl create
$ $EDITOR envctl.yaml
$ envctl login # do stuff, then exit
$ envctl destroy
```
