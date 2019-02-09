## infra-breakout

Converts infra config files to their more complex companion files, namely:

- build.yaml to Jenkinsfile
- k8s.yaml to Kubernetes manifests

### Why?

Because you want to "breakout" from the simple configuration files and use Jenkinsfile or Kubernetes manifests directly.

### Usage

```shell
docker run fareoffice/infra-breakout help
docker run -v ${PWD}:/output fareoffice/infra-breakout build /output
docker run -v ${PWD}:/output fareoffice/infra-breakout k8s /output
```
