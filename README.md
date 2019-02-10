## infra-breakout

Converts infra config files to their more complex companion files, namely:

- build.yaml to Jenkinsfile
- k8s.yaml to Kubernetes manifests

### Why?

Because you want to "breakout" from the simple configuration files and use Jenkinsfile or Kubernetes manifests directly.

### Usage

```shell
bin/cli init
bin/cli breakout build.yaml
bin/cli breakout k8s.yaml ./deploy

or:

docker run -v ${PWD}:/cwd fareoffice/infra-cli init
docker run -v ${PWD}:/cwd fareoffice/infra-cli breakout /cwd/build.yaml
docker run -v ${PWD}:/cwd fareoffice/infra-cli breakout /cwd/k8s.yaml /cwd/deploy
```
