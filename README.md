# rig

`rig` is a kubernetes manifest pre-processor and templating tool.

- Templates are stored in remote git repositories, enabling the ability
  to share templates across apps.
- Built on go templates with support for all the templating features of Helm.

## Example usage

Install a rig template:

```shell
rig install https://github.com/foo/bar/simple-app#simple-app/v1.0.0
```

A `rig.yaml` file will be created in the current working directory referencing
the remote template. The file will also contain template values. Edit `rig.yaml`
to your liking and build the template to stdout:

```shell
rig build
```

To apply or overwrite template values when building and applying the manifests
to kubernetes use `--value` and pipe the output to `kubectl`:

```shell
rig build --value deployment.tag=$(git rev-parse HEAD) | kubectl apply -f -
```

## Installing

- Make sure git is installed
- Download release assets from the latest github release or install using go:

```shell
go get github.com/gonstr/rig
```
