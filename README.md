# rig

Rig is a Helm inspired kubernetes manifest preprocessor and templating tool that
allows you to share versioned manifest templates across kubernetes apps.

### How does it work

Templates are stored in separate github repositories and can be shared by many
kubernetes applications.

Template url, gitref and values are stored in a `rig.yaml` file in the applications
repository.

The command line tool can download (using git) and build the remote template using
the values stored in `rig.yaml` and values passed as command line arguments.

Examples:

```shell
rig install https://github.com/foo/bar/simple-app#simple-app/v1.0.0

rig build --value deployment.tag=$(git rev-parse HEAD) | kubectl apply -f -
```

## Installing

- Install git
- Download release assets from the latest github release or install using go:

```shell
go get github.com/gonstr/rig
```
