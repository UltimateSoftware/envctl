# envctl

Every codebase has a set of tools developers need to work with it. This is usually language dependent.

For example, a Ruby codebase will need a particular version of Ruby, probably with Bundler installed.
There might be some additional required tools for deployment of the code in question, or other things
developers need to do with the code.

Usually what ends up happening is one (or a couple) of engineers will get together and start a project.
Each will have their own workstation set up in a particular way. The project gets going, progress is
being made quickly, and then finally, someone sits down to document how to get started with the codebase.

Inevitably, some steps will be missed. Most of the time, developers aren't starting really starting from
scratch.

Even then, documentation isn't fun to write.

So now, three months later, another team member is getting started with the project. Inevitably, they
have some trouble following along with the docs. Things don't quite work the same way because, just
like all other developers, their workstation is set up in a particular way. An original developer has
to stop what they are doing to help the new contributor get set up, and it's always painful because
many times the original developer doesn't even remember how his own workstation is set up.

Step back and examine how much time has been wasted. There was the original time spent writing the docs,
then the time spent by the new contributor reading them and struggling to get set up, followed by the
combined time of two or more people trying to get it done together.

Envctl is a time machine. It's purpose is to give you back that time.

## Quick Start

1. Install `envctl`.
```
$ go get github.com/UltimateSoftware/envctl
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
