# rig

`rig` is a kubernetes manifest pre-processor and templating tool.

- Templates can be stored in remote git repositories, enabling the ability
  to share templates across apps.
- Built on go templates with support for all the templating features of Helm.

## Usage

```
rig help
```

## Installing

Install using the install script:

```shell
curl https://raw.githubusercontent.com/gonstr/rig/master/install.sh | sh
```

Install using `go get`:

```shell
go get github.com/gonstr/rig
```

Or just download binaries from the latest release.

## Docker image

[gonstr/rig](https://cloud.docker.com/u/gonstr/repository/docker/gonstr/rig) has `rig` and `git` installed.
