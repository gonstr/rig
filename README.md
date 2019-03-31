# rig

`rig` is a kubernetes manifest pre-processor and templating tool.

- Templates can be stored in remote git repositories, enabling the ability
  to share templates across apps.
- Built on go templates with support for all the templating features of Helm.

## Example usage with templates stored remotely

Install the remote template:

```shell
rig install https://github.com/foo/bar/simple-app#simple-app/v1.0.0
```

A `rig.yaml` file will be created in the current working directory referencing
the remote repository and gitref. The file will also contain template values.
Edit `rig.yaml` to your liking and build the template to stdout:

```shell
rig build
```

To apply or overwrite template values when building use `--value`:

```shell
rig build --value deployment.tag=$(git rev-parse HEAD) | kubectl apply -f -
```

## Example usage with templates stored locally

```yaml
echo "
apiVersion: v1
kind: Service
metadata:
  name: {{ .values.name }}
spec:
  ports:
  - port: {{ .values.host_port }}
    targetPort: {{ env "TARGET_PORT" }}
" > ./templates/service.yaml

echo "
template:
  path: templates
values:
  name: myservice
" > rig.yaml

TARGET_PORT=8080 rig build --string-value host_port=80
```

## Using the templating engine without rig.yaml

Rig can be used without rig.yaml if manifest path is set with `--path`:

```
rig build --path manifests/prod --value ingress.host=my-app.prod.com
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
