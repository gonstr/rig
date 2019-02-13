## rig

Command line tool for generating and managing infra configuration files.

```shell
Usage
  $ rig init
  $ rig validate [file]
  $ rig gen-k8s [file] [cluster] [image tag]
  $ rig breakout [file] [output path] <template type>

Commands
  init
    Creates build.yaml and k8s.yaml files in the current working directory.

  validate
    Validates a build.yaml or k8s.yaml file.

  gen-k8s
    Generates kubernetes manifests from a k8s.ysml file and prints them to
    stdout.

  breakout
    Converts build.yaml or k8s.yaml to Jenkinsfile or Kubernetes manifest
    files respectivly and writes them to the specified output path. Output
    template can be specified as an optional third argument. Valid templates
    are:

    k8s.yaml:
      Kustomize (default) - converts k8s.yaml to kustomize yaml files.
      GoTemplates - converts k8s.yaml to go-template yaml files.

    build.yaml
      Jenkinsfile - converts build.yaml to Jenkinsfile.

Examples
  $ rig init
  $ rig validate build.yaml
  $ rig gen-k8s k8s.yaml int v1.2.3
  $ rig breakout k8s.yaml .
```

#### Install and use with npm

```shell
$ npm install -g git+ssh://git@github.com/fareoffice/rig.git
$ rig
```

#### Install and use with yarn

```shell
$ yarn global add git+ssh://git@github.com/fareoffice/rig.git
$ rig
```

#### Use with docker run

```shell
$ docker run -v ${PWD}:/cwd fareoffice/rig init /cwd
```
